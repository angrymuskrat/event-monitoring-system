package connector

import (
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"github.com/visheratin/unilog"
	"strings"

	types "github.com/angrymuskrat/event-monitoring-system/services/data-storage/data"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	_ "github.com/lib/pq"
)

type Storage struct {
	db     *sql.DB
	config Configuration
}


var (
	ErrDBConnecting = errors.New("do not be able to connect with db")
	ErrPushStatement = errors.New("one or more posts wasn't pushed")
	ErrSelectStatement = errors.New("don't be able to return posts")
	ErrPullGridStatement = errors.New("don't be able to return grid")
	ErrDuplicatedKey = errors.New("duplicated id, object hadn't saved to db")
)

func NewStorage(confPath string) (*Storage, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", conf.AuthDB)
	if err != nil {
		return nil, err
	}
	dbc := &Storage{db: db, config: conf}

	// create needed extension for PostGIS and TimescaleDB
	_, err = dbc.db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;")
	if err != nil {
		dbc.Close()
		return nil, err
	}
	_, err = dbc.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")
	if err != nil {
		dbc.Close()
		return nil, err
	}
	_, err = dbc.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis_topology;")
	if err != nil {
		dbc.Close()
		return nil, err
	}

	// create table posts with it's environment (hypertable and integer time now function)
	_, err = dbc.db.Exec(PostTable)
	if err != nil {
		dbc.Close()
		return nil, err
	}
	_, err = dbc.db.Exec("SELECT create_hypertable('posts', 'timestamp', chunk_time_interval => 86400, if_not_exists => TRUE);")
	if err != nil {
		dbc.Close()
		return nil, err
	}
	_, err = dbc.db.Exec("CREATE OR REPLACE FUNCTION unix_now() returns BIGINT LANGUAGE SQL STABLE as $$ SELECT extract(epoch from now())::BIGINT $$;")
	if err != nil {
		dbc.Close()
		return nil, err
	}
	_, err = dbc.db.Exec("SELECT set_integer_now_func('posts', 'unix_now', replace_if_exists => true);")
	if err != nil {
		dbc.Close()
		return nil, err
	}

	// create continuous aggregation of posts
	_, err = dbc.db.Exec("DROP VIEW aggr_posts CASCADE;")
	if err != nil {
		dbc.Close()
		return nil, err
	}
	createAggregationPosts := fmt.Sprintf(AggregationPosts, conf.AggrPostsGRIDSize, conf.AggrPostsGRIDSize) // set grid size
	_, err = dbc.db.Exec(createAggregationPosts)
	if err != nil {
		dbc.Close()
		return nil, err
	}

	// create table for grids
	_, err = dbc.db.Exec(GridTable)
	if err != nil {
		dbc.Close()
		return nil, err
	}

	unilog.Logger().Info("db connector has started")
	return dbc, nil
}

func (c *Storage) PushPosts(posts []data.Post) (ids []int32, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	wasError := false
	for _, v := range posts {
		statement := `
			INSERT INTO posts (ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, 
				AuthorID, LocationID, Location)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, ST_GeometryFromText($12, 4326));`
		point := fmt.Sprintf("POINT(%v %v)", v.Lat, v.Lon)
		_, err = c.db.Exec(statement, v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp,
			v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, point)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				ids = append(ids, types.DuplicatedPostId.Int32())
			} else {
				unilog.Logger().Error("don't be able to push post", zap.Any("post Shortcode", v.Shortcode), zap.Error(err))
				wasError = true
				ids = append(ids, types.DBError.Int32())
			}
		} else {
			ids = append(ids, types.PostPushed.Int32())
		}
	}
	if wasError {
		err = ErrPushStatement
	} else {
		err = nil
	}
	return ids, err
}

func (c Storage) SelectPosts(irv data.SpatioTemporalInterval) (posts []data.Post, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	poly := makePoly(irv.TopLeft, irv.BotRight)

	statement := fmt.Sprintf(`
		SELECT ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, 
			LocationID, ST_X(Location) as Lat, ST_Y(Location) as Lon
		FROM posts
		WHERE ST_Contains(%v, Location) AND (Timestamp BETWEEN %v AND %v)
	`, poly, irv.MinTime, irv.MaxTime)

	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in select posts", zap.Error(err))
		return nil, ErrSelectStatement
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select posts", zap.Error(err))
		}
	}()

	for rows.Next() {
		p := new(data.Post)
		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
			&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select posts", zap.Error(err))
			return nil, ErrSelectStatement
		}
		posts = append(posts, *p)
	}
	return posts, nil
}

func (c Storage) SelectAggrPosts(interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	poly := makePoly(interval.TopLeft, interval.BotRight)
	statement := fmt.Sprintf(SelectAggrPostsTemplate, interval.Hour, poly)

	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in select aggr_posts", zap.Error(err))
		return nil, ErrSelectStatement
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select aggr_posts", zap.Error(err))
		}
	}()

	for rows.Next() {
		p := new(struct{
			count int64
			lat float64
			lon float64
		})
		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&p.count, &p.lat, &p.lon)
		if err != nil {
			unilog.Logger().Error("error in select aggr_posts", zap.Error(err))
			return nil, ErrSelectStatement
		}
		post := data.AggregatedPost{Count:p.count, Center: data.Point{ Lat:p.lat, Lon:p.lon }}
		posts = append(posts, post)
	}
	return posts, nil
}

func (c *Storage) PushGrid(id string, blob []byte) (err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	statement := `
			INSERT INTO grids(ID, Blob)
			VALUES ($1, $2);`

	_, err = c.db.Exec(statement, id, blob)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicatedKey
		} else {
			unilog.Logger().Error("don't be able to push grid", zap.String("id", id), zap.Error(err))
		}
	}
	return err
}

func (c *Storage) PullGrid(id string) (blob []byte, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	statement := fmt.Sprintf("SELECT Blob FROM grids WHERE '%v' = Id;", id)
	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in pull grid", zap.Error(err))
		return nil, ErrPullGridStatement
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in pull grid", zap.Error(err))
		}
	}()

	ans := new(struct{Blob []byte})
	for rows.Next() {
		err = rows.Scan(&ans.Blob)
		if err != nil {
			unilog.Logger().Error("error in select", zap.Error(err))
			return nil, ErrPullGridStatement
		}
		break
	}
	blob = ans.Blob
	return
}

func (c *Storage) Close() {
	err := c.db.Close()
	if err != nil {
		unilog.Logger().Error("don't be able to close db", zap.Error(err))
	}
}


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

	createPostTable := `CREATE TABLE IF NOT EXISTS posts(
		ID varchar (30) not null primary key,
		Shortcode varchar (15),
		ImageURL varchar (300),
		IsVideo boolean not null,
		Caption varchar (2200), -- max size of text in Instagram
		CommentsCount bigint,
		Timestamp bigint,
		LikesCount bigint,
		IsAd boolean,
		AuthorID varchar (15),
		LocationID varchar (20),
		--Lat real,
		--Lon real
		Location geometry 
	)`

	_, err = dbc.db.Exec(createPostTable)
	if err != nil {
		dbc.Close()
		return nil, err
	}

	unilog.Logger().Info("db connector has started")
	return dbc, nil
}

func (c *Storage) Push(posts []data.Post) (ids []int32, err error) {
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
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, ST_GeometryFromText($12));`
		point := fmt.Sprintf("POINT(%v %v)", v.Lat, v.Lon)
		_, err = c.db.Exec(statement, v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp,
			v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, point)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				ids = append(ids, types.DuplicatedPostId.Int32())
			} else {
				unilog.Logger().Error("don't be able to push post", zap.String("post", fmt.Sprint(v)), zap.Error(err))
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

func (c Storage) Select(irv data.SpatioTemporalInterval) (posts []data.Post, err error) {

	poly := fmt.Sprintf("ST_GeometryFromText('POLYGON((%v %v, %v %v, %v %v, %v %v, %v %v))')",
		irv.MinLat, irv.MinLon, irv.MaxLat, irv.MinLon, irv.MaxLat, irv.MaxLon, irv.MinLat, irv.MaxLon,
		irv.MinLat, irv.MinLon)

	statement := fmt.Sprintf(`
		SELECT ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, 
			LocationID, ST_X(Location) as Lat, ST_Y(Location) as Lon
		FROM posts
		WHERE ST_Contains(%v, Location) AND (Timestamp BETWEEN %v AND %v)
	`, poly, irv.MinTime, irv.MaxTime)

	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in select", zap.Error(err))
		return nil, ErrSelectStatement
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select", zap.Error(err))
		}
	}()

	for rows.Next() {
		p := new(data.Post)
		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
			&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select", zap.Error(err))
			return nil, ErrSelectStatement
		}
		posts = append(posts, *p)
	}
	return posts, nil
}

func (c *Storage) Close() {
	err := c.db.Close()
	if err != nil {
		unilog.Logger().Error("don't be able to close db", zap.Error(err))
	}
}


package connector

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"strings"

	types "github.com/angrymuskrat/event-monitoring-system/services/data-storage/data"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/lib/pq"
)

type Storage struct {
	db     *sql.DB
	config Configuration
}

var (
	ErrDBConnecting    = errors.New("do not be able to connect with db")
	ErrDBTransaction   = errors.New("error with transaction")
	ErrPushPosts       = errors.New("one or more posts wasn't pushed")
	ErrSelectPosts     = errors.New("don't be able to return posts")
	ErrPullGrid        = errors.New("don't be able to return grid")
	ErrDuplicatedKey   = errors.New("duplicated id, object hadn't saved to db")
	ErrPushEvents      = errors.New("do not be able to insert events")
	ErrSelectEvents    = errors.New("don't be able to return events")
	ErrPushCity        = errors.New("do not be able to insert city")
	ErrPushLocations   = errors.New("do not be able to insert locations")
	ErrSelectLocations = errors.New("don't be able to return locations")
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

	err = setDBEnvironment(dbc, conf.AggrPostsGRIDSize)
	if err != nil {
		dbc.Close()
		return nil, err
	}

	unilog.Logger().Info("db connector has started")
	return dbc, nil
}

func setDBEnvironment(dbc *Storage, GRIDSize float64) (err error) {
	// create needed extension for PostGIS and TimescaleDB
	_, err = dbc.db.Exec(ExtensionTimescaleDB)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(ExtensionPostGIS)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(ExtensionPostGISTopology)
	if err != nil {
		return
	}

	_, err = dbc.db.Exec(CreateTimeFunction)
	if err != nil {
		return
	}

	// create table posts with it's environment (hypertable and integer time now function)
	_, err = dbc.db.Exec(PostTable)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(CreateHyperTablePosts)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(SetTimeFunctionForPosts)
	if err != nil {
		return
	}

	// create continuous aggregation of posts
	_, err = dbc.db.Exec(DropAggregationPosts)
	if err != nil {
		return
	}
	createAggregationPosts := fmt.Sprintf(AggregationPosts, GRIDSize, GRIDSize) // set grid size
	_, err = dbc.db.Exec(createAggregationPosts)
	if err != nil {
		return
	}

	// create events table
	_, err = dbc.db.Exec(EventsTable)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(CreateHyperTableEvents)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(SetTimeFunctionForEvents)
	if err != nil {
		return
	}

	// create tables for cities and locations
	_, err = dbc.db.Exec(CitiesTable)
	if err != nil {
		return
	}
	_, err = dbc.db.Exec(LocationsTable)
	if err != nil {
		return
	}

	// create table for grids
	_, err = dbc.db.Exec(GridTable)
	if err != nil {
		return
	}
	return nil
}

func (c *Storage) PushPosts(posts []data.Post) (ids []int32, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	tx, err := c.db.Begin()
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return nil, ErrDBTransaction
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(PushPostsTemplate)
	if err != nil {
		unilog.Logger().Error("can not prepare statement", zap.Error(err))
		return nil, ErrDBTransaction
	}

	for _, v := range posts {
		_, err = stmt.Exec(v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp, v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, v.Lat, v.Lon)
		if err != nil {
			unilog.Logger().Error("is not able to exec event", zap.Error(err))
			return nil, ErrPushPosts
		} else {
			ids = append(ids, types.PostPushed.Int32()) // TODO now this is useless
		}
	}
	if err := tx.Commit(); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return nil, ErrPushPosts
	}
	return
}

func (c Storage) SelectPosts(irv data.SpatioTemporalInterval) (posts []data.Post, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	poly := makePoly(irv.TopLeft, irv.BotRight)
	statement := fmt.Sprintf(SelectPostsTemplate, poly, irv.MinTime, irv.MaxTime)

	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in select posts", zap.Error(err))
		return nil, ErrSelectPosts
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
			return nil, ErrSelectPosts
		}
		posts = append(posts, *p)
	}
	return
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
		return nil, ErrSelectPosts
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select aggr_posts", zap.Error(err))
		}
	}()

	for rows.Next() {
		p := new(data.Point)
		ap := new(data.AggregatedPost)

		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&ap.Count, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select aggr_posts", zap.Error(err))
			return nil, ErrSelectPosts
		}
		ap.Center = *p
		posts = append(posts, *ap)
	}
	return posts, nil
}

func (c *Storage) PullTimeline(cityId string, start, finish int64) (timeline []data.Timestamp, err error) {
	return nil, nil
}

func (c *Storage) PushGrid(id string, blob []byte) (err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	statement := PushGridTemplate

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
	statement := fmt.Sprintf(PullGridTemplate, id)
	rows, err := c.db.Query(statement)
	if err != nil {
		unilog.Logger().Error("error in pull grid", zap.Error(err))
		return nil, ErrPullGrid
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in pull grid", zap.Error(err))
		}
	}()

	ans := new(struct{ Blob []byte })
	for rows.Next() {
		err = rows.Scan(&ans.Blob)
		if err != nil {
			unilog.Logger().Error("error in select", zap.Error(err))
			return nil, ErrPullGrid
		}
		break
	}
	blob = ans.Blob
	return
}

func (c *Storage) PushEvents(events []data.Event) (err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	tx, err := c.db.Begin()
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(PushEventsTemplate)
	if err != nil {
		unilog.Logger().Error("can not prepare statement", zap.Error(err))
		return ErrDBTransaction
	}
	for _, event := range events {
		_, err = stmt.Exec(event.Title, event.Start, event.Finish, event.Center.Lat, event.Center.Lon, pq.Array(event.PostCodes), pq.Array(event.Tags))
		if err != nil {
			unilog.Logger().Error("is not able to exec event", zap.Error(err))
			return ErrPushEvents
		}
	}
	if err := tx.Commit(); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return ErrPushEvents
	}
	return
}

func (c *Storage) PullEvents(interval data.SpatioHourInterval) (events []data.Event, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	poly := makePoly(interval.TopLeft, interval.BotRight)
	statement := fmt.Sprintf(SelectEventsTemplate, poly, interval.Hour, interval.Hour+Hour)
	rows, err := c.db.Query(statement)

	if err != nil {
		unilog.Logger().Error("error in select events", zap.Error(err))
		return nil, ErrSelectEvents
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select events", zap.Error(err))
		}
	}()

	for rows.Next() {
		e := new(data.Event)
		p := new(data.Point)
		err = rows.Scan(&e.Title, &e.Start, &e.Finish, pq.Array(&e.PostCodes), pq.Array(&e.Tags), &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select events", zap.Error(err))
			return nil, ErrSelectEvents
		}
		e.Center = *p
		events = append(events, *e)
	}
	return
}

func (c *Storage) PushLocations(city data.City, locations []data.Location) (err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	_, err = c.db.Exec(PushCityIfNotExists, city.Title, city.ID)
	if err != nil {
		unilog.Logger().Error("is not able to insert city", zap.Error(err))
		return ErrPushCity
	}

	tx, err := c.db.Begin()
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(PushLocationTemplate)
	if err != nil {
		unilog.Logger().Error("can not prepare statement", zap.Error(err))
		return ErrDBTransaction
	}
	for _, l := range locations {
		_, err = stmt.Exec(l.ID, city.ID, l.Position.Lat, l.Position.Lon, l.Title, l.Slug)
		if err != nil {
			unilog.Logger().Error("is not able to exec location", zap.Error(err))
			return ErrPushLocations
		}
	}
	if err := tx.Commit(); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return ErrPushLocations
	}
	return
}

func (c *Storage) PullLocations(cityId string) (locations []data.Location, err error) {
	err = c.db.Ping()
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	statement := fmt.Sprintf(SelectLocationsTemplate, cityId)
	rows, err := c.db.Query(statement)

	if err != nil {
		unilog.Logger().Error("error in select locations", zap.Error(err))
		return nil, ErrSelectLocations
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			unilog.Logger().Error("don't be able to close rows in select locations", zap.Error(err))
		}
	}()

	for rows.Next() {
		l := new(data.Location)
		p := new(data.Point)
		err = rows.Scan(&l.ID, &l.Title, &l.Slug, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select events", zap.Error(err))
			return nil, ErrSelectLocations
		}
		l.Position = p
		locations = append(locations, *l)
	}
	return
}

func (c *Storage) Close() {
	err := c.db.Close()
	if err != nil {
		unilog.Logger().Error("don't be able to close db", zap.Error(err))
	}
}

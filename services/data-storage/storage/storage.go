package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"strings"

	types "github.com/angrymuskrat/event-monitoring-system/services/data-storage/data"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	db     *pgx.Conn
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

func NewStorage(ctx context.Context, confPath string) (*Storage, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}
	connCongig, err := pgx.ParseConfig(conf.AuthDB)
	if err != nil {
		return nil, err
	}
	db, err := pgx.ConnectConfig(ctx, connCongig)
	if err != nil {
		return nil, err
	}
	dbc := &Storage{db: db, config: conf}

	err = dbc.setDBEnvironment(ctx, conf)
	if err != nil {
		dbc.Close(ctx)
		return nil, err
	}

	unilog.Logger().Info("db storage has started")
	return dbc, nil
}

func (c *Storage) setDBEnvironment(ctx context.Context, conf Configuration) (err error) {
	// create needed extension for PostGIS and TimescaleDB
	_, err = c.db.Exec(ctx, ExtensionTimescaleDB)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, ExtensionPostGIS)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, ExtensionPostGISTopology)
	if err != nil {
		return
	}

	_, err = c.db.Exec(ctx, CreateTimeFunction)
	if err != nil {
		return
	}

	// create table posts with it's environment (hypertable and integer time now function)
	_, err = c.db.Exec(ctx, PostTable)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, CreateHyperTablePosts)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, SetTimeFunctionForPosts)
	if err != nil {
		return
	}

	// create continuous aggregation of posts
	_, err = c.db.Exec(ctx, DropAggregationPosts)
	if err != nil {
		return
	}
	createAggregationPosts := fmt.Sprintf(AggregationPosts, conf.GRIDSize, conf.GRIDSize) // set grid size
	_, err = c.db.Exec(ctx, createAggregationPosts)
	if err != nil {
		return
	}

	// create events table
	_, err = c.db.Exec(ctx, EventsTable)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, CreateHyperTableEvents)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, SetTimeFunctionForEvents)
	if err != nil {
		return
	}

	// create tables for cities and locations
	_, err = c.db.Exec(ctx, CitiesTable)
	if err != nil {
		return
	}
	_, err = c.db.Exec(ctx, LocationsTable)
	if err != nil {
		return
	}

	// create table for grids
	_, err = c.db.Exec(ctx, GridTable)
	if err != nil {
		return
	}
	return nil
}

func (c *Storage) PushPosts(ctx context.Context, cityId string, posts []data.Post) (ids []int32, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	tx, err := c.db.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return nil, ErrDBTransaction
	}
	defer tx.Rollback(ctx)

	for _, v := range posts {
		_, err = tx.Exec(ctx, InsertPost, v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp, v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, v.Lat, v.Lon)
		if err != nil {
			unilog.Logger().Error("is not able to exec event", zap.Error(err))
			return nil, ErrPushPosts
		} else {
			ids = append(ids, types.PostPushed.Int32()) // TODO now this is useless
		}
	}
	if err := tx.Commit(ctx); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return nil, ErrPushPosts
	}
	return
}

func (c Storage) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) (posts []data.Post, cityArea *data.Area, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, nil, ErrDBConnecting
	}

	statement := fmt.Sprintf(SelectPosts, startTime, finishTime)

	rows, err := c.db.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in select posts", zap.Error(err))
		return nil, nil, ErrSelectPosts
	}
	defer rows.Close()

	for rows.Next() {
		p := new(data.Post)
		//(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
		err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
			&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select posts", zap.Error(err))
			return nil, nil, ErrSelectPosts
		}
		posts = append(posts, *p)
	}
	return
}

func (c Storage) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}
	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectAggrPosts, interval.Hour, poly)

	rows, err := c.db.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in select aggr_posts", zap.Error(err))
		return nil, ErrSelectPosts
	}

	defer rows.Close()

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

func (c *Storage) PullTimeline(ctx context.Context, cityId string, start, finish int64) (timeline []data.Timestamp, err error) {
	return nil, nil
}

func (c *Storage) PushGrid(ctx context.Context, cityId string, ids []int64, blobs [][]byte) (err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}

	for i, blob := range blobs {
		id := ids[i]
		_, err = c.db.Exec(ctx, PushGrid, id, blob)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return ErrDuplicatedKey
			} else {
				unilog.Logger().Error("don't be able to push grid", zap.Int64("id", id), zap.Error(err))
			}
			return err
		}
	}
	return err
}

func (c *Storage) PullGrid(ctx context.Context, cityId string, startId, finishId int64) (ids []int64, blobs [][]byte, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, nil, ErrDBConnecting
	}
	statement := fmt.Sprintf(PullGrid, startId, finishId)
	rows, err := c.db.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in pull grid", zap.Error(err))
		return nil, nil, ErrPullGrid
	}
	defer rows.Close()

	var id int64
	var blob []byte
	for rows.Next() {
		err = rows.Scan(&id, &blob)
		if err != nil {
			unilog.Logger().Error("error in select", zap.Error(err))
			return nil,nil, ErrPullGrid
		}
		ids = append(ids, id)
		blobs = append(blobs, blob)

	}
	return
}

func (c *Storage) PushEvents(ctx context.Context, cityId string, events []data.Event) (err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	tx, err := c.db.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)
	if err != nil {
		unilog.Logger().Error("can not prepare statement", zap.Error(err))
		return ErrDBTransaction
	}
	for _, event := range events {
		_, err = tx.Exec(ctx, PushEvents, event.Title, event.Start, event.Finish, event.Center.Lat, event.Center.Lon, pq.Array(event.PostCodes), pq.Array(event.Tags))
		if err != nil {
			unilog.Logger().Error("is not able to exec event", zap.Error(err))
			return ErrPushEvents
		}
	}
	if err := tx.Commit(ctx); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return ErrPushEvents
	}
	return
}

func (c *Storage) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) (events []data.Event, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectEvents, poly, interval.Hour, interval.Hour+Hour)
	rows, err := c.db.Query(ctx, statement)

	if err != nil {
		unilog.Logger().Error("error in select events", zap.Error(err))
		return nil, ErrSelectEvents
	}
	defer rows.Close()

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

func (c *Storage) PushLocations(ctx context.Context, cityId string, locations []data.Location) (err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return ErrDBConnecting
	}
	_, err = c.db.Exec(ctx, PushCityIfNotExists, cityId, cityId)
	if err != nil {
		unilog.Logger().Error("is not able to insert city", zap.Error(err))
		return ErrPushCity
	}

	tx, err := c.db.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)
	for _, l := range locations {
		_, err = tx.Exec(ctx, PushLocation, l.ID, cityId, l.Position.Lat, l.Position.Lon, l.Title, l.Slug)
		if err != nil {
			unilog.Logger().Error("is not able to exec location", zap.Error(err))
			return ErrPushLocations
		}
	}
	if err := tx.Commit(ctx); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return ErrPushLocations
	}
	return
}

func (c *Storage) PullLocations(ctx context.Context, cityId string) (locations []data.Location, err error) {
	err = c.db.Ping(ctx)
	if err != nil {
		unilog.Logger().Error("db error", zap.Error(err))
		return nil, ErrDBConnecting
	}

	statement := fmt.Sprintf(SelectLocations, cityId)
	rows, err := c.db.Query(ctx, statement)

	if err != nil {
		unilog.Logger().Error("error in select locations", zap.Error(err))
		return nil, ErrSelectLocations
	}
	defer rows.Close()

	for rows.Next() {
		l := new(data.Location)
		p := new(data.Point)
		err = rows.Scan(&l.ID, &l.Title, &l.Slug, &p.Lat, &p.Lon)
		if err != nil {
			unilog.Logger().Error("error in select events", zap.Error(err))
			return nil, ErrSelectLocations
		}
		l.Position = *p
		locations = append(locations, *l)
	}
	return
}

func (c *Storage) Close(ctx context.Context) {
	err := c.db.Close(ctx)
	if err != nil {
		unilog.Logger().Error("don't be able to close db", zap.Error(err))
	}
}

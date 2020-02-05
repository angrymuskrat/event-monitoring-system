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
	general *pgx.Conn
	cities  map[string]*pgx.Conn
	config  Configuration
}

var (
	ErrUnexpectedCity  = errors.New("this city doesn't exist")
	ErrDBTransaction   = errors.New("error with transaction")
	ErrPushPosts       = errors.New("one or more posts wasn't pushed")
	ErrSelectPosts     = errors.New("don't be able to return posts")
	ErrPullGrid        = errors.New("don't be able to return grid")
	ErrDuplicatedKey   = errors.New("duplicated id, object hadn't saved to db")
	ErrPushEvents      = errors.New("do not be able to insert events")
	ErrSelectEvents    = errors.New("don't be able to return events")
	ErrPushLocations   = errors.New("do not be able to insert locations")
	ErrSelectLocations = errors.New("don't be able to return locations")
)

func NewStorage(ctx context.Context, confPath string) (*Storage, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}
	s := &Storage{config: conf}
	err = s.initGeneral(ctx)
	if err != nil {
		s.Close(ctx)
		return nil, err
	}

	err = s.initCities(ctx)
	if err != nil {
		s.Close(ctx)
		return nil, err
	}

	unilog.Logger().Info("db storage has started")
	return s, nil
}

func (s *Storage) initGeneral(ctx context.Context) (err error) {
	connConfig, err := pgx.ParseConfig(s.config.makeAuthToken(PostgresDBName))
	if err != nil {
		return err
	}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return err
	}

	var name string
	row := conn.QueryRow(ctx, makeSelectDB(GeneralDBName))
	err = row.Scan(&name)
	if err == pgx.ErrNoRows {
		_, err = conn.Exec(ctx, makeCreateDB(GeneralDBName))
	}
	if err != nil {
		return err
	}

	connConfig, err = pgx.ParseConfig(s.config.makeAuthToken(GeneralDBName))
	if err != nil {
		return err
	}
	conn, err = pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, ExtensionPostGIS)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, CreateCitiesTable)
	if err != nil {
		return err
	}

	s.general = conn
	return nil
}

func (s *Storage) initCities(ctx context.Context) (err error) {
	cities, err := s.GetCities(ctx)
	if err != nil {
		return err
	}
	s.cities = make(map[string]*pgx.Conn)

	for _, city := range cities {
		cityId := city.Code
		row := s.general.QueryRow(ctx, makeSelectDB(cityId))
		err = row.Scan(&cityId)
		if err == pgx.ErrNoRows {
			_, err = s.general.Exec(ctx, makeCreateDB(cityId))
		}
		if err != nil {
			return err
		}

		connConfig, err := pgx.ParseConfig(s.config.makeAuthToken(cityId))
		if err != nil {
			return nil
		}
		conn, err := pgx.ConnectConfig(ctx, connConfig)
		if err != nil {
			return nil
		}
		s.cities[cityId] = conn

		err = s.setCityEnvironment(ctx, cityId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) getCityConn(cityId string) (*pgx.Conn, error) {
	if conn, isExist := s.cities[cityId]; isExist {
		return conn, nil
	} else {
		return nil, ErrUnexpectedCity
	}
}

func (s *Storage) setCityEnvironment(ctx context.Context, cityId string) (err error) {
	conn, err := s.getCityConn(cityId)

	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, ExtensionTimescaleDB)
	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, ExtensionPostGIS)
	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, ExtensionPostGISTopology)
	if err != nil {
		return
	}

	_, err = conn.Exec(ctx, CreateTimeFunction)
	if err != nil {
		return
	}

	// create table posts with it's environment (hypertable and integer time now function)
	_, err = conn.Exec(ctx, CreatePostsTable)
	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, CreateHyperTablePosts)
	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, SetTimeFunctionForPosts)
	if err != nil {
		return
	}

	createAggrPostsView := fmt.Sprintf(CreateAggrPostsView, s.config.GRIDSize, s.config.GRIDSize) // set grid size
	_, err = conn.Exec(ctx, createAggrPostsView)
	if err != nil && isNotAlreadyExistsError(err) {
		return
	}

	// create events table
	_, err = conn.Exec(ctx, CreateEventsTable)
	if err != nil {
		return
	}
	_, err = conn.Exec(ctx, CreateHyperTableEvents)
	if err != nil {
		return
	}

	_, err = conn.Exec(ctx, SetTimeFunctionForEvents)
	if err != nil {
		return
	}

	_, err = conn.Exec(ctx, CreatePostsTimelineView)
	if err != nil && isNotAlreadyExistsError(err) {
		return
	}

	_, err = conn.Exec(ctx, CreateEventsTimelineView)
	if err != nil && isNotAlreadyExistsError(err) {
		return
	}

	// create table for locations
	_, err = conn.Exec(ctx, CreateLocationsTable)
	if err != nil {
		return
	}

	// create table for grids
	_, err = conn.Exec(ctx, CreateGridsTable)
	if err != nil {
		return
	}
	return nil
}

func (s *Storage) GetCities(ctx context.Context) (cities []data.City, err error) {
	rows, err := s.general.Query(ctx, SelectCities)
	if err != nil {
		unilog.Logger().Error("error in GetCities - Query", zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		var (
			city   data.City
			tl, br data.Point
		)
		err = rows.Scan(&city.Title, &city.Code, &tl.Lat, &tl.Lon, &br.Lat, &br.Lon)
		if err != nil {
			unilog.Logger().Error("error in GetCities - Scan", zap.Error(err))
			return nil, err
		}
		city.Area = data.Area{TopLeft: &tl, BotRight: &br}
		cities = append(cities, city)
	}
	return cities, nil
}

func (s *Storage) PushPosts(ctx context.Context, cityId string, posts []data.Post) (ids []int32, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return nil, ErrDBTransaction
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			unilog.Logger().Error("do not be able to close transaction", zap.Error(err))
		}
	}()

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

func (s Storage) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) (posts []data.Post, cityArea *data.Area, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, nil, err
	}

	statement := fmt.Sprintf(SelectPosts, startTime, finishTime)
	rows, err := conn.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in select posts", zap.Error(err))
		return nil, nil, ErrSelectPosts
	}
	defer rows.Close()

	for rows.Next() {
		p := new(data.Post)
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

func (s Storage) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectAggrPosts, interval.Hour, poly)

	rows, err := conn.Query(ctx, statement)
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

func (s *Storage) PullTimeline(ctx context.Context, cityId string, start, finish int64) (timeline []data.Timestamp, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	statement := fmt.Sprintf(SelectTimeline, start, finish, start, finish)
	rows, err := conn.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in pull timeline", zap.Error(err))
		return nil, ErrPullGrid
	}
	for rows.Next() {
		var timestamp data.Timestamp
		err = rows.Scan(&timestamp.PostsNumber, &timestamp.EventsNumber, &timestamp.Time)
		if err != nil {
			unilog.Logger().Error("error in pull timeline", zap.Error(err))
			return  nil, err
		}
		timeline = append(timeline, timestamp)
	}
	return
}

func (s *Storage) PushGrid(ctx context.Context, cityId string, ids []int64, blobs [][]byte) (err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}
	// TODO add transaction
	for i, blob := range blobs {
		id := ids[i]
		_, err = conn.Exec(ctx, InsertGrid, id, blob)
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

func (s *Storage) PullGrid(ctx context.Context, cityId string, startId, finishId int64) (ids []int64, blobs [][]byte, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, nil, err
	}
	statement := fmt.Sprintf(SelectGrid, startId, finishId)
	rows, err := conn.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in pull grid", zap.Error(err))
		return nil, nil, ErrPullGrid
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var blob []byte
		err = rows.Scan(&id, &blob)
		if err != nil {
			unilog.Logger().Error("error in pull grid", zap.Error(err))
			return nil, nil, ErrPullGrid
		}
		ids = append(ids, id)
		blobs = append(blobs, blob)

	}
	return
}

func (s *Storage) PushEvents(ctx context.Context, cityId string, events []data.Event) (err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			unilog.Logger().Error("do not be able to close transaction", zap.Error(err))
		}
	}()
	for _, event := range events {
		_, err = tx.Exec(ctx, InsertEvent, event.Title, event.Start, event.Finish, event.Center.Lat, event.Center.Lon, pq.Array(event.PostCodes), pq.Array(event.Tags))
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

func (s *Storage) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) (events []data.Event, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectEvents, poly, interval.Hour, interval.Hour+Hour)
	rows, err := conn.Query(ctx, statement)

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

func (s *Storage) PushLocations(ctx context.Context, cityId string, locations []data.Location) (err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return  err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			unilog.Logger().Error("do not be able to close transaction", zap.Error(err))
		}
	}()
	for _, l := range locations {
		_, err = tx.Exec(ctx, InsertLocation, l.ID, l.Position.Lat, l.Position.Lon, l.Title, l.Slug)
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

func (s *Storage) PullLocations(ctx context.Context, cityId string) (locations []data.Location, err error) {
	conn, err := s.getCityConn(cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}
	rows, err := conn.Query(ctx, SelectLocations)

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

func (s *Storage) Close(ctx context.Context) {
	if s.general == nil {
		return
	}
	err := s.general.Close(ctx)
	if err != nil {
		unilog.Logger().Error("don't be able to close general conn", zap.Error(err))
	}
	for cityId, conn := range s.cities {
		err = conn.Close(ctx)
		if err != nil {
			unilog.Logger().Error("don't be able to city conn", zap.String("cityId", cityId), zap.Error(err))
		}
	}
}

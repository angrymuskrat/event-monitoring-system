package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

type Storage struct {
	general *pgxpool.Pool
	cities  map[string]*pgxpool.Pool
	config  Configuration
}

var (
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

func New(ctx context.Context, confPath string) (*Storage, error) {
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
	connConfig, err := pgxpool.ParseConfig(s.config.makeAuthToken(PostgresDBName))
	if err != nil {
		return err
	}
	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return err
	}

	var name string
	row := conn.QueryRow(ctx, MakeSelectDB(GeneralDBName))
	err = row.Scan(&name)
	if err == pgx.ErrNoRows {
		_, err = conn.Exec(ctx, MakeCreateDB(GeneralDBName))
	}
	if err != nil {
		return err
	}

	connConfig, err = pgxpool.ParseConfig(s.config.makeAuthToken(GeneralDBName))
	if err != nil {
		return err
	}
	conn, err = pgxpool.ConnectConfig(ctx, connConfig)
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
	s.cities = make(map[string]*pgxpool.Pool)
	for _, city := range cities {
		cityId := city.Code
		err = s.initCity(ctx, cityId)
		if err != nil {
			return
		}
	}
	return
}

func (s *Storage) initCity(ctx context.Context, cityID string) error {
	row := s.general.QueryRow(ctx, MakeSelectDB(cityID))
	err := row.Scan(&cityID)
	if err == pgx.ErrNoRows {
		_, err = s.general.Exec(ctx, MakeCreateDB(cityID))
	}
	if err != nil {
		unilog.Logger().Error("unable to create database for the city")
		return err
	}
	connConfig, err := pgxpool.ParseConfig(s.config.makeAuthToken(cityID))
	if err != nil {
		unilog.Logger().Error("unable to parse config for database connection")
		return err
	}
	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		unilog.Logger().Error("unable to connect to the city database")
		return nil
	}
	s.cities[cityID] = conn
	err = s.setCityEnvironment(ctx, cityID)
	if err != nil {
		unilog.Logger().Error("unable to set city database environment")
	}
	return err
}

func (s *Storage) getCityConn(ctx context.Context, cityID string) (*pgxpool.Pool, error) {
	if conn, isExist := s.cities[cityID]; isExist {
		return conn, nil
	}
	if err := s.initCity(ctx, cityID); err != nil { //TODO: possible concurrent map write, not a deal for current use case.
		return nil, err
	}
	if conn, isExist := s.cities[cityID]; isExist {
		return conn, nil
	}
	return nil, errors.New("specified city does not exist in the database")
}

func (s *Storage) setCityEnvironment(ctx context.Context, cityId string) (err error) {
	conn, err := s.getCityConn(ctx, cityId)

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

	createAggrPostsView := fmt.Sprintf(CreateAggrPostsView, s.config.RefreshInterval, s.config.GRIDSize, s.config.GRIDSize) // set grid size
	_, err = conn.Exec(ctx, createAggrPostsView)
	if err != nil && IsNotAlreadyExistsError(err) {
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

	createPostsTimelineView := fmt.Sprintf(CreatePostsTimelineView, s.config.RefreshInterval)
	_, err = conn.Exec(ctx, createPostsTimelineView)
	if err != nil && IsNotAlreadyExistsError(err) {
		return
	}

	createEventsTimelineView := fmt.Sprintf(CreateEventsTimelineView, s.config.RefreshInterval)
	_, err = conn.Exec(ctx, createEventsTimelineView)
	if err != nil && IsNotAlreadyExistsError(err) {
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

func (s *Storage) InsertCity(ctx context.Context, city data.City, updateIfExist bool) (err error) {
	tl := city.Area.TopLeft
	br := city.Area.BotRight
	var statement string
	if updateIfExist {
		statement = UpsertCity
	} else {
		statement = InsertCity
	}
	_, err = s.general.Exec(ctx, statement, city.Title, city.Code, tl.Lat, tl.Lon, br.Lat, br.Lon)
	if err != nil {
		unilog.Logger().Error("error in InsertCity", zap.Error(err))
		return
	}
	return nil
}

func (s *Storage) SelectCity(ctx context.Context, cityId string) (city *data.City, err error) {
	statement := fmt.Sprintf(SelectCity, cityId)
	row := s.general.QueryRow(ctx, statement)
	var tl, br data.Point
	city = &data.City{}

	err = row.Scan(&city.Title, &city.Code, &tl.Lat, &tl.Lon, &br.Lat, &br.Lon)
	if err != nil {
		unilog.Logger().Error("error in selectCity", zap.Error(err))
		return nil, err
	}
	city.Area = data.Area{TopLeft: &tl, BotRight: &br}
	return city, nil
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

func (s *Storage) PushPosts(ctx context.Context, cityId string, posts []data.Post) (err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)

	for _, v := range posts {
		_, err = tx.Exec(ctx, InsertPost, v.ID, v.Shortcode, v.ImageURL, v.IsVideo, v.Caption, v.CommentsCount, v.Timestamp, v.LikesCount, v.IsAd, v.AuthorID, v.LocationID, v.Lat, v.Lon)
		if err != nil {
			unilog.Logger().Error("is not able to exec event", zap.Error(err))
			return ErrPushPosts
		}
	}
	if err := tx.Commit(ctx); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return ErrPushPosts
	}
	return
}

func (s Storage) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) (posts []data.Post, cityArea *data.Area, err error) {
	conn, err := s.getCityConn(ctx, cityId)
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
	city, err := s.SelectCity(ctx, cityId)
	if err != nil {
		// unilog inn't needed due to unilog in SelectCity
		return nil, nil, err
	}
	return posts, &city.Area, nil
}

func (s Storage) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	poly := MakePoly(interval.Area)
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
	conn, err := s.getCityConn(ctx, cityId)
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
			return nil, err
		}
		timeline = append(timeline, timestamp)
	}
	return
}

func (s *Storage) PushGrid(ctx context.Context, cityId string, grids map[int64][]byte) (err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)

	for id, blob := range grids {
		_, err = tx.Exec(ctx, InsertGrid, id, blob)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return ErrDuplicatedKey
			} else {
				unilog.Logger().Error("don't be able to push grid", zap.Int64("id", id), zap.Error(err))
			}
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		unilog.Logger().Error("is not able to commit events transaction", zap.Error(err))
		return err
	}
	return err
}

func (s *Storage) PullGrid(ctx context.Context, cityId string, ids []int64) (grids map[int64][]byte, err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}
	grids = make(map[int64][]byte)
	statement := formSelectGrids(ids)
	rows, err := conn.Query(ctx, statement)
	if err != nil {
		unilog.Logger().Error("error in pull grid", zap.Error(err))
		return nil, ErrPullGrid
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var blob []byte
		err = rows.Scan(&id, &blob)
		if err != nil {
			unilog.Logger().Error("error in pull grid", zap.Error(err))
			return nil, ErrPullGrid
		}
		grids[id] = blob
	}
	return
}

func formSelectGrids(ids []int64) string {
	s := "SELECT id, blob FROM grids WHERE id IN ("
	f := ");"
	res := s
	for i := range ids {
		res += strconv.FormatInt(ids[i], 10)
		if i < (len(ids) - 1) {
			res += ","
		}
	}
	res += f
	return res
}

func (s *Storage) PushEvents(ctx context.Context, cityId string, events []data.Event) (err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)

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
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	poly := MakePoly(interval.Area)
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

func (s *Storage) PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) (events []data.Event, err error) {
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return nil, err
	}

	statement := MakeSelectEventsTags(tags, startTime, finishTime)
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
	conn, err := s.getCityConn(ctx, cityId)
	if err != nil {
		unilog.Logger().Error("unexpected cityId", zap.String("cityId", cityId), zap.Error(err))
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		unilog.Logger().Error("can not begin transaction", zap.Error(err))
		return ErrDBTransaction
	}
	defer tx.Rollback(ctx)

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
	conn, err := s.getCityConn(ctx, cityId)
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
	s.general.Close()
	for _, conn := range s.cities {
		conn.Close()
	}
}

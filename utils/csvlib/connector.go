package csvlib

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"os"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

var (
	ErrSelectPosts  = errors.New("don't be able to return posts")
	ErrSelectEvents = errors.New("don't be able to return events")
)

type connector struct {
	conn   *pgxpool.Pool
	dbName string
	config *configuration
}

func NewConnector(ctx context.Context, configPath *string, dbName string) (c *connector, err error) {
	config, err := readConfig(*configPath)
	connConfig, err := pgxpool.ParseConfig(config.makeAuthToken(dbName))
	if err != nil {
		unilog.Logger().Error("unable to parse pg config",
			zap.String("dbName", dbName), zap.Error(err))
		return
	}
	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		unilog.Logger().Error("unable to connect to db",
			zap.String("dbName", dbName), zap.Error(err))
		return
	}
	c = &connector{
		conn:   conn,
		dbName: dbName,
		config: config,
	}
	return c, nil
}

func (c *connector) selectPostsByCodes(ctx context.Context, codes []string,
	start, finish int64) (posts []data.Post, err error) {
	rows, err := c.conn.Query(ctx, makeSelectPostsByCodesSQL(codes, start, finish, c.config.HasNoiseAndUtility))
	if err != nil {
		unilog.Logger().Error("error in select posts", zap.Error(err))
		return nil, ErrSelectPosts
	}
	defer rows.Close()

	for rows.Next() {
		p := new(data.Post)
		if c.config.HasNoiseAndUtility {
			err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
				&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lon, &p.Lat, &p.NoiseProbability, &p.EventUtility)
		} else {
			err = rows.Scan(&p.ID, &p.Shortcode, &p.ImageURL, &p.IsVideo, &p.Caption, &p.CommentsCount, &p.Timestamp,
				&p.LikesCount, &p.IsAd, &p.AuthorID, &p.LocationID, &p.Lon, &p.Lat)
		}
		if err != nil {
			unilog.Logger().Error("error in select posts", zap.Error(err))
			return nil, ErrSelectPosts
		}
		posts = append(posts, *p)
	}
	return posts, nil
}

func (c *connector) selectAllEvents(ctx context.Context, eventTable string) (events []data.Event, err error) {
	rows, err := c.conn.Query(ctx, makeSelectAllEventsSQL(eventTable))

	if err != nil {
		unilog.Logger().Error("error in select events", zap.Error(err))
		return nil, ErrSelectEvents
	}
	defer rows.Close()

	for rows.Next() {
		e := new(data.Event)
		p := new(data.Point)
		err = rows.Scan(&e.Title, &e.Start, &e.Finish, pq.Array(&e.PostCodes))
		if err != nil {
			unilog.Logger().Error("error in select events", zap.Error(err))
			return nil, ErrSelectEvents
		}
		e.Center = *p
		events = append(events, *e)
	}
	return events, nil
}

func (c *connector) eventPosts(ctx context.Context, rootPath, eventTableName, outputFile string) error {
	unilog.Logger().Info("process status: start collecting events")
	events, err := c.selectAllEvents(ctx, eventTableName)
	unilog.Logger().Info("process status: finish collecting events")
	if err != nil {
		return err
	}
	percentage := 0
	percentageStep := 5
	sum := 0

	f, err := os.Create(rootPath + outputFile)
	defer func() {
		err := f.Close()
		if err != nil {
			unilog.Logger().Error("don't be able to close file", zap.Error(err))
		}
	}()

	if err != nil {
		unilog.Logger().Error("don't be able to open file", zap.Error(err))
		return err
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	columns := []string{"code", "caption", "lat", "lon", "author_id", "location_id", "timestamp", "city", "event_ind", "event_title"}
	if c.config.HasNoiseAndUtility {
		columns = append(columns, []string{"noise_probability", "event_utility"}...)
	}
	if err := w.Write(columns); err != nil {
		unilog.Logger().Error("error writing record to file", zap.Error(err))
		return err
	}

	unilog.Logger().Info("process status: start collecting event posts")
	for ind, e := range events {
		posts, err := c.selectPostsByCodes(ctx, e.PostCodes, e.Start, e.Finish)
		if err != nil {
			unilog.Logger().Error("don't be able to load posts for the event",
				zap.String("city", c.dbName), zap.Int("eventInd", ind), zap.Error(err))
			return err
		}

		for _, p := range posts {
			record := []string{p.Shortcode, p.Caption, fmt.Sprintf("%v", p.Lat), fmt.Sprintf("%v", p.Lon),
				p.AuthorID, p.LocationID, fmt.Sprintf("%v", p.Timestamp), c.dbName, fmt.Sprintf("%v", ind),
				events[ind].Title}
			if c.config.HasNoiseAndUtility {
				record = append(record, []string{
					fmt.Sprintf("%v", p.NoiseProbability),
					fmt.Sprintf("%v", p.EventUtility)}...)
			}
			if err := w.Write(record); err != nil {
				unilog.Logger().Error("error writing record to file", zap.Error(err))
				return err
			}
		}

		sum += len(posts)
		progress := int(float32(ind) / (float32(len(events)) / 100))
		if percentage+percentageStep <= progress {
			percentage = progress
			unilog.Logger().Info("process status: collecting event posts", zap.Int("progress, %", percentage),
				zap.Int("eventInd", ind), zap.Int("posts written", sum))
		}
	}
	return nil
}

func (c *connector) ExecuteRequest(ctx context.Context, requestType string, rootPath string,
	additionalParams map[string]*string) (err error) {
	switch requestType {
	case "LoadEventPosts":
		err = c.eventPosts(ctx, rootPath, *additionalParams["EventTableName"], *additionalParams["EventPostsOutput"])
	default:
		err = errors.New("unknown type flag option")
	}
	return err
}

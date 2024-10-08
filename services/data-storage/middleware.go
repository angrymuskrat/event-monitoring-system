package service

import (
	"context"
	"time"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   Service
}

func (mw loggingMiddleware) InsertCity(ctx context.Context, city data.City, updateIfExists bool) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("insert city request",
			zap.Any("city", city),
			zap.Bool("update option", updateIfExists),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.InsertCity(ctx, city, updateIfExists)
	return
}

func (mw loggingMiddleware) GetAllCities(ctx context.Context) (cities []data.City, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("get all cities request",
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	cities, err = mw.next.GetAllCities(ctx)
	return
}

func (mw loggingMiddleware) GetCity(ctx context.Context, cityId string) (city *data.City, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("get city by id",
			zap.String("cityId", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	city, err = mw.next.GetCity(ctx, cityId)
	return
}

func (mw loggingMiddleware) PushPosts(ctx context.Context, cityId string, posts []data.Post) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push posts request",
			zap.Int("count of posts", len(posts)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushPosts(ctx, cityId, posts)
	return
}

func (mw loggingMiddleware) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) (posts []data.Post, area *data.Area, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select posts",
			zap.Int64("start time", startTime),
			zap.Int64("finish time", finishTime),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, area, err = mw.next.SelectPosts(ctx, cityId, startTime, finishTime)
	return
}

func (mw loggingMiddleware) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select aggregated posts",
			zap.Any("interval", interval),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.SelectAggrPosts(ctx, cityId, interval)
	return
}

func (mw loggingMiddleware) PullTimeline(ctx context.Context, cityId string, start, finish int64) (timeline []data.Timestamp, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select aggregated posts",
			zap.String("City", cityId),
			zap.Int64("Start time", start),
			zap.Int64("Finish time", finish),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	timeline, err = mw.next.PullTimeline(ctx, cityId, start, finish)
	return
}

func (mw loggingMiddleware) PushGrid(ctx context.Context, cityId string, grids map[int64][]byte) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push grid",
			zap.Int("len grids", len(grids)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushGrid(ctx, cityId, grids)
	return
}

func (mw loggingMiddleware) PullGrid(ctx context.Context, cityId string, ids []int64) (grids map[int64][]byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull grid",
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	grids, err = mw.next.PullGrid(ctx, cityId, ids)
	return
}

func (mw loggingMiddleware) PushEvents(ctx context.Context, cityId string, events []data.Event) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push events",
			zap.Int("len of events", len(events)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushEvents(ctx, cityId, events)
	return
}

func (mw loggingMiddleware) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) (events []data.Event, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull events",
			zap.Any("interval", interval),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	events, err = mw.next.PullEvents(ctx, cityId, interval)
	return
}

func (mw loggingMiddleware) PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) (events []data.Event, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull events by tags",
			zap.String("city id", cityId),
			zap.Int64("Start time", startTime),
			zap.Int64("Finish time", finishTime),
			zap.Any("tags", tags))
	}(time.Now())
	events, err = mw.next.PullEventsTags(ctx, cityId, tags, startTime, finishTime)
	return
}

func (mw loggingMiddleware) PushLocations(ctx context.Context, cityId string, locations []data.Location) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push locations",
			zap.Int("locations size", len(locations)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushLocations(ctx, cityId, locations)
	return
}

func (mw loggingMiddleware) PullLocations(ctx context.Context, cityId string) (locations []data.Location, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull locations",
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	locations, err = mw.next.PullLocations(ctx, cityId)
	return
}

func (mw loggingMiddleware) PullShortPostInInterval(ctx context.Context, cityId string, shortCodes []string,
	startTimestamp int64, endTimestamp int64) (posts []data.ShortPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull short posts in interval",
			zap.String("city id", cityId),
			zap.Int("len of codes", len(shortCodes)),
			zap.Int64("start of interval", startTimestamp),
			zap.Int64("end of interval", endTimestamp),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.PullShortPostInInterval(ctx, cityId, shortCodes, startTimestamp, endTimestamp)
	return
}

func (mw loggingMiddleware) PullSingleShortPost(ctx context.Context, cityId, shortcode string) (post *data.ShortPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull short posts in interval",
			zap.String("city id", cityId),
			zap.String("shortcode", shortcode),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	post, err = mw.next.PullSingleShortPost(ctx, cityId, shortcode)
	return
}

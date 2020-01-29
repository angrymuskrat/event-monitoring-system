package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"go.uber.org/zap"
	"time"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   Service
}

func (mw loggingMiddleware) PushPosts(ctx context.Context, posts []data.Post) (ids []int32, err error) {
	func(begin time.Time) {
		mw.logger.Info("push posts request",
			zap.Int("count of posts", len(posts)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	ids, err = mw.next.PushPosts(ctx, posts)
	return
}

func (mw loggingMiddleware) SelectPosts(ctx context.Context, interval data.SpatioTemporalInterval) (posts []data.Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select posts",
			zap.Any("interval", interval),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.SelectPosts(ctx, interval)
	return
}

func (mw loggingMiddleware) SelectAggrPosts(ctx context.Context, interval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select aggregated posts",
			zap.Any("interval", interval),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.SelectAggrPosts(ctx, interval)
	return
}

func (mw loggingMiddleware) PullTimeline(ctx context.Context, cityId string, start, finish int64) (timeline []data.Timestamp, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select aggregated posts",
			zap.String("City", cityId),
			zap.Int64("Start time", start),
			zap.Int64("Finish time", finish),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	timeline, err = mw.next.PullTimeline(ctx, cityId, start, finish)
	return
}

func (mw loggingMiddleware) PushGrid(ctx context.Context, id string, blob []byte) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push grid",
			zap.String("id", id),
			zap.Int("capacity", len(blob)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushGrid(ctx, id, blob)
	return
}

func (mw loggingMiddleware) PullGrid(ctx context.Context, id string) (blob []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull grid",
			zap.String("id", id),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	blob, err = mw.next.PullGrid(ctx, id)
	return
}

func (mw loggingMiddleware) PushEvents(ctx context.Context, events []data.Event) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push events",
			zap.Int("len of events", len(events)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushEvents(ctx, events)
	return
}

func (mw loggingMiddleware) PullEvents(ctx context.Context, interval data.SpatioHourInterval) (events []data.Event, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull events",
			zap.Any("interval", interval),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	events, err = mw.next.PullEvents(ctx, interval)
	return
}

func (mw loggingMiddleware) PushLocations(ctx context.Context, city data.City, locations []data.Location) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push locations",
			zap.Any("city", city),
			zap.Int("locations size", len(locations)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushLocations(ctx, city, locations)
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

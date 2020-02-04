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

func (mw loggingMiddleware) PushPosts(ctx context.Context, cityId string, posts []data.Post) (ids []int32, err error) {
	func(begin time.Time) {
		mw.logger.Info("push posts request",
			zap.Int("count of posts", len(posts)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	ids, err = mw.next.PushPosts(ctx, cityId, posts)
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

func (mw loggingMiddleware) PushGrid(ctx context.Context, cityId string, ids []int64, blobs [][]byte) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info("push grid",
			zap.Int("len ids", len(ids)),
			zap.Int("len blobs", len(blobs)),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	err = mw.next.PushGrid(ctx, cityId, ids, blobs)
	return
}

func (mw loggingMiddleware) PullGrid(ctx context.Context, cityId string, startId, finishId int64) (ids []int64, blobs [][]byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("pull grid",
			zap.Int64("start id", startId),
			zap.Int64("start id", finishId),
			zap.String("city id", cityId),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	ids, blobs, err = mw.next.PullGrid(ctx, cityId, startId, finishId)
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

package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next dbsvc.Service) dbsvc.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   dbsvc.Service
	logger log.Logger
}

func (mw loggingMiddleware) Push(ctx context.Context, posts []dbsvc.Post) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Push", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Push(ctx, posts)
}

func (mw loggingMiddleware) Select(ctx context.Context, interval dbsvc.SpatialTemporalInterval) (posts []dbsvc.Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Select", "Interval", interval, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Select(ctx, interval)
}

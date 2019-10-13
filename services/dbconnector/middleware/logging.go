package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/angrymuskrat/event-monitoring-system/services/dbconnector"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next dbconnector.Service) dbconnector.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   dbconnector.Service
	logger log.Logger
}

func (mw loggingMiddleware) Push(ctx context.Context, posts []dbconnector.Post) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Push", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Push(ctx, posts)
}

func (mw loggingMiddleware) Select(ctx context.Context, interval dbconnector.SpatialTemporalInterval) (posts []dbconnector.Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Select", "Interval", interval, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Select(ctx, interval)
}

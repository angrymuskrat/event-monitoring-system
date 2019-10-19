package middleware

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next dbservice.Service) dbservice.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   dbservice.Service
	logger log.Logger
}

func (mw loggingMiddleware) Push(ctx context.Context, posts []pb.Post) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Push", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Push(ctx, posts)
}

func (mw loggingMiddleware) Select(ctx context.Context, interval pb.SpatialTemporalInterval) (posts []pb.Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Select", "Interval", interval, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Select(ctx, interval)
}

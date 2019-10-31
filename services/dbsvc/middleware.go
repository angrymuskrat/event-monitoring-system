package dbsvc

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Push(ctx context.Context, posts []data.Post) (ids []string, err error) {
	defer func() {
		mw.logger.Log("method", "Push", "err", err)
	}()
	return mw.next.Push(ctx, posts)
}

func (mw loggingMiddleware) Select(ctx context.Context, interval data.SpatioTemporalInterval) (posts []data.Post, err error) {
	defer func() {
		mw.logger.Log("method", "Select", "interval", interval.String(), "err", err)
	}()
	return mw.next.Select(ctx, interval)
}

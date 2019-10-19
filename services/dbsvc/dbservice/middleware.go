package dbservice

import (
	"context"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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

func (mw loggingMiddleware) Push(ctx context.Context, posts []pb.Post) (err error) {
	defer func() {
		mw.logger.Log("method", "Push", "err", err)
	}()
	return mw.next.Push(ctx, posts)
}

func (mw loggingMiddleware) Select(ctx context.Context, interval pb.SpatialTemporalInterval) (posts []pb.Post, err error) {
	defer func() {
		mw.logger.Log("method", "Select", "interval", interval.ForLog(), "err", err)
	}()
	return mw.next.Select(ctx, interval)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(pushed, selected metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			pushed:   pushed,
			selected: selected,
			next:     next,
		}
	}
}

type instrumentingMiddleware struct {
	pushed   metrics.Counter
	selected metrics.Counter
	next     Service
}

func (mw instrumentingMiddleware) Push(ctx context.Context, posts []pb.Post) (err error) {
	err = mw.next.Push(ctx, posts)
	mw.pushed.Add(float64(len(posts)))
	return err
}

func (mw instrumentingMiddleware) Select(ctx context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error) {
	posts, err := mw.next.Select(ctx, interval)
	mw.selected.Add(float64(len(posts)))
	return posts, err
}

package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   Service
}

func (mw loggingMiddleware) Push(ctx context.Context, posts []data.Post) (ids []int32, err error) {
	func(begin time.Time) {
		mw.logger.Info("push request",
			zap.String("len(posts)", strconv.Itoa(len(posts))),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	ids, err = mw.next.Push(ctx, posts)
	return
}

func (mw loggingMiddleware) Select(ctx context.Context, interval data.SpatioTemporalInterval) (posts []data.Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("session status",
			zap.Any("interval", interval),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.Select(ctx, interval)
	return
}

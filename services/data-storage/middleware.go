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

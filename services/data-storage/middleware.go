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

func (mw loggingMiddleware) SelectAggrPosts(ctx context.Context, intreval data.SpatioHourInterval) (posts []data.AggregatedPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("select aggregated posts",
			zap.Any("interval", intreval),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, err = mw.next.SelectAggrPosts(ctx, intreval)
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

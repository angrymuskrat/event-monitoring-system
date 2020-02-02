package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"go.uber.org/zap"
	"time"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   CrawlerService
}

func (mw *loggingMiddleware) New(p crawler.Parameters) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("new session",
			zap.Any("params", p),
			zap.String("id", id),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	id, err = mw.next.New(p)
	return
}

func (mw *loggingMiddleware) Status(id string) (s crawler.OutStatus, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("session status",
			zap.String("id", id),
			zap.Any("status", s),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	s, err = mw.next.Status(id)
	return
}

func (mw *loggingMiddleware) Stop(id string) (ok bool, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("session stop",
			zap.String("id", id),
			zap.Bool("stopped", ok),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	ok, err = mw.next.Stop(id)
	return
}

func (mw *loggingMiddleware) Entities(id string) (data []data.Entity, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("session stop",
			zap.String("id", id),
			zap.Int("entities count", len(data)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	data, err = mw.next.Entities(id)
	return
}

func (mw *loggingMiddleware) Posts(id, cursor string, num int) (posts []data.Post, newCur string, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("session stop",
			zap.String("id", id),
			zap.Int("posts count", len(posts)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	posts, newCur, err = mw.next.Posts(id, cursor, num)
	return
}

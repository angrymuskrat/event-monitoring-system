package service

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"go.uber.org/zap"
	"time"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   BackendService
}

func (mw *loggingMiddleware) HeatmapPosts(req HeatmapRequest) (res []data.AggregatedPost, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("heatmap request",
			zap.Any("request", req),
			zap.Int("count", len(res)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	res, err = mw.next.HeatmapPosts(req)
	return
}

func (mw *loggingMiddleware) Timeline(req TimelineRequest) (res Timeline, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("timeline request",
			zap.Any("request", req),
			zap.Int("count", len(res)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	res, err = mw.next.Timeline(req)
	return
}

func (mw *loggingMiddleware) Events(req EventsRequest) (res []data.Event, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("events request",
			zap.Any("request", req),
			zap.Int("count", len(res)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	res, err = mw.next.Events(req)
	return
}

func (mw *loggingMiddleware) SearchEvents(req SearchRequest) (res []data.Event, err error) {
	defer func(begin time.Time) {
		mw.logger.Info("events search",
			zap.Any("request", req),
			zap.Int("count", len(res)),
			zap.Error(err),
			zap.String("took", time.Since(begin).String()))
	}(time.Now())
	res, err = mw.next.SearchEvents(req)
	return
}

package service

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type BackendService interface {
	HeatmapPosts(req HeatmapRequest) ([]data.AggregatedPost, error)
	Timeline(req TimelineRequest) (Timeline, error)
	Events(req EventsRequest) ([]data.Event, error)
	SearchEvents(req SearchRequest) ([]data.Event, error)
}

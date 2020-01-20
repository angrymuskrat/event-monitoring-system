package service

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type StorageConnector interface {
	HeatmapPosts(city string, topLeft, botRight data.Point, hour int64) ([]data.AggregatedPost, error)
	Timeline(city string, start, finish int64) (Timeline, error)
	Events(city string, topLeft, botRight data.Point, hour int64) ([]data.Event, error)
	EventsByTags(city string, keytags []string, start, finish int64) ([]data.Event, error)
}

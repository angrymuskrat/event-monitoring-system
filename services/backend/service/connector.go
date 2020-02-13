package service

import (
	"errors"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

type StorageConnector interface {
	HeatmapPosts(city string, topLeft, botRight data.Point, hour int64) ([]data.AggregatedPost, error)
	Timeline(city string, start, finish int64) (Timeline, error)
	Events(city string, topLeft, botRight data.Point, hour int64) ([]data.Event, error)
	EventsByTags(city string, keytags []string, start, finish int64) ([]data.Event, error)
}

func setConnector(cType string, params map[string]string) (StorageConnector, error) {
	switch cType {
	case "mock":
		return MockConnector{}, nil
	case "storage":
		dsAddr, ok := params["address"]
		if !ok {
			return nil, errors.New("unable to get data storage address")
		}
		return NewDataConnector(dsAddr)
	default:
		return nil, errors.New("unknown connector type")
	}
}

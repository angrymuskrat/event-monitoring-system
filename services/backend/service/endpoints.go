package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeHeatmapEndpoint(svc BackendService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(HeatmapRequest)
		res, err := svc.HeatmapPosts(p)
		return res, err
	}
}

func makeTimelineEndpoint(svc BackendService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(TimelineRequest)
		res, err := svc.Timeline(p)
		return res, err
	}
}

func makeEventsEndpoint(svc BackendService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(EventsRequest)
		res, err := svc.Events(p)
		return res, err
	}
}

func makeEventsSearchEndpoint(svc BackendService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(SearchRequest)
		res, err := svc.SearchEvents(p)
		return res, err
	}
}

package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/go-kit/kit/endpoint"
)

func makeNewEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(crawler.Parameters)
		id, err := svc.New(p)
		if err != nil {
			return NewEpResponse{"", err.Error()}, nil
		}
		return NewEpResponse{id, ""}, nil
	}
}

func makeStatusEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(IDEpRequest)
		s, err := svc.Status(req.ID)
		if err != nil {
			return StatusEpResponse{crawler.OutStatus{}, err.Error()}, nil
		}
		return StatusEpResponse{s, ""}, nil
	}
}

func makeStopEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(IDEpRequest)
		ok, err := svc.Stop(req.ID)
		if err != nil {
			return StopEpResponse{false, err.Error()}, nil
		}
		return StopEpResponse{ok, ""}, nil
	}
}

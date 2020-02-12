package service

import (
	"context"

	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/go-kit/kit/endpoint"
)

func makeHistoricGridsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.HistoricRequest)
		id, err := s.HistoricGrids(ctx, req)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.HistoricResponse{Id: id, Err: msg}, nil
	}
}

func makeHistoricStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.StatusRequest)
		s, f, err := s.HistoricStatus(ctx, req)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.StatusResponse{Status: s, Finished: f, Err: msg}, nil
	}
}

func makeFindEventsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.EventRequest)
		id, err := s.FindEvents(ctx, req)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.EventResponse{Id: id, Err: msg}, nil
	}
}

func makeEventsStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.StatusRequest)
		s, f, err := s.EventsStatus(ctx, req)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.StatusResponse{Status: s, Finished: f, Err: msg}, nil
	}
}

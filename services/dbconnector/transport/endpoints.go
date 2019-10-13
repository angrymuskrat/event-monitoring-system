package transport

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/dbconnector"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Push   endpoint.Endpoint
	Select endpoint.Endpoint
}

func MakeEndpoints(s dbconnector.Service) Endpoints {
	return Endpoints{
		Push:   makePushEndpoint(s),
		Select: makeSelectEndpoint(s),
	}
}

func makePushEndpoint(s dbconnector.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PushRequest) // type assertion
		err := s.Push(ctx, req.Posts)
		return PushResponse{ Err: err }, nil
	}
}

func makeSelectEndpoint(s dbconnector.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SelectRequest)
		posts, err := s.Select(ctx, req.Interval)
		return SelectResponse{Posts: posts, Err: err}, nil
	}
}

package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"

	"github.com/go-kit/kit/endpoint"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	PushPostsEndpoint   endpoint.Endpoint
	SelectPostsEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func NewEndpoint(svc Service) Set {
	var pushEndpoint = makePushPostsEndpoint(svc)
	var selectEndpoint = makeSelectPostsEndpoint(svc)
	return Set{
		PushPostsEndpoint:   pushEndpoint,
		SelectPostsEndpoint: selectEndpoint,
	}
}

// Push implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) PushPosts(ctx context.Context, posts []data.Post) ([]int32, error) {
	resp, err := s.PushPostsEndpoint(ctx, proto.PushPostsRequest{Posts: posts})
	if err != nil {
		return nil, err
	}
	response := resp.(PushPostsResponse)
	return response.Ids, response.Err
}

// Select implements the service interface, so Set may be used as a
// service. This is primarily useful in the context of a client library.
func (s Set) SelectPosts(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	resp, err := s.SelectPostsEndpoint(ctx, proto.SelectPostsRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(SelectPostsResponse)
	return response.Posts, response.Err
}

// makePushEndpoint constructs a Push endpoint wrapping the service.
func makePushPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushPostsRequest)
		ids, err := s.PushPosts(ctx, req.Posts)
		return PushPostsResponse{Ids: ids, Err: err}, nil
	}
}

// makeSelectEndpoint constructs a Select endpoint wrapping the service.
func makeSelectPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectPostsRequest)
		posts, err := s.SelectPosts(ctx, req.Interval)
		return SelectPostsResponse{Posts: posts, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = PushPostsResponse{}
	_ endpoint.Failer = SelectPostsResponse{}
)

// PushResponse collects the response values for the Push method.
type PushPostsResponse struct {
	Ids []int32 `json:"ids"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r PushPostsResponse) Failed() error { return r.Err }

// SelectResponse collects the response values for the Select method.
type SelectPostsResponse struct {
	Posts []data.Post `json:"posts"`
	Err   error       `json:"-"`
}

// Failed implements endpoint.Failer.
func (r SelectPostsResponse) Failed() error { return r.Err }

package dbsvc

import (
	"context"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/proto"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"

	//"time"

	//"golang.org/x/time/rate"

	"github.com/go-kit/kit/endpoint"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//stdzipkin "github.com/openzipkin/zipkin-go"
	//"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	PushEndpoint   endpoint.Endpoint
	SelectEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func NewEndpoint(svc Service) Set {
	var pushEndpoint = makePushEndpoint(svc)
	var selectEndpoint = makeSelectEndpoint(svc)
	return Set{
		PushEndpoint:   pushEndpoint,
		SelectEndpoint: selectEndpoint,
	}
}

// Push implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Push(ctx context.Context, posts []data.Post) ([]string, error) {
	resp, err := s.PushEndpoint(ctx, proto.PushRequest{Posts: posts})
	fmt.Print(resp, " resp!!!\n")
	if err != nil {
		return nil, err
	}
	response := resp.(PushResponse)
	return response.Ids, response.Err
}

// Select implements the service interface, so Set may be used as a
// service. This is primarily useful in the context of a client library.
func (s Set) Select(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	resp, err := s.SelectEndpoint(ctx, proto.SelectRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(SelectResponse)
	return response.Posts, response.Err
}

// makePushEndpoint constructs a Push endpoint wrapping the service.
func makePushEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushRequest)
		ids, err := s.Push(ctx, req.Posts)
		return PushResponse{Ids: ids, Err: err}, nil
	}
}

// makeSelectEndpoint constructs a Select endpoint wrapping the service.
func makeSelectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectRequest)
		posts, err := s.Select(ctx, req.Interval)
		return SelectResponse{Posts: posts, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = PushResponse{}
	_ endpoint.Failer = SelectResponse{}
)

// PushResponse collects the response values for the Push method.
type PushResponse struct {
	Ids []string `json:"ids"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r PushResponse) Failed() error { return r.Err }

// SelectResponse collects the response values for the Select method.
type SelectResponse struct {
	Posts []data.Post `json:"posts"`
	Err   error       `json:"-"`
}

// Failed implements endpoint.Failer.
func (r SelectResponse) Failed() error { return r.Err }

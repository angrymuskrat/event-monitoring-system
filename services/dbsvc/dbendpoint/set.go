package dbendpoint

import (
	"context"
	//"time"

	//"golang.org/x/time/rate"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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
func New(svc dbservice.Service, logger log.Logger, duration metrics.Histogram/*, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer*/) Set {
	var pushEndpoint endpoint.Endpoint
	{
		pushEndpoint = MakePushEndpoint(svc)
		//pushEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(pushEndpoint)
		/*pushEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(pushEndpoint)
		pushEndpoint = opentracing.TraceServer(otTracer, "Push")(pushEndpoint)
		if zipkinTracer != nil {
			pushEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Push")(pushEndpoint)
		}*/
		pushEndpoint = LoggingMiddleware(log.With(logger, "method", "Push"))(pushEndpoint)
		pushEndpoint = InstrumentingMiddleware(duration.With("method", "Push"))(pushEndpoint)
	}
	var selectEndpoint endpoint.Endpoint
	{
		selectEndpoint = MakeSelectEndpoint(svc)
		//selectEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(selectEndpoint)
		/*selectEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(selectEndpoint)
		selectEndpoint = opentracing.TraceServer(otTracer, "Select")(selectEndpoint)
		if zipkinTracer != nil {
			selectEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Se;ect")(selectEndpoint)
		}*/
		selectEndpoint = LoggingMiddleware(log.With(logger, "method", "Select"))(selectEndpoint)
		selectEndpoint = InstrumentingMiddleware(duration.With("method", "Select"))(selectEndpoint)
	}
	return Set{
		PushEndpoint:    pushEndpoint,
		SelectEndpoint: selectEndpoint,
	}
}

// Push implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Push(ctx context.Context, posts []pb.Post) error {
	resp, err := s.PushEndpoint(ctx, PushRequest{Posts: posts})
	if err != nil {
		return err
	}
	response := resp.(PushResponse)
	return response.Err
}

// Select implements the service interface, so Set may be used as a
// service. This is primarily useful in the context of a client library.
func (s Set) Select(ctx context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error) {
	resp, err := s.SelectEndpoint(ctx, SelectRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(SelectResponse)
	return response.Posts, response.Err
}


// MakePushEndpoint constructs a Push endpoint wrapping the service.
func MakePushEndpoint(s dbservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PushRequest)
		err = s.Push(ctx, req.Posts)
		return PushResponse{Err: err}, nil
	}
}

// MakeSelectEndpoint constructs a Select endpoint wrapping the service.
func MakeSelectEndpoint(s dbservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SelectRequest)
		posts, err := s.Select(ctx, req.Interval)
		return SelectResponse{Posts:posts, Err: err}, nil
	}
}


// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = PushResponse{}
	_ endpoint.Failer = SelectResponse{}
)

// PushRequest collects the request parameters for the Push method.
type PushRequest struct {
	Posts []pb.Post
}

// PushResponse collects the response values for the Push method.
type PushResponse struct {
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r PushResponse) Failed() error { return r.Err }

// SelectRequest collects the request parameters for the Select method.
type SelectRequest struct {
	Interval pb.SpatialTemporalInterval
}

// SelectResponse collects the response values for the Select method.
type SelectResponse struct {
	Posts []pb.Post `json:"posts"`
	Err   error     `json:"-"`
}

// Failed implements endpoint.Failer.
func (r SelectResponse) Failed() error { return r.Err }

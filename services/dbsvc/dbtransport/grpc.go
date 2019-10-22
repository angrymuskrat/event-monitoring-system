package dbtransport

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbendpoint"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
	"github.com/go-kit/kit/circuitbreaker"
	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//stdzipkin "github.com/openzipkin/zipkin-go"
	//"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	push     grpctransport.Handler
	mySelect grpctransport.Handler //select isn't be name of field
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints dbendpoint.Set/*, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer*/, logger log.Logger) pb.DbsvcServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	/*if zipkinTracer != nil {
		// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a
		// ServerOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		//
		// In this example, we demonstrate a global Zipkin tracing service with
		// Go kit gRPC Interceptor.
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}*/

	return &grpcServer{
		push: grpctransport.NewServer(
			endpoints.PushEndpoint,
			decodeGRPCPushRequest,
			encodeGRPCPushResponse,
			append(options/*, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Push", logger))*/)...,
		),
		mySelect: grpctransport.NewServer(
			endpoints.SelectEndpoint,
			decodeGRPCSelectRequest,
			encodeGRPCSelectResponse,
			append(options/*, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Select", logger))*/)...,
		),
	}
}

func (s *grpcServer) Push(ctx context.Context, req *pb.PushRequest) (*pb.PushReply, error) {
	_, rep, err := s.push.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.PushReply), nil
}

func (s *grpcServer) Select(ctx context.Context, req *pb.SelectRequest) (*pb.SelectReply, error) {
	_, rep, err := s.mySelect.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SelectReply), nil
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn/*, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer*/, logger log.Logger) dbservice.Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	//limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []grpctransport.ClientOption


	/*if zipkinTracer != nil {
		// Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing client can be instantiated
		// without an operation name and fed to each Go kit client as ClientOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path.
		//
		// In this example, we demonstrace a global tracing client.
		options = append(options, zipkin.GRPCClientTrace(zipkinTracer))

	}*/

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var pushEndpoint endpoint.Endpoint
	{
		pushEndpoint = grpctransport.NewClient(
			conn,
			"pb.Dbsvc",
			"Push",
			encodeGRPCPushRequest,
			decodeGRPCPushResponse,
			pb.PushReply{},
			append(options/*, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger))*/)...,
		).Endpoint()
		//pushEndpoint = opentracing.TraceClient(otTracer, "Push")(pushEndpoint)
		//pushEndpoint = limiter(pushEndpoint)
		pushEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Push",
			Timeout: 30 * time.Second,
		}))(pushEndpoint)
	}

	// The Select endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var selectEndpoint endpoint.Endpoint
	{
		selectEndpoint = grpctransport.NewClient(
			conn,
			"pb.Dbsvc",
			"Select",
			encodeGRPCSelectRequest,
			decodeGRPCSelectResponse,
			pb.SelectReply{},
			append(options/*, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger))*/)...,
		).Endpoint()
		//selectEndpoint = opentracing.TraceClient(otTracer, "Select")(selectEndpoint)
		//selectEndpoint = limiter(selectEndpoint)
		selectEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Select",
			Timeout: 10 * time.Second,
		}))(selectEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return dbendpoint.Set{
		PushEndpoint:   pushEndpoint,
		SelectEndpoint: selectEndpoint,
	}
}

// decodeGRPCPushRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC push request to a user-domain push request. Primarily useful in a server.
func decodeGRPCPushRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PushRequest)
	return dbendpoint.PushRequest{Posts: toPostArray(req.Posts)}, nil
}

// decodeGRPCSelectRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC select request to a user-domain select request. Primarily useful in a
// server.
func decodeGRPCSelectRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SelectRequest)
	return dbendpoint.SelectRequest{Interval: *req.Interval}, nil
}

// decodeGRPCPushResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC push reply to a user-domain push response. Primarily useful in a client.
func decodeGRPCPushResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.PushReply)
	return dbendpoint.PushResponse{Err: str2err(reply.Err)}, nil
}

// decodeGRPCSelectResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC select reply to a user-domain select response. Primarily useful in a
// client.
func decodeGRPCSelectResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SelectReply)
	return dbendpoint.SelectResponse{Posts: toPostArray(reply.Posts), Err: str2err(reply.Err)}, nil
}

// encodeGRPCPushResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain push response to a gRPC push reply. Primarily useful in a server.
func encodeGRPCPushResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(dbendpoint.PushResponse)
	return &pb.PushReply{Err: err2str(resp.Err)}, nil
}

// encodeGRPCSelectResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain select response to a gRPC select reply. Primarily useful in a
// server.
func encodeGRPCSelectResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(dbendpoint.SelectResponse)
	return &pb.SelectReply{Posts: toPointerArray(resp.Posts), Err: err2str(resp.Err)}, nil
}

// encodeGRPCPushRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain push request to a gRPC push request. Primarily useful in a client.
func encodeGRPCPushRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dbendpoint.PushRequest)
	return &pb.PushRequest{Posts: toPointerArray(req.Posts)}, nil
}

// encodeGRPCSelectRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain select request to a gRPC select request. Primarily useful in a
// client.
func encodeGRPCSelectRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dbendpoint.SelectRequest)
	return &pb.SelectRequest{Interval: &req.Interval}, nil
}

func toPointerArray(array []pb.Post) (posts []*pb.Post) {
	for _, post := range array {
		posts = append(posts, &post)
	}
	return posts
}

func toPostArray(array []*pb.Post) (posts []pb.Post) {
	for _, post := range array {
		posts = append(posts, *post)
	}
	return posts
}

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

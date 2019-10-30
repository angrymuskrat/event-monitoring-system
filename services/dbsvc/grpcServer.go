package dbsvc

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/proto"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type grpcServer struct {
	push     grpctransport.Handler
	mySelect grpctransport.Handler //select isn't be name of field
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Set, logger log.Logger) proto.DBsvcServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		push: grpctransport.NewServer(
			endpoints.PushEndpoint,
			decodeGRPCPushRequest,
			encodeGRPCPushResponse,
			append(options /*, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Push", logger))*/)...,
		),
		mySelect: grpctransport.NewServer(
			endpoints.SelectEndpoint,
			decodeGRPCSelectRequest,
			encodeGRPCSelectResponse,
			append(options /*, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Select", logger))*/)...,
		),
	}
}

func (s *grpcServer) Push(ctx context.Context, req *proto.PushRequest) (*proto.PushReply, error) {
	_, rep, err := s.push.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PushReply), nil
}

func (s *grpcServer) Select(ctx context.Context, req *proto.SelectRequest) (*proto.SelectReply, error) {
	_, rep, err := s.mySelect.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.SelectReply), nil
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn /*, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer*/, logger log.Logger) Service {
	var pushEndpoint endpoint.Endpoint
	{

		pushEndpoint = grpctransport.NewClient(
			conn,
			"proto.DBsvc",
			"Push",
			encodeGRPCPushRequest,
			decodeGRPCPushResponse,
			proto.PushReply{},
		).Endpoint()
		pushEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Push",
			Timeout: 30 * time.Second,
		}))(pushEndpoint)
	}

	var selectEndpoint endpoint.Endpoint
	{
		selectEndpoint = grpctransport.NewClient(
			conn,
			"proto.DBsvc",
			"Select",
			encodeGRPCSelectRequest,
			decodeGRPCSelectResponse,
			proto.SelectReply{},
		).Endpoint()
		selectEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Select",
			Timeout: 10 * time.Second,
		}))(selectEndpoint)
	}

	return Set{
		PushEndpoint:   pushEndpoint,
		SelectEndpoint: selectEndpoint,
	}
}

// decodeGRPCPushRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC push request to a user-domain push request. Primarily useful in a server.
func decodeGRPCPushRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushRequest)
	return proto.PushRequest{Posts: req.Posts}, nil
}

// decodeGRPCSelectRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC select request to a user-domain select request. Primarily useful in a
// server.
func decodeGRPCSelectRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectRequest)
	return proto.SelectRequest{Interval: req.Interval}, nil
}

// decodeGRPCPushResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC push reply to a user-domain push response. Primarily useful in a client.
func decodeGRPCPushResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushReply)
	return PushResponse{Err: str2err(reply.Err)}, nil
}

// decodeGRPCSelectResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC select reply to a user-domain select response. Primarily useful in a
// client.
func decodeGRPCSelectResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectReply)
	return SelectResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}

// encodeGRPCPushResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain push response to a gRPC push reply. Primarily useful in a server.
func encodeGRPCPushResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushResponse)
	return &proto.PushReply{Err: err2str(resp.Err)}, nil
}

// encodeGRPCSelectResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain select response to a gRPC select reply. Primarily useful in a
// server.
func encodeGRPCSelectResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectResponse)
	return &proto.SelectReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}

// encodeGRPCPushRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain push request to a gRPC push request. Primarily useful in a client.
func encodeGRPCPushRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushRequest)
	return &proto.PushRequest{Posts: req.Posts}, nil
}

// encodeGRPCSelectRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain select request to a gRPC select request. Primarily useful in a
// client.
func encodeGRPCSelectRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectRequest)
	return &proto.SelectRequest{Interval: req.Interval}, nil
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

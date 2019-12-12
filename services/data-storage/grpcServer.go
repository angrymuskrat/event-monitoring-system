package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
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
	mySelect grpctransport.Handler
}

func NewGRPCServer(endpoints Set, logger log.Logger) proto.DataStorageServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		push: grpctransport.NewServer(
			endpoints.PushEndpoint,
			decodeGRPCPushRequest,
			encodeGRPCPushResponse,
			append(options)...,
		),
		mySelect: grpctransport.NewServer(
			endpoints.SelectEndpoint,
			decodeGRPCSelectRequest,
			encodeGRPCSelectResponse,
			append(options)...,
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

func NewGRPCClient(conn *grpc.ClientConn) Service {
	var pushEndpoint endpoint.Endpoint
	{

		pushEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
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
			"proto.DataStorage",
			"Select",
			encodeGRPCSelectRequest,
			decodeGRPCSelectResponse,
			proto.SelectReply{},
		).Endpoint()
		selectEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Select",
			Timeout: 30 * time.Second,
		}))(selectEndpoint)
	}

	return Set{
		PushEndpoint:   pushEndpoint,
		SelectEndpoint: selectEndpoint,
	}
}

func decodeGRPCPushRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushRequest)
	return proto.PushRequest{Posts: req.Posts}, nil
}

func decodeGRPCSelectRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectRequest)
	return proto.SelectRequest{Interval: req.Interval}, nil
}

func decodeGRPCPushResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushReply)
	return PushResponse{Ids: reply.Ids, Err: str2err(reply.Err)}, nil
}

func decodeGRPCSelectResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectReply)
	return SelectResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}

func encodeGRPCPushResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushResponse)
	return &proto.PushReply{Ids: resp.Ids, Err: err2str(resp.Err)}, nil
}

func encodeGRPCSelectResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectResponse)
	return &proto.SelectReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}

func encodeGRPCPushRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushRequest)
	return &proto.PushRequest{Posts: req.Posts}, nil
}

func encodeGRPCSelectRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectRequest)
	return &proto.SelectRequest{Interval: req.Interval}, nil
}

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

package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pushPosts     grpctransport.Handler
	selectPosts grpctransport.Handler
}

func NewGRPCServer(endpoints Set) proto.DataStorageServer {

	return &grpcServer{
		pushPosts: grpctransport.NewServer(
			endpoints.PushPostsEndpoint,
			decodeGRPCPushPostsRequest,
			encodeGRPCPushPostsResponse,
		),
		selectPosts: grpctransport.NewServer(
			endpoints.SelectPostsEndpoint,
			decodeGRPCSelectPostsRequest,
			encodeGRPCSelectPostsResponse,
		),
	}
}

func (s *grpcServer) PushPosts(ctx context.Context, req *proto.PushPostsRequest) (*proto.PushPostsReply, error) {
	_, rep, err := s.pushPosts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PushPostsReply), nil
}

func (s *grpcServer) SelectPosts(ctx context.Context, req *proto.SelectPostsRequest) (*proto.SelectPostsReply, error) {
	_, rep, err := s.selectPosts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.SelectPostsReply), nil
}

func NewGRPCClient(conn *grpc.ClientConn) Service {
	var pushPostsEndpoint endpoint.Endpoint
	{

		pushPostsEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"PushPosts",
			encodeGRPCPushPostsRequest,
			decodeGRPCPushPostsResponse,
			proto.PushPostsReply{},
		).Endpoint()
		pushPostsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PushPosts",
			Timeout: 30 * time.Second,
		}))(pushPostsEndpoint)
	}

	var selectPostsEndpoint endpoint.Endpoint
	{
		selectPostsEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"Select",
			encodeGRPCSelectPostsRequest,
			decodeGRPCSelectPostsResponse,
			proto.SelectPostsReply{},
		).Endpoint()
		selectPostsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Select",
			Timeout: 30 * time.Second,
		}))(selectPostsEndpoint)
	}

	return Set{
		PushPostsEndpoint:   pushPostsEndpoint,
		SelectPostsEndpoint: selectPostsEndpoint,
	}
}

func decodeGRPCPushPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushPostsRequest)
	return proto.PushPostsRequest{Posts: req.Posts}, nil
}

func decodeGRPCSelectPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectPostsRequest)
	return proto.SelectPostsRequest{Interval: req.Interval}, nil
}

func decodeGRPCPushPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushPostsReply)
	return PushPostsResponse{Ids: reply.Ids, Err: str2err(reply.Err)}, nil
}

func decodeGRPCSelectPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectPostsReply)
	return SelectPostsResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}

func encodeGRPCPushPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushPostsResponse)
	return &proto.PushPostsReply{Ids: resp.Ids, Err: err2str(resp.Err)}, nil
}

func encodeGRPCSelectPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectPostsResponse)
	return &proto.SelectPostsReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}

func encodeGRPCPushPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushPostsRequest)
	return &proto.PushPostsRequest{Posts: req.Posts}, nil
}

func encodeGRPCSelectPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectPostsRequest)
	return &proto.SelectPostsRequest{Interval: req.Interval}, nil
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

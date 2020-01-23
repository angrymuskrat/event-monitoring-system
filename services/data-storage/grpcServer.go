package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pushPosts       grpctransport.Handler
	selectPosts     grpctransport.Handler
	selectAggrPosts grpctransport.Handler
	pushGrid        grpctransport.Handler
	pullGrid        grpctransport.Handler
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
		selectAggrPosts: grpctransport.NewServer(
			endpoints.SelectAggrPostsEndpoint,
			decodeGRPCSelectAggrPostsRequest,
			encodeGRPCSelectAggrPostsResponse,
		),
		pushGrid: grpctransport.NewServer(
			endpoints.PushGridEndpoint,
			decodeGRPCPushGridRequest,
			encodeGRPCPushGridResponse,
		),
		pullGrid: grpctransport.NewServer(
			endpoints.PullGridEndpoint,
			decodeGRPCPullGridRequest,
			encodeGRPCPullGridResponse,
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

func (s *grpcServer) SelectAggrPosts(ctx context.Context, req *proto.SelectAggrPostsRequest) (*proto.SelectAggrPostsReply, error) {
	_, rep, err := s.selectAggrPosts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.SelectAggrPostsReply), nil
}

func (s *grpcServer) PushGrid(ctx context.Context, req *proto.PushGridRequest) (*proto.PushGridReply, error) {
	_, rep, err := s.pushGrid.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PushGridReply), nil
}

func (s *grpcServer) PullGrid(ctx context.Context, req *proto.PullGridRequest) (*proto.PullGridReply, error) {
	_, rep, err := s.pullGrid.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PullGridReply), nil
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
			Timeout: TimeWaitingClient,
		}))(pushPostsEndpoint)
	}

	var selectPostsEndpoint endpoint.Endpoint
	{
		selectPostsEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"SelectPosts",
			encodeGRPCSelectPostsRequest,
			decodeGRPCSelectPostsResponse,
			proto.SelectPostsReply{},
		).Endpoint()
		selectPostsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "SelectPosts",
			Timeout: TimeWaitingClient,
		}))(selectPostsEndpoint)
	}

	var selectAggrPostsEndpoint endpoint.Endpoint
	{
		selectAggrPostsEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"SelectAggrPosts",
			encodeGRPCSelectAggrPostsRequest,
			decodeGRPCSelectAggrPostsResponse,
			proto.SelectAggrPostsReply{},
		).Endpoint()
		selectAggrPostsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "SelectAggrPosts",
			Timeout: TimeWaitingClient,
		}))(selectAggrPostsEndpoint)
	}

	var pushGridEndpoint endpoint.Endpoint
	{
		pushGridEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"PushGrid",
			encodeGRPCPushGridRequest,
			decodeGRPCPushGridResponse,
			proto.PushGridReply{},
		).Endpoint()
		pushGridEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PushGrid",
			Timeout: TimeWaitingClient,
		}))(pushGridEndpoint)
	}

	var pullGridEndpoint endpoint.Endpoint
	{
		pullGridEndpoint = grpctransport.NewClient(
			conn,
			"proto.DataStorage",
			"PullGrid",
			encodeGRPCPullGridRequest,
			decodeGRPCPullGridResponse,
			proto.PullGridReply{},
		).Endpoint()
		pullGridEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PullGrid",
			Timeout: TimeWaitingClient,
		}))(pullGridEndpoint)
	}

	return Set{
		PushPostsEndpoint:       pushPostsEndpoint,
		SelectPostsEndpoint:     selectPostsEndpoint,
		SelectAggrPostsEndpoint: selectAggrPostsEndpoint,
		PushGridEndpoint:        pushGridEndpoint,
		PullGridEndpoint:        pullGridEndpoint,
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

func decodeGRPCSelectAggrPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectAggrPostsRequest)
	return proto.SelectAggrPostsRequest{Interval: req.Interval}, nil
}

func decodeGRPCPushGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushGridRequest)
	return proto.PushGridRequest{Id: req.Id, Blob: req.Blob}, nil
}

func decodeGRPCPullGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullGridRequest)
	return proto.PullGridRequest{Id: req.Id}, nil
}

func decodeGRPCPushPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushPostsReply)
	return PushPostsResponse{Ids: reply.Ids, Err: str2err(reply.Err)}, nil
}

func decodeGRPCSelectPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectPostsReply)
	return SelectPostsResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}

func decodeGRPCSelectAggrPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectAggrPostsReply)
	return SelectAggrPostsResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}

func decodeGRPCPushGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushGridReply)
	return PushGridResponse{Err: str2err(reply.Err)}, nil
}

func decodeGRPCPullGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullGridReply)
	return PullGridResponse{Blob: reply.Blob, Err: str2err(reply.Err)}, nil
}

func encodeGRPCPushPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushPostsResponse)
	return &proto.PushPostsReply{Ids: resp.Ids, Err: err2str(resp.Err)}, nil
}

func encodeGRPCSelectPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectPostsResponse)
	return &proto.SelectPostsReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}

func encodeGRPCSelectAggrPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectAggrPostsResponse)
	return &proto.SelectAggrPostsReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}

func encodeGRPCPushGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushGridResponse)
	return &proto.PushGridReply{Err: err2str(resp.Err)}, nil
}

func encodeGRPCPullGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PullGridResponse)
	return &proto.PullGridReply{Blob: resp.Blob, Err: err2str(resp.Err)}, nil
}

func encodeGRPCPushPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushPostsRequest)
	return &proto.PushPostsRequest{Posts: req.Posts}, nil
}

func encodeGRPCSelectPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectPostsRequest)
	return &proto.SelectPostsRequest{Interval: req.Interval}, nil
}

func encodeGRPCSelectAggrPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectAggrPostsRequest)
	return &proto.SelectAggrPostsRequest{Interval: req.Interval}, nil
}

func encodeGRPCPushGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushGridRequest)
	return &proto.PushGridRequest{Id: req.Id, Blob: req.Blob}, nil
}

func encodeGRPCPullGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullGridRequest)
	return &proto.PullGridRequest{Id: req.Id}, nil
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

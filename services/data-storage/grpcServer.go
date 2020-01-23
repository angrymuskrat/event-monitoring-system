package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pushPosts       grpctransport.Handler
	selectPosts     grpctransport.Handler
	selectAggrPosts grpctransport.Handler
	pushGrid        grpctransport.Handler
	pullGrid        grpctransport.Handler
	pushEvents      grpctransport.Handler
	pullEvents      grpctransport.Handler
	pushLocations   grpctransport.Handler
	pullLocations   grpctransport.Handler
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
		pushEvents: grpctransport.NewServer(
			endpoints.PushEventsEndpoint,
			decodeGRPCPushEventsRequest,
			encodeGRPCPushEventsResponse,
		),
		pullEvents: grpctransport.NewServer(
			endpoints.PullEventsEndpoint,
			decodeGRPCPullEventsRequest,
			encodeGRPCPullEventsResponse,
		),
		pushLocations: grpctransport.NewServer(
			endpoints.PushLocationsEndpoint,
			decodeGRPCPushLocationsRequest,
			encodeGRPCPushLocationsResponse,
		),
		pullLocations: grpctransport.NewServer(
			endpoints.PullLocationsEndpoint,
			decodeGRPCPullLocationsRequest,
			encodeGRPCPullLocationsResponse,
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

func (s *grpcServer) PushEvents(ctx context.Context, req *proto.PushEventsRequest) (*proto.PushEventsReply, error) {
	_, rep, err := s.pushEvents.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PushEventsReply), nil
}

func (s *grpcServer) PullEvents(ctx context.Context, req *proto.PullEventsRequest) (*proto.PullEventsReply, error) {
	_, rep, err := s.pullEvents.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PullEventsReply), nil
}

func (s *grpcServer) PushLocations(ctx context.Context, req *proto.PushLocationsRequest) (*proto.PushLocationsReply, error) {
	_, rep, err := s.pushLocations.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PushLocationsReply), nil
}

func (s *grpcServer) PullLocations(ctx context.Context, req *proto.PullLocationsRequest) (*proto.PullLocationsReply, error) {
	_, rep, err := s.pullLocations.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PullLocationsReply), nil
}


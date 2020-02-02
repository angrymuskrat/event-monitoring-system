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
	pullTimeline    grpctransport.Handler
	pushGrid        grpctransport.Handler
	pullGrid        grpctransport.Handler
	pushEvents      grpctransport.Handler
	pullEvents      grpctransport.Handler
	pushLocations   grpctransport.Handler
	pullLocations   grpctransport.Handler
}

func NewGRPCServer(svc Service) proto.DataStorageServer {
	return &grpcServer{
		pushPosts: grpctransport.NewServer(
			makePushPostsEndpoint(svc),
			decodeGRPCPushPostsRequest,
			encodeGRPCPushPostsResponse,
		),
		selectPosts: grpctransport.NewServer(
			makeSelectPostsEndpoint(svc),
			decodeGRPCSelectPostsRequest,
			encodeGRPCSelectPostsResponse,
		),
		selectAggrPosts: grpctransport.NewServer(
			makeSelectAggrPostsEndpoint(svc),
			decodeGRPCSelectAggrPostsRequest,
			encodeGRPCSelectAggrPostsResponse,
		),
		pullTimeline: grpctransport.NewServer(
			makePullTimelineEndpoint(svc),
			decodeGRPCPullTimelineRequest,
			encodeGRPCPullTimelineResponse,
		),
		pushGrid: grpctransport.NewServer(
			makePushGridEndpoint(svc),
			decodeGRPCPushGridRequest,
			encodeGRPCPushGridResponse,
		),
		pullGrid: grpctransport.NewServer(
			makePullGridEndpoint(svc),
			decodeGRPCPullGridRequest,
			encodeGRPCPullGridResponse,
		),
		pushEvents: grpctransport.NewServer(
			makePushEventsEndpoint(svc),
			decodeGRPCPushEventsRequest,
			encodeGRPCPushEventsResponse,
		),
		pullEvents: grpctransport.NewServer(
			makePullEventsEndpoint(svc),
			decodeGRPCPullEventsRequest,
			encodeGRPCPullEventsResponse,
		),
		pushLocations: grpctransport.NewServer(
			makePushLocationsEndpoint(svc),
			decodeGRPCPushLocationsRequest,
			encodeGRPCPushLocationsResponse,
		),
		pullLocations: grpctransport.NewServer(
			makePullLocationsEndpoint(svc),
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

func (s *grpcServer) PullTimeline(ctx context.Context, req *proto.PullTimelineRequest) (*proto.PullTimelineReply, error) {
	_, rep, err := s.pullTimeline.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.PullTimelineReply), nil
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

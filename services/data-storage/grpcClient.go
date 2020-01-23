package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn) Service {
	var pushPostsEndpoint endpoint.Endpoint
	{
		pushPostsEndpoint = grpctransport.NewClient(
			conn, "proto.DataStorage", "PushPosts",
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
			conn, "proto.DataStorage", "SelectPosts",
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
			conn, "proto.DataStorage", "SelectAggrPosts",
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
			conn, "proto.DataStorage", "PushGrid",
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
			conn, "proto.DataStorage", "PullGrid",
			encodeGRPCPullGridRequest,
			decodeGRPCPullGridResponse,
			proto.PullGridReply{},
		).Endpoint()
		pullGridEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PullGrid",
			Timeout: TimeWaitingClient,
		}))(pullGridEndpoint)
	}

	var pushEventsEndpoint endpoint.Endpoint
	{
		pushEventsEndpoint = grpctransport.NewClient(
			conn, "proto.DataStorage", "PushEvents",
			encodeGRPCPushEventsRequest,
			decodeGRPCPushEventsResponse,
			proto.PushEventsReply{},
		).Endpoint()
		pushEventsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PushEvents",
			Timeout: TimeWaitingClient,
		}))(pushEventsEndpoint)
	}
	var pullEventsEndpoint endpoint.Endpoint
	{
		pullEventsEndpoint = grpctransport.NewClient(
			conn, "proto.DataStorage", "PullEvents",
			encodeGRPCPullEventsRequest,
			decodeGRPCPullEventsResponse,
			proto.PullEventsReply{},
		).Endpoint()
		pullEventsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PullEvents",
			Timeout: TimeWaitingClient,
		}))(pullEventsEndpoint)
	}

	var pushLocationsEndpoint endpoint.Endpoint
	{
		pushLocationsEndpoint = grpctransport.NewClient(
			conn, "proto.DataStorage", "PushLocations",
			encodeGRPCPushLocationsRequest,
			decodeGRPCPushLocationsResponse,
			proto.PushLocationsReply{},
		).Endpoint()
		pushLocationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PushLocations",
			Timeout: TimeWaitingClient,
		}))(pushLocationsEndpoint)
	}
	var pullLocationsEndpoint endpoint.Endpoint
	{
		pullLocationsEndpoint = grpctransport.NewClient(
			conn, "proto.DataStorage", "PullLocations",
			encodeGRPCPullLocationsRequest,
			decodeGRPCPullLocationsResponse,
			proto.PullLocationsReply{},
		).Endpoint()
		pullLocationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "PullLocations",
			Timeout: TimeWaitingClient,
		}))(pullLocationsEndpoint)
	}

	return Set{
		PushPostsEndpoint:       pushPostsEndpoint,
		SelectPostsEndpoint:     selectPostsEndpoint,
		SelectAggrPostsEndpoint: selectAggrPostsEndpoint,
		PushGridEndpoint:        pushGridEndpoint,
		PullGridEndpoint:        pullGridEndpoint,
		PushEventsEndpoint:      pushEventsEndpoint,
		PullEventsEndpoint:      pullEventsEndpoint,
		PushLocationsEndpoint:   pushLocationsEndpoint,
		PullLocationsEndpoint:   pullLocationsEndpoint,
	}
}
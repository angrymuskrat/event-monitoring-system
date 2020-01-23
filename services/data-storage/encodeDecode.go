package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
)

// encode/decode for PushPosts
func decodeGRPCPushPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushPostsRequest)
	return proto.PushPostsRequest{Posts: req.Posts}, nil
}
func decodeGRPCPushPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushPostsReply)
	return PushPostsResponse{Ids: reply.Ids, Err: str2err(reply.Err)}, nil
}
func encodeGRPCPushPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushPostsResponse)
	return &proto.PushPostsReply{Ids: resp.Ids, Err: err2str(resp.Err)}, nil
}
func encodeGRPCPushPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushPostsRequest)
	return &proto.PushPostsRequest{Posts: req.Posts}, nil
}

// encode/decode for SelectPosts
func decodeGRPCSelectPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectPostsRequest)
	return proto.SelectPostsRequest{Interval: req.Interval}, nil
}
func decodeGRPCSelectPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectPostsReply)
	return SelectPostsResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}
func encodeGRPCSelectPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectPostsResponse)
	return &proto.SelectPostsReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}
func encodeGRPCSelectPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectPostsRequest)
	return &proto.SelectPostsRequest{Interval: req.Interval}, nil
}

// encode/decode for SelectAggrPosts
func decodeGRPCSelectAggrPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectAggrPostsRequest)
	return proto.SelectAggrPostsRequest{Interval: req.Interval}, nil
}
func decodeGRPCSelectAggrPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectAggrPostsReply)
	return SelectAggrPostsResponse{Posts: reply.Posts, Err: str2err(reply.Err)}, nil
}
func encodeGRPCSelectAggrPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SelectAggrPostsResponse)
	return &proto.SelectAggrPostsReply{Posts: resp.Posts, Err: err2str(resp.Err)}, nil
}
func encodeGRPCSelectAggrPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectAggrPostsRequest)
	return &proto.SelectAggrPostsRequest{Interval: req.Interval}, nil
}

// encode/decode for PushGrid
func decodeGRPCPushGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushGridRequest)
	return proto.PushGridRequest{Id: req.Id, Blob: req.Blob}, nil
}
func decodeGRPCPushGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushGridReply)
	return PushGridResponse{Err: str2err(reply.Err)}, nil
}
func encodeGRPCPushGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushGridResponse)
	return &proto.PushGridReply{Err: err2str(resp.Err)}, nil
}
func encodeGRPCPushGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushGridRequest)
	return &proto.PushGridRequest{Id: req.Id, Blob: req.Blob}, nil
}

// encode/decode for PullGrid
func decodeGRPCPullGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullGridRequest)
	return proto.PullGridRequest{Id: req.Id}, nil
}
func decodeGRPCPullGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullGridReply)
	return PullGridResponse{Blob: reply.Blob, Err: str2err(reply.Err)}, nil
}
func encodeGRPCPullGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PullGridResponse)
	return &proto.PullGridReply{Blob: resp.Blob, Err: err2str(resp.Err)}, nil
}
func encodeGRPCPullGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullGridRequest)
	return &proto.PullGridRequest{Id: req.Id}, nil
}

// encode/decode for PushEvents
func decodeGRPCPushEventsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushEventsRequest)
	return proto.PushEventsRequest{Events: req.Events}, nil
}
func decodeGRPCPushEventsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushEventsReply)
	return PushEventsResponse{Err: str2err(reply.Err)}, nil
}
func encodeGRPCPushEventsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushEventsResponse)
	return &proto.PushEventsReply{Err: err2str(resp.Err)}, nil
}
func encodeGRPCPushEventsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushEventsRequest)
	return &proto.PushEventsRequest{Events: req.Events}, nil
}

// encode/decode for PullEvents
func decodeGRPCPullEventsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullEventsRequest)
	return proto.PullEventsRequest{Interval: req.Interval}, nil
}
func decodeGRPCPullEventsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullEventsReply)
	return PullEventsResponse{Events: reply.Events, Err: str2err(reply.Err)}, nil
}
func encodeGRPCPullEventsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PullEventsResponse)
	return &proto.PullEventsReply{Events: resp.Events, Err: err2str(resp.Err)}, nil
}
func encodeGRPCPullEventsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullEventsRequest)
	return &proto.PullEventsRequest{Interval: req.Interval}, nil
}

// encode/decode for PushLocations
func decodeGRPCPushLocationsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushLocationsRequest)
	return proto.PushLocationsRequest{City: req.City, Locations: req.Locations}, nil
}
func decodeGRPCPushLocationsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushLocationsReply)
	return PushLocationsResponse{Err: str2err(reply.Err)}, nil
}
func encodeGRPCPushLocationsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PushLocationsResponse)
	return &proto.PushLocationsReply{Err: err2str(resp.Err)}, nil
}
func encodeGRPCPushLocationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushLocationsRequest)
	return &proto.PushLocationsRequest{City: req.City, Locations: req.Locations}, nil
}

// encode/decode for PullLocations
func decodeGRPCPullLocationsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullLocationsRequest)
	return proto.PullLocationsRequest{CityId: req.CityId}, nil
}
func decodeGRPCPullLocationsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullLocationsReply)
	return PullLocationsResponse{Locations: reply.Locations, Err: str2err(reply.Err)}, nil
}
func encodeGRPCPullLocationsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(PullLocationsResponse)
	return &proto.PullLocationsReply{Locations: resp.Locations, Err: err2str(resp.Err)}, nil
}
func encodeGRPCPullLocationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullLocationsRequest)
	return &proto.PullLocationsRequest{CityId: req.CityId}, nil
}

// encode/decode errors
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

package service

import (
	"context"

	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
)

func encodeGRPCInsertCityRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.InsertCityRequest)
	return &req, nil
}

func decodeGRPCInsertCityRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.InsertCityRequest)
	return *req, nil
}

func encodeGRPCInsertCityResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.InsertCityReply)
	return resp, nil
}

func decodeGRPCInsertCityResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.InsertCityReply)
	return *reply, nil
}

func encodeGRPCGetAllCitiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.GetAllCitiesRequest)
	return &req, nil
}

func decodeGRPCGetAllCitiesRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.GetAllCitiesRequest)
	return *req, nil
}

func encodeGRPCGetAllCitiesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.GetAllCitiesReply)
	return resp, nil
}

func decodeGRPCGetAllCitiesResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.GetAllCitiesReply)
	return *reply, nil
}

func encodeGRPCGetCityRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.GetCityRequest)
	return &req, nil
}

func decodeGRPCGetCityRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.GetCityRequest)
	return *req, nil
}

func encodeGRPCGetCityResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.GetCityReply)
	return resp, nil
}

func decodeGRPCGetCityResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.GetCityReply)
	return *reply, nil
}

func encodeGRPCPushPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushPostsRequest)
	return &req, nil
}

func decodeGRPCPushPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushPostsRequest)
	return *req, nil
}

func encodeGRPCPushPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PushPostsReply)
	return resp, nil
}

func decodeGRPCPushPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushPostsReply)
	return *reply, nil
}

func encodeGRPCSelectPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectPostsRequest)
	return &req, nil
}

func decodeGRPCSelectPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectPostsRequest)
	return *req, nil
}

func encodeGRPCSelectPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.SelectPostsReply)
	return &resp, nil
}

func decodeGRPCSelectPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectPostsReply)
	return *reply, nil
}

func encodeGRPCSelectAggrPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.SelectAggrPostsRequest)
	return &req, nil
}

func decodeGRPCSelectAggrPostsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SelectAggrPostsRequest)
	return *req, nil
}

func encodeGRPCSelectAggrPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.SelectAggrPostsReply)
	return &resp, nil
}

func decodeGRPCSelectAggrPostsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SelectAggrPostsReply)
	return *reply, nil
}

func encodeGRPCPullTimelineRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullTimelineRequest)
	return &req, nil
}

func decodeGRPCPullTimelineRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullTimelineRequest)
	return *req, nil
}

func encodeGRPCPullTimelineResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullTimelineReply)
	return &resp, nil
}

func decodeGRPCPullTimelineResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullTimelineReply)
	return *reply, nil
}

func encodeGRPCPushGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushGridRequest)
	return &req, nil
}

func decodeGRPCPushGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushGridRequest)
	return *req, nil
}

func encodeGRPCPushGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PushGridReply)
	return &resp, nil
}

func decodeGRPCPushGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushGridReply)
	return *reply, nil
}

func encodeGRPCPullGridRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullGridRequest)
	return &req, nil
}

func decodeGRPCPullGridRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullGridRequest)
	return *req, nil
}

func encodeGRPCPullGridResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullGridReply)
	return &resp, nil
}

func decodeGRPCPullGridResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullGridReply)
	return *reply, nil
}

func encodeGRPCPushEventsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushEventsRequest)
	return &req, nil
}

func decodeGRPCPushEventsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushEventsRequest)
	return *req, nil
}

func encodeGRPCPushEventsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PushEventsReply)
	return &resp, nil
}

func decodeGRPCPushEventsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushEventsReply)
	return *reply, nil
}

func encodeGRPCPullEventsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullEventsRequest)
	return &req, nil
}

func decodeGRPCPullEventsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullEventsRequest)
	return *req, nil
}

func encodeGRPCPullEventsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullEventsReply)
	return &resp, nil
}

func decodeGRPCPullEventsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullEventsReply)
	return *reply, nil
}

func encodeGRPCPullEventsTagsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullEventsTagsRequest)
	return &req, nil
}

func decodeGRPCPullEventsTagsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullEventsTagsRequest)
	return *req, nil
}

func encodeGRPCPullEventsTagsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullEventsTagsReply)
	return &resp, nil
}

func decodeGRPCPullEventsTagsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullEventsTagsReply)
	return *reply, nil
}

func encodeGRPCPushLocationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PushLocationsRequest)
	return &req, nil
}

func decodeGRPCPushLocationsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PushLocationsRequest)
	return *req, nil
}

func encodeGRPCPushLocationsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PushLocationsReply)
	return &resp, nil
}

func decodeGRPCPushLocationsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PushLocationsReply)
	return *reply, nil
}

func encodeGRPCPullLocationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullLocationsRequest)
	return &req, nil
}

func decodeGRPCPullLocationsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullLocationsRequest)
	return *req, nil
}

func encodeGRPCPullLocationsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullLocationsReply)
	return &resp, nil
}

func decodeGRPCPullLocationsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullLocationsReply)
	return *reply, nil
}

func encodeGRPCPullShortPostInIntervalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullShortPostInIntervalRequest)
	return &req, nil
}

func decodeGRPCPullShortPostInIntervalRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullShortPostInIntervalRequest)
	return *req, nil
}

func encodeGRPCPullShortPostInIntervalResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullShortPostInIntervalReply)
	return &resp, nil
}

func decodeGRPCPullShortPostInIntervalResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullShortPostInIntervalReply)
	return *reply, nil
}

func encodeGRPCPullSingleShortPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.PullSingleShortPostRequest)
	return &req, nil
}

func decodeGRPCPullSingleShortPostRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.PullSingleShortPostRequest)
	return *req, nil
}

func encodeGRPCPullSingleShortPostResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(proto.PullSingleShortPostReply)
	return &resp, nil
}

func decodeGRPCPullSingleShortPostResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.PullSingleShortPostReply)
	return *reply, nil
}

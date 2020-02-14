package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
)

func encodeGRPCHistoricGridsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.HistoricRequest)
	return &req, nil
}

func encodeGRPCHistoricGridsResponse(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.HistoricResponse)
	return &req, nil
}

func encodeGRPCFindEventsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.EventRequest)
	return &req, nil
}

func encodeGRPCFindEventsResponse(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.EventResponse)
	return req, nil
}

func encodeGRPCStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.StatusRequest)
	return &req, nil
}

func encodeGRPCStatusResponse(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(proto.StatusResponse)
	return req, nil
}

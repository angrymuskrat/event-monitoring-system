package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
)

func decodeGRPCHistoricGridsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.HistoricRequest)
	return *req, nil
}

func decodeGRPCHistoricGridsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.HistoricResponse)
	return *reply, nil
}

func decodeGRPCFindEventsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.EventRequest)
	return *req, nil
}

func decodeGRPCFindEventsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.EventResponse)
	return *reply, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.StatusRequest)
	return *req, nil
}

func decodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.StatusResponse)
	return *reply, nil
}

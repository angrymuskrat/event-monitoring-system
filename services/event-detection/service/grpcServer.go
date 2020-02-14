package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"time"
)

const (
	TimeWaitingClient = 30 * time.Second // in seconds
	MaxMsgSize        = 1000000000       // in bytes
)

type server struct {
	historicGrids  grpctransport.Handler
	historicStatus grpctransport.Handler
	findEvents     grpctransport.Handler
	eventsStatus   grpctransport.Handler
}

func Server(svc Service) proto.EventDetectionServer {
	return &server{
		historicGrids: grpctransport.NewServer(
			makeHistoricGridsEndpoint(svc),
			decodeGRPCHistoricGridsRequest,
			encodeGRPCHistoricGridsResponse,
		),
		findEvents: grpctransport.NewServer(
			makeFindEventsEndpoint(svc),
			decodeGRPCFindEventsRequest,
			encodeGRPCFindEventsResponse,
		),
		historicStatus: grpctransport.NewServer(
			makeHistoricStatusEndpoint(svc),
			decodeGRPCStatusRequest,
			encodeGRPCStatusResponse,
		),
		eventsStatus: grpctransport.NewServer(
			makeEventsStatusEndpoint(svc),
			decodeGRPCStatusRequest,
			encodeGRPCStatusResponse,
		),
	}
}

func (gs *server) HistoricGrids(ctx context.Context, histReq *proto.HistoricRequest) (*proto.HistoricResponse, error) {
	_, rep, err := gs.historicGrids.ServeGRPC(ctx, histReq)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.HistoricResponse), nil
}

func (gs *server) HistoricStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	_, rep, err := gs.historicStatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	tmp := rep.(proto.StatusResponse)
	return &tmp, nil
}

func (gs *server) FindEvents(ctx context.Context, eventReq *proto.EventRequest) (*proto.EventResponse, error) {
	_, rep, err := gs.findEvents.ServeGRPC(ctx, eventReq)
	if err != nil {
		return nil, err
	}
	tmp := rep.(proto.EventResponse)
	return &tmp, nil
}

func (gs *server) EventsStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	_, rep, err := gs.eventsStatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	tmp := rep.(proto.StatusResponse)
	return &tmp, nil
}

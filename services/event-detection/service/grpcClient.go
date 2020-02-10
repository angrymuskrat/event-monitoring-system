package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type Client struct {
	HistoricGrids  endpoint.Endpoint
	HistoricStatus endpoint.Endpoint
	FindEvents     endpoint.Endpoint
	EventsStatus   endpoint.Endpoint
}

func NewClient(conn *grpc.ClientConn) Client {
	svc := Client{}

	hisoricGridsEndpoint := grpctransport.NewClient(
		conn, "proto.EventDetection", "HistoricGrids",
		encodeGRPCHistoricGridsRequest,
		decodeGRPCHistoricGridsResponse,
		proto.HistoricResponse{},
	).Endpoint()
	svc.HistoricGrids = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "HistoricGrids",
		Timeout: TimeWaitingClient,
	}))(hisoricGridsEndpoint)

	historicStatusEndpoint := grpctransport.NewClient(
		conn, "proto.EventDetection", "HistoricStatus",
		encodeGRPCStatusRequest,
		decodeGRPCStatusResponse,
		proto.StatusResponse{},
	).Endpoint()
	svc.HistoricStatus = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "HistoricStatus",
		Timeout: TimeWaitingClient,
	}))(historicStatusEndpoint)

	findEventsEndpoint := grpctransport.NewClient(
		conn, "proto.EventDetection", "FindEvents",
		encodeGRPCFindEventsRequest,
		decodeGRPCFindEventsResponse,
		proto.EventResponse{},
	).Endpoint()
	svc.FindEvents = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "FindEvents",
		Timeout: TimeWaitingClient,
	}))(findEventsEndpoint)

	eventsStatusEndpoint := grpctransport.NewClient(
		conn, "proto.EventDetection", "EventsStatus",
		encodeGRPCStatusRequest,
		decodeGRPCStatusResponse,
		proto.StatusResponse{},
	).Endpoint()
	svc.EventsStatus = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "EventsStatus",
		Timeout: TimeWaitingClient,
	}))(eventsStatusEndpoint)

	return svc
}

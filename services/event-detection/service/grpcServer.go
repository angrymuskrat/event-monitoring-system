package service

import (
	"context"
	"net"
	"sync"

	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/google/uuid"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const MaxMsgSize = 1000000000 //maximum grpc message size in bytes

func ServerStart(cfg Config) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		unilog.Logger().Error("unable to create tcp listener", zap.Error(err))
		panic(err)
	}
	srv := grpc.NewServer(grpc.MaxRecvMsgSize(MaxMsgSize))
	gsrv := newGRPCServer(cfg)
	proto.RegisterEventDetectionServer(srv, gsrv)
	err = srv.Serve(listener)
	if err != nil {
		unilog.Logger().Error("error in server execution", zap.Error(err))
		panic(err)
	}
}

type grpcServer struct {
	cfg           Config
	histSesssions []*historicSession
	eventSessions []*eventSession
	mut           sync.Mutex
}

func newGRPCServer(cfg Config) *grpcServer {
	return &grpcServer{cfg: cfg}
}

func (gs *grpcServer) HistoricGrids(ctx context.Context, histReq *proto.HistoricRequest) (*proto.HistoricResponse, error) {
	id := uuid.New().String()
	session := newHistoricSession(gs.cfg, *histReq, id)
	gs.mut.Lock()
	gs.histSesssions = append(gs.histSesssions, session)
	gs.mut.Unlock()
	go session.generateGrids()
	return &proto.HistoricResponse{Id: id, Err: ""}, nil
}

func (gs *grpcServer) HistoricStatus(context.Context, *proto.StatusRequest) (*proto.StatusResponse, error) {
	return nil, nil
}

func (gs *grpcServer) FindEvents(ctx context.Context, eventReq *proto.EventRequest) (*proto.EventResponse, error) {
	id := uuid.New().String()
	session := newEventSession(gs.cfg, *eventReq, id)
	gs.mut.Lock()
	gs.eventSessions = append(gs.eventSessions, session)
	gs.mut.Unlock()
	go session.detectEvents()
	return &proto.EventResponse{Id: id, Err: ""}, nil
}

func (gs *grpcServer) EventsStatus(context.Context, *proto.StatusRequest) (*proto.StatusResponse, error) {
	return nil, nil
}

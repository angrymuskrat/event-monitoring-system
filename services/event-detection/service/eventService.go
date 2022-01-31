package service

import (
	"context"
	"sync"

	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/google/uuid"
)

type eventService struct {
	cfg           Config
	histSessions  map[string]*historicSession
	eventSessions map[string]*eventSession
	mut           sync.Mutex
}

func newEventService(cfg Config) *eventService {
	return &eventService{cfg: cfg, histSessions: make(map[string]*historicSession), eventSessions: make(map[string]*eventSession)}
}

func (svc *eventService) HistoricGrids(_ context.Context, histReq proto.HistoricRequest) (string, error) {
	id := uuid.New().String()
	session := newHistoricSession(svc.cfg, histReq, id)
	svc.mut.Lock()
	svc.histSessions[id] = session
	svc.mut.Unlock()
	go session.generateGrids()
	return id, nil
}

func (svc *eventService) HistoricStatus(_ context.Context, req proto.StatusRequest) (string, bool, error) {
	finished := false
	if svc.histSessions[req.Id].status == FinishedStatus {
		finished = true
	}
	return svc.histSessions[req.Id].status.String(), finished, nil
}

func (svc *eventService) FindEvents(_ context.Context, eventReq proto.EventRequest) (string, error) {
	id := uuid.New().String()
	session := newEventSession(svc.cfg, eventReq, id)
	svc.mut.Lock()
	svc.eventSessions[id] = session
	svc.mut.Unlock()
	go session.detectEvents()
	return id, nil
}

func (svc *eventService) EventsStatus(_ context.Context, req proto.StatusRequest) (string, bool, error) {
	finished := false
	if svc.eventSessions[req.Id].status == FinishedStatus {
		finished = true
	}
	return svc.eventSessions[req.Id].status.String(), finished, nil
}

package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/google/uuid"
	"sync"
)

type eventService struct {
	cfg           Config
	histSesssions []*historicSession
	eventSessions []*eventSession
	mut           sync.Mutex
}

func (svc eventService) HistoricGrids(ctx context.Context, histReq proto.HistoricRequest) (string, error) {
	id := uuid.New().String()
	session := newHistoricSession(svc.cfg, histReq, id)
	svc.mut.Lock()
	svc.histSesssions = append(svc.histSesssions, session)
	svc.mut.Unlock()
	go session.generateGrids()
	return id, nil
}

func (svc eventService) HistoricStatus(ctx context.Context, req proto.StatusRequest) (string, error) {

}

func (svc eventService) FindEvents(ctx context.Context, eventReq proto.EventRequest) (string, error) {
	id := uuid.New().String()
	session := newEventSession(svc.cfg, eventReq, id)
	svc.mut.Lock()
	svc.eventSessions = append(svc.eventSessions, session)
	svc.mut.Unlock()
	go session.detectEvents()
	return id, nil
}

func (svc eventService) EventsStatus(ctx context.Context, req proto.StatusRequest) (string, error) {

}

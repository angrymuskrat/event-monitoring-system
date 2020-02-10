package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
)

type Service interface {
	HistoricGrids(ctx context.Context, histReq proto.HistoricRequest) (string, error)
	HistoricStatus(context.Context, proto.StatusRequest) (string, error)
	FindEvents(ctx context.Context, eventReq proto.EventRequest) (string, error)
	EventsStatus(context.Context, proto.StatusRequest) (string, error)
}

package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/detection"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type eventSession struct {
	id       string
	status   StatusType
	cfg      Config
	eventReq proto.EventRequest
}

func newEventSession(config Config, eventReq proto.EventRequest, id string) *eventSession {
	return &eventSession{
		id:       id,
		status:   RunningStatus,
		cfg:      config,
		eventReq: eventReq,
	}
}

func (es *eventSession) detectEvents() {
	conn, err := grpc.Dial(es.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data storage", zap.Error(err))
		es.status = FailedStatus
		return
	}
	cl := service.NewGRPCClient(conn)
	ids := generateGridIds(es.eventReq.StartTime, es.eventReq.FinishTime)
	rawGrids, err := cl.PullGrid(context.Background(), es.eventReq.CityId, ids)
	if err != nil {
		unilog.Logger().Error("unable to get grids from data storage", zap.Error(err))
		es.status = FailedStatus
		return
	}

	grids := make(map[int64]convtree.ConvTree)
	for n, rawGrid := range rawGrids {
		buf := bytes.NewBuffer(rawGrid)
		dec := gob.NewDecoder(buf)
		var g convtree.ConvTree
		if err := dec.Decode(&g); err != nil {
			unilog.Logger().Error("unable to decode grid", zap.Error(err))
			es.status = FailedStatus
			return
		}
		grids[n] = g
	}

	times, err := getTimes(es.eventReq.StartTime, es.eventReq.FinishTime, es.eventReq.Timezone)
	if err != nil {
		unilog.Logger().Error("unable to generate intervals", zap.Error(err))
		es.status = FailedStatus
		return
	}

	for _, t := range times {
		center := t[0].Add(t[1].Sub(t[0]) / 2)
		timeNum := getGridNum(center.Month(), center.Weekday(), center.Hour())
		grid := grids[timeNum]
		startTime := t[0].Unix()
		finishTime := t[1].Unix()
		posts, _, err := cl.SelectPosts(context.Background(), es.eventReq.CityId, startTime, finishTime)
		if err != nil {
			unilog.Logger().Error("unable to get posts from data storage", zap.Error(err))
			continue
		}

		dsi := data.SpatioHourInterval{
			Hour: startTime - 3600,
			Area: data.Area{
				TopLeft: &data.Point{
					Lat: grid.TopLeft.Y,
					Lon: grid.TopLeft.X,
				},
				BotRight: &data.Point{
					Lat: grid.BottomRight.Y,
					Lon: grid.BottomRight.X,
				},
			},
		}
		oldEvents, err := cl.PullEvents(context.Background(), es.eventReq.CityId, dsi)
		if err != nil {
			unilog.Logger().Error("unable to get events from data storage", zap.Error(err))
			continue
		}

		filterTags := filterTags(es.eventReq.FilterTags)
		evs, found := detection.FindEvents(grid, posts, es.cfg.MaxPoints, filterTags, startTime, finishTime, oldEvents)
		if found {
			unilog.Logger().Info("found events", zap.String("session", es.id),
				zap.Int("num", len(evs)), zap.String("timestamp", t[0].String()))

			err = cl.PushEvents(context.Background(), es.eventReq.CityId, evs)
			if err != nil {
				unilog.Logger().Error("unable to push events to data storage", zap.Error(err))
				return
			}
		} else {
			unilog.Logger().Error("no events found", zap.String("session", es.id), zap.String("timestamp", t[0].String()), zap.Error(err))
		}
	}
	es.status = FinishedStatus
}

func getTimes(start, finish int64, tz string) ([][2]time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		unilog.Logger().Error("unable to load timezone", zap.Error(err))
		return nil, err
	}
	res := [][2]time.Time{}
	s := time.Unix(start, 0)
	s = s.In(loc)
	f := time.Unix(finish, 0)
	f = f.In(loc)
	c := s.Truncate(10 * time.Minute).Add(10 * time.Minute)
	for c.Before(f) {
		r := [2]time.Time{c, c.Add(-1 * time.Hour)}
		res = append(res, r)
		c = c.Add(10 * time.Minute)
	}
	return res, nil
}

func generateGridIds(startDate, finishDate int64) []int64 {
	startTime := time.Unix(startDate, 0)
	finishTime := time.Unix(finishDate, 0)
	t := startTime
	res := []int64{}
	for t.Before(finishTime) {
		v := getGridNum(t.Month(), t.Weekday(), t.Hour())
		res = append(res, v)
		t = t.Add(time.Hour)
	}
	return res
}

func getGridNum(month time.Month, day time.Weekday, hour int) int64 {
	monthNum := int64(month)
	var dayNum int64
	switch day {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		dayNum = 1
	case time.Saturday, time.Sunday:
		dayNum = 2
	}
	gridNum := monthNum*1000 + dayNum*100 + int64(hour)
	return gridNum
}

func filterTags(tags []string) map[string]bool {
	filterTags := map[string]bool{}
	for _, t := range tags {
		filterTags[t] = true
	}
	return filterTags
}

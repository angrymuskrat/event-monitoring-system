package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/detection"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
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

	ids, err := generateGridIds(es.eventReq.StartTime, es.eventReq.FinishTime, es.eventReq.Timezone)
	if err != nil {
		unilog.Logger().Error("unable to generate grid ids", zap.Error(err))
		es.status = FailedStatus
		return
	}
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
		grid, ok := grids[timeNum]
		if !ok {
			unilog.Logger().Error("unable to get grid", zap.Int64("grid number", timeNum))
			continue
		}

		startTime := t[0].Unix()
		finishTime := t[1].Unix()
		posts, _, err := cl.SelectPosts(context.Background(), es.eventReq.CityId, startTime, finishTime)
		if err != nil {
			unilog.Logger().Error("unable to get posts from data storage", zap.Error(err))
			continue
		}

		expireTime := t[0].Add(-25 * time.Hour).Unix()
		oldEvents, err := cl.PullEventsWithIDs(context.Background(), es.eventReq.CityId, expireTime, startTime)
		if err != nil {
			unilog.Logger().Error("unable to get events from data storage", zap.Error(err))
			continue
		}

		filterTags := filterTags(es.eventReq.FilterTags)
		newEvents, updatedEvents, deletedEvents, found := detection.FindEvents(grid, posts, es.cfg.MaxPoints, filterTags, startTime, finishTime, oldEvents)
		if found {
			unilog.Logger().Info("found changed events", zap.String("session", es.id),
				zap.Int("num of new events", len(newEvents)), zap.Int("num of updated events", len(updatedEvents)),
				zap.Int("num of deleted events", len(deletedEvents)), zap.String("timestamp", t[0].String()))

			if len(newEvents) > 0 {
				err = cl.PushEvents(context.Background(), es.eventReq.CityId, newEvents)
				if err != nil {
					unilog.Logger().Error("unable to push events to data storage", zap.Error(err))
					return
				}
			}

			if len(updatedEvents) > 0 {
				err = cl.UpdateEvents(context.Background(), es.eventReq.CityId, updatedEvents)
				if err != nil {
					unilog.Logger().Error("unable to update events in data storage", zap.Error(err))
					return
				}
			}

			if len(deletedEvents) > 0 {
				ids := make([]int64, len(deletedEvents))
				for i, deletedEvent := range deletedEvents {
					ids[i] = deletedEvent.ID
				}
				err = cl.DeleteEvents(context.Background(), es.eventReq.CityId, ids)
				if err != nil {
					unilog.Logger().Error("unable to delete events in data storage", zap.Error(err))
					return
				}
			}
		} else {
			unilog.Logger().Info("no events found", zap.String("session", es.id), zap.String("timestamp", t[0].String()), zap.Error(err))
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
		r := [2]time.Time{c.Add(-1 * time.Hour), c}
		res = append(res, r)
		c = c.Add(10 * time.Minute)
	}
	return res, nil
}

func generateGridIds(startDate, finishDate int64, tz string) ([]int64, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		unilog.Logger().Error("unable to load timezone", zap.Error(err))
		return nil, err
	}
	startTime := time.Unix(startDate, 0)
	startTime = startTime.In(loc)
	finishTime := time.Unix(finishDate, 0)
	finishTime = finishTime.In(loc)
	t := startTime.Add(-1 * time.Hour)
	res := []int64{}
	set := map[int64]bool{}
	for t.Before(finishTime) {
		v := getGridNum(t.Month(), t.Weekday(), t.Hour())
		t = t.Add(time.Hour)
		if _, ok := set[v]; !ok {
			set[v] = true
			res = append(res, v)
		}
	}
	return res, nil
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

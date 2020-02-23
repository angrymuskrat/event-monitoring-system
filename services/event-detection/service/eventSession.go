package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"sync"
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

const (
	tagsPath     = "/home/alexvish/monitoring/event-detection/spb_tags_base.txt"
	firstGridNum = 1100
	lastGridNum  = 12224
)

type eventSession struct {
	id       string
	status   StatusType
	cfg      Config
	eventReq proto.EventRequest
	grids    map[int64][]byte
}

func newEventSession(config Config, eventReq proto.EventRequest, id string) *eventSession {
	return &eventSession{
		id:       id,
		status:   RunningStatus,
		cfg:      config,
		eventReq: eventReq,
		grids:    make(map[int64][]byte),
	}
}

func (es *eventSession) detectEvents() {
	conn, err := grpc.Dial(es.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data storage", zap.Error(err))
		es.status = FailedStatus
		return
	}
	client := service.NewGRPCClient(conn)
	ids := generateGridIds(es.eventReq.StartTime, es.eventReq.FinishTime)
	es.grids, err = client.PullGrid(context.Background(), es.eventReq.CityId, ids)
	if err != nil {
		unilog.Logger().Error("unable to get grids from data storage", zap.Error(err))
		es.status = FailedStatus
		return
	}

	times, err := getTimes(es.eventReq.StartTime, es.eventReq.FinishTime, es.eventReq.Timezone)
	if err != nil {
		unilog.Logger().Error("unable to generate intervals", zap.Error(err))
		es.status = FailedStatus
		return
	}

	wg := &sync.WaitGroup{}
	ewg := &sync.WaitGroup{}
	evChan := make(chan []data.Event)
	go loadEvents(evChan, es.cfg.DataStorageAddress, es.eventReq.CityId, ewg, &err)
	ewg.Add(1)

	timeChan := make(chan time.Time)
	for w := 0; w < es.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go es.eventWorker(wg, timeChan, evChan)
	}
	for _, t := range times {
		timeChan <- t
	}
	close(timeChan)
	wg.Wait()
	close(evChan)
	ewg.Wait()

	if err != nil {
		unilog.Logger().Error("error during pushing events to data storage", zap.Error(err))
		es.status = FailedStatus
		return
	}
	es.status = FinishedStatus
}

func getTimes(start, finish int64, tz string) ([]time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		unilog.Logger().Error("unable to load timezone", zap.Error(err))
		return nil, err
	}
	res := []time.Time{}
	s := time.Unix(start, 0)
	s = s.In(loc)
	f := time.Unix(finish, 0)
	f = f.In(loc)
	c := s
	for c.Before(f) {
		res = append(res, c)
		c = c.Add(time.Hour)
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

func (es *eventSession) eventWorker(wg *sync.WaitGroup, timeChan chan time.Time, eChan chan []data.Event) {
	defer wg.Done()
	conn, err := grpc.Dial(es.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		return
	}
	cl := service.NewGRPCClient(conn)

	for t := range timeChan {
		timeNum := getGridNum(t.Month(), t.Weekday(), t.Hour())
		buf := bytes.NewBuffer(es.grids[timeNum])
		dec := gob.NewDecoder(buf)

		var grid convtree.ConvTree

		if err := dec.Decode(&grid); err != nil {
			unilog.Logger().Error("unable to decode grid", zap.Error(err))
			es.status = FailedStatus
			return
		}

		startTime := t.Unix()
		finishTime := t.Unix() + 3600
		posts, _, err := cl.SelectPosts(context.Background(), es.eventReq.CityId, startTime, finishTime)
		if err != nil {
			unilog.Logger().Error("unable to get posts from data storage", zap.Error(err))
			continue
		}

		filterTags := map[string]bool{}
		f, err := os.OpenFile(tagsPath, os.O_RDONLY, 0644)
		if err == nil {
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				tag := "#" + scanner.Text()
				filterTags[tag] = false
			}
		}
		evs, found := detection.FindEvents(grid, posts, es.cfg.MaxPoints, filterTags, startTime, finishTime)
		if found {
			unilog.Logger().Info("found events", zap.String("session", es.id), zap.Int("num", len(evs)))
			eChan <- evs
		}
	}
}

func loadEvents(eChan chan []data.Event, storageEp string, cityID string, wg *sync.WaitGroup, outErr *error) {
	defer wg.Done()
	conn, err := grpc.Dial(storageEp, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data storage", zap.Error(err))
		outErr = &err
		return
	}
	client := service.NewGRPCClient(conn)
	for evs := range eChan {
		err = client.PushEvents(context.Background(), cityID, evs)
		if err != nil {
			unilog.Logger().Error("unable to push events to data storage", zap.Error(err))
			outErr = &err
			return
		}
	}
}

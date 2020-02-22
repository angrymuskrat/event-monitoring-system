package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/detection"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type historicSession struct {
	id       string
	status   StatusType
	cfg      Config
	histReq  proto.HistoricRequest
	gridChan chan interval
	grids    map[int64][]byte
	mut      sync.Mutex
}

func newHistoricSession(config Config, histReq proto.HistoricRequest, id string) *historicSession {

	return &historicSession{
		id:       id,
		status:   RunningStatus,
		cfg:      config,
		histReq:  histReq,
		gridChan: make(chan interval),
		grids:    make(map[int64][]byte),
	}
}

func (hs *historicSession) generateGrids() {
	intervals, err := getIntervals(hs.histReq.StartTime, hs.histReq.FinishTime, hs.histReq.Timezone)
	if err != nil {
		unilog.Logger().Error("unable to generate intervals", zap.Error(err))
		hs.status = FailedStatus
		return
	}
	area := hs.histReq.Area
	wg := &sync.WaitGroup{}
	for w := 1; w <= hs.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go hs.gridWorker(wg, *area)
	}
	for k, ils := range intervals {
		il := interval{
			key:   k,
			value: ils,
		}
		hs.gridChan <- il
	}
	close(hs.gridChan)
	wg.Wait()
	conn, err := grpc.Dial(hs.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		return
	}
	cl := service.NewGRPCClient(conn)
	err = cl.PushGrid(context.Background(), hs.histReq.CityId, hs.grids)
	if err != nil {
		unilog.Logger().Error("unable to push grid to data storage", zap.Error(err))
		hs.status = FailedStatus
		return
	}
	hs.status = FinishedStatus
}

func getIntervals(start, finish int64, tz string) (map[int64][][2]int64, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		unilog.Logger().Error("unable to load timezone", zap.Error(err))
		return nil, err
	}
	res := map[int64][][2]int64{}
	s := time.Unix(start, 0)
	s = s.In(loc)
	if s.Month() != time.January || s.Day() != 1 || s.Hour() != 0 || s.Minute() != 0 || s.Second() != 0 {
		return nil, errors.New("start timestamp must be 1 January 00:00:00")
	}
	f := time.Unix(finish, 0)
	f = f.In(loc)
	if f.Month() != time.January || f.Day() != 1 || f.Hour() != 0 || f.Minute() != 0 || f.Second() != 0 {
		return nil, errors.New("finish timestamp must be 1 January 00:00:00")
	}
	c := s
	for !c.Equal(f) {
		wd := c.Weekday()
		var w int64
		switch wd {
		case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
			w = 1
		case time.Saturday, time.Sunday:
			w = 2
		}
		k := 1000*int64(c.Month()) + 100*w + int64(c.Hour())
		v, ok := res[k]
		if !ok {
			v = [][2]int64{}
		}
		t := c.Unix()
		v = append(v, [2]int64{t, t + 3600})
		res[k] = v
		c = c.Add(time.Hour)
	}
	return res, nil
}

func (hs *historicSession) gridWorker(wg *sync.WaitGroup, area data.Area) {
	conn, err := grpc.Dial(hs.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		return
	}
	cl := service.NewGRPCClient(conn)
	defer wg.Done()
	for id := range hs.gridChan {
		posts := []data.Post{}
		for _, i := range id.value {
			ps, _, err := cl.SelectPosts(context.Background(), hs.histReq.CityId, i[0], i[1])
			if err != nil {
				unilog.Logger().Error("unable to get posts from data strorage", zap.Error(err))
				continue
			}
			posts = append(posts, ps...)
		}
		if len(posts) == 0 {
			continue
		}
		grid, err := detection.HistoricGrid(posts, *area.TopLeft, *area.BotRight, hs.cfg.MaxPoints, hs.histReq.Timezone, hs.histReq.GridSize)
		if err != nil {
			unilog.Logger().Error("can't generate grid", zap.Error(err))
			hs.status = FailedStatus
			return
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)

		if err := enc.Encode(grid); err != nil {
			hs.status = FailedStatus
			unilog.Logger().Error("can't encode grid", zap.Error(err))
			return
		}

		hs.mut.Lock()
		hs.grids[id.key] = buf.Bytes()
		hs.mut.Unlock()
	}
}

type interval struct {
	key   int64
	value [][2]int64
}

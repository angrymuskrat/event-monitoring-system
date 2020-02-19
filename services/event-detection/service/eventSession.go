package service

import (
	"bytes"
	"context"
	"encoding/gob"
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
	firstGridNum = 1100
	lastGridNum  = 12224
)

type eventSession struct {
	id          string
	status      StatusType
	cfg         Config
	eventReq    proto.EventRequest
	postsChan   chan data.Post
	gridsChan   chan int64
	sortedPosts map[int64][]data.Post
	grids       map[int64][]byte
	events      []data.Event
	mut         sync.Mutex
}

func newEventSession(config Config, eventReq proto.EventRequest, id string) *eventSession {
	return &eventSession{
		id:          id,
		status:      RunningStatus,
		cfg:         config,
		eventReq:    eventReq,
		postsChan:   make(chan data.Post),
		gridsChan:   make(chan int64),
		sortedPosts: make(map[int64][]data.Post),
		grids:       make(map[int64][]byte),
	}
}

func (es *eventSession) detectEvents() {
	conn, err := grpc.Dial(es.cfg.DataStorageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		es.status = FailedStatus
		panic(err)
	}
	client := service.NewGRPCClient(conn)

	posts, area, err := client.SelectPosts(context.Background(), es.eventReq.CityId, es.eventReq.StartTime, es.eventReq.FinishDate)
	if err != nil {
		unilog.Logger().Error("unable to get posts from data strorage", zap.Error(err))
		es.status = FailedStatus
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for w := 1; w <= es.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go es.readWorker(wg, *area)
	}
	for _, post := range posts {
		es.postsChan <- post
	}
	close(es.postsChan)
	wg.Wait()

	startId, finishId := convertDatesToGridIds(es.eventReq.StartTime, es.eventReq.FinishDate)

	es.grids, err = client.PullGrid(context.Background(), es.eventReq.CityId, startId, finishId)
	if err != nil {
		unilog.Logger().Error("unable to get grids from data strorage", zap.Error(err))
		es.status = FailedStatus
		panic(err)
	}

	wg = &sync.WaitGroup{}
	for w := 1; w <= es.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go es.eventWorker(wg)
	}
	for id, _ := range es.grids {
		es.gridsChan <- id
	}
	close(es.gridsChan)
	wg.Wait()

	err = client.PushEvents(context.Background(), es.eventReq.CityId, es.events)
	if err != nil {
		unilog.Logger().Error("unable to push events to data storage", zap.Error(err))
		es.status = FailedStatus
		panic(err)
	}
	es.status = FinishedStatus

}

func (es *eventSession) readWorker(wg *sync.WaitGroup, area data.Area) {
	defer wg.Done()
	for post := range es.postsChan {
		if post.Lat <= area.TopLeft.Lat && post.Lat >= area.BotRight.Lat && post.Lon >= area.TopLeft.Lon && post.Lon <= area.BotRight.Lon {
			postTime := time.Unix(post.Timestamp, 0)
			loc, err := time.LoadLocation(es.eventReq.Timezone)
			if err != nil {
				unilog.Logger().Error("can't load location", zap.Error(err))
				es.status = FailedStatus
				panic(err)
			}
			postTime = postTime.In(loc)
			gridNum := getGridNum(postTime.Month(), postTime.Weekday(), postTime.Hour())
			es.mut.Lock()
			es.sortedPosts[gridNum] = append(es.sortedPosts[gridNum], post)
			es.mut.Unlock()
		}
	}
}

func convertDatesToGridIds(startDate, finishDate int64) (int64, int64) {
	startTime := time.Unix(startDate, 0)
	finishTime := time.Unix(finishDate, 0)
	if finishTime.Sub(startTime) > time.Hour*24*365 {
		return firstGridNum, lastGridNum
	}
	startNum := getGridNum(startTime.Month(), startTime.Weekday(), startTime.Hour())
	finishNum := getGridNum(finishTime.Month(), finishTime.Weekday(), finishTime.Hour())
	return startNum, finishNum
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

func (es *eventSession) eventWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range es.gridsChan {
		buf := bytes.NewBuffer(es.grids[id])
		dec := gob.NewDecoder(buf)

		var grid convtree.ConvTree

		if err := dec.Decode(&grid); err != nil {
			unilog.Logger().Error("can't decode grid", zap.Error(err))
			es.status = FailedStatus
			panic(err)
		}

		event, err := detection.FindEvents(grid, es.sortedPosts[id], es.cfg.MaxPoints, make(map[string]bool), es.eventReq.StartTime, es.eventReq.FinishDate)
		if err != nil {
			unilog.Logger().Error("can't generate grid", zap.Error(err))
			es.status = FailedStatus
			panic(err)
		}

		es.mut.Lock()
		es.events = append(es.events, event...)
		es.mut.Unlock()
	}
}

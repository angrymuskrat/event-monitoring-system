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
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type historicSession struct {
	id          string
	status      StatusType
	cfg         Config
	histReq     proto.HistoricRequest
	postsChan   chan data.Post
	gridChan    chan int64
	sortedPosts map[int64][]data.Post
	grids       map[int64][]byte
	mut         sync.Mutex
}

func newHistoricSession(config Config, histReq proto.HistoricRequest, id string) *historicSession {
	return &historicSession{
		id:          id,
		status:      RunningStatus,
		cfg:         config,
		histReq:     histReq,
		postsChan:   make(chan data.Post),
		gridChan:    make(chan int64),
		sortedPosts: make(map[int64][]data.Post),
		grids:       make(map[int64][]byte),
	}
}

func (hs *historicSession) generateGrids() {
	unilog.Logger().Info(hs.cfg.DataStorageAddress)
	conn, err := grpc.Dial(hs.cfg.DataStorageAddress, grpc.WithInsecure())
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		hs.status = FailedStatus
		panic(err)
	}
	client := service.NewGRPCClient(conn)

	posts, area, err := client.SelectPosts(context.Background(), hs.histReq.CityId, hs.histReq.StartTime, hs.histReq.FinishDate)
	if err != nil {
		unilog.Logger().Error("unable to get posts from data strorage", zap.Error(err))
		hs.status = FailedStatus
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for w := 1; w <= hs.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go hs.readWorker(wg, *area)
	}
	for _, post := range posts {
		hs.postsChan <- post
	}
	close(hs.postsChan)
	wg.Wait()
	wg = &sync.WaitGroup{}
	for w := 1; w <= hs.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go hs.gridWorker(wg, *area)
	}

	for gridID := range hs.sortedPosts {
		hs.gridChan <- gridID
	}
	close(hs.gridChan)
	wg.Wait()

	err = client.PushGrid(context.Background(), hs.histReq.CityId, hs.grids)
	if err != nil {
		unilog.Logger().Error("can't push grid to data storage", zap.Error(err))
		hs.status = FailedStatus
		panic(err)
	}
	hs.status = FinishedStatus
}

func (hs *historicSession) readWorker(wg *sync.WaitGroup, area data.Area) {
	defer wg.Done()
	for post := range hs.postsChan {
		//if post.Lat <= area.TopLeft.Lat && post.Lat >= area.BotRight.Lat && post.Lon >= area.TopLeft.Lon && post.Lon <= area.BotRight.Lon {
		postTime := time.Unix(post.Timestamp, 0)
		loc, err := time.LoadLocation(hs.histReq.Timezone)
		if err != nil {
			unilog.Logger().Error("can't load location", zap.Error(err))
			hs.status = FailedStatus
			panic(err)
		}
		postTime = postTime.In(loc)
		gridNum := getGridNum(postTime.Month(), postTime.Weekday(), postTime.Hour())
		hs.mut.Lock()
		hs.sortedPosts[gridNum] = append(hs.sortedPosts[gridNum], post)
		hs.mut.Unlock()
		//}
	}
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

func (hs *historicSession) gridWorker(wg *sync.WaitGroup, area data.Area) {
	defer wg.Done()

	for id := range hs.gridChan {
		grid, err := detection.HistoricGrid(hs.sortedPosts[id], *area.TopLeft, *area.BotRight, hs.cfg.MaxPoints, hs.histReq.Timezone, hs.histReq.GridSize)
		if err != nil {
			unilog.Logger().Error("can't generate grid", zap.Error(err))
			hs.status = FailedStatus
			panic(err)
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)

		if err := enc.Encode(grid); err != nil {
			hs.status = FailedStatus
			unilog.Logger().Error("can't encode grid", zap.Error(err))
			panic(err)
		}

		hs.mut.Lock()
		hs.grids[id] = buf.Bytes()
		hs.mut.Unlock()
	}
}

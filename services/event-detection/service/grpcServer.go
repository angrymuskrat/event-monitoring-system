package service

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/detection"

	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
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
	cfg Config
}

func newGRPCServer(cfg Config) grpcServer {
	return grpcServer{cfg: cfg}
}

func (gs grpcServer) HistoricGrids(ctx context.Context, histReq *proto.HistoricRequest) (*proto.HistoricResponse, error) {
	conn, err := grpc.Dial(gs.cfg.DataStorageAddress)
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		panic(err)
	}
	ds := service.NewGRPCClient(conn)

	posts, area, err := ds.SelectPosts(context.Background(), histReq.CityId, histReq.StartTime, histReq.FinishDate)
	if err != nil {
		unilog.Logger().Error("unable to get posts from data strorage", zap.Error(err))
		panic(err)
	}

	var postsChan chan data.Post
	var sortedPosts map[string][]data.Post
	var mut sync.Mutex

	wg := &sync.WaitGroup{}
	for w := 1; w <= gs.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go readWorker(wg, histReq.Timezone, postsChan, *area, sortedPosts, &mut)
	}

	for _, post := range posts {
		postsChan <- post
	}
	wg.Wait()

	var gridChan chan []data.Post
	for w := 1; w <= gs.cfg.WorkersNumber; w++ {
		wg.Add(1)
		go gridWorker(wg, gridChan, *area.TopLeft, *area.BotRight, gs.cfg.MaxPoints, histReq.Timezone, histReq.GridSize)
	}

	for _, gridPosts := range sortedPosts {
		gridChan <- gridPosts
	}

	return nil, nil
}

func readWorker(wg *sync.WaitGroup, timezone string, postsChan chan data.Post, area data.Area, sortedPosts map[string][]data.Post, mut *sync.Mutex) {
	defer wg.Done()
	for post := range postsChan {
		if post.Lat <= area.TopLeft.Lat && post.Lat >= area.BotRight.Lat && post.Lon >= area.TopLeft.Lon && post.Lon <= area.BotRight.Lon {
			postTime := time.Unix(post.Timestamp, 0)
			loc, err := time.LoadLocation(timezone)
			if err != nil {
				unilog.Logger().Error("can't load location", zap.Error(err))
				panic(err)
			}
			postTime = postTime.In(loc)
			gridName := getGridName(postTime.Month(), postTime.Weekday(), postTime.Hour())
			mut.Lock()
			sortedPosts[gridName] = append(sortedPosts[gridName], post)
			mut.Unlock()
		}
	}
}

func getGridName(month time.Month, day time.Weekday, hour int) string {
	monthString := month.String()
	var dayString string
	switch day {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		dayString = "working"
	case time.Saturday, time.Sunday:
		dayString = "weekend"
	}
	hourString := strconv.Itoa(hour)
	gridName := monthString + "-" + dayString + "-" + hourString
	return gridName
}

func gridWorker(wg *sync.WaitGroup, gridChan chan []data.Post, topLeft, bottomRight data.Point, maxPoints int, tz string, gridSize float64) {
	defer wg.Done()
	for gridPosts := range gridChan {
		detection.HistoricGrid(gridPosts, topLeft, bottomRight, maxPoints, tz, gridSize)
	}
}

func (gs grpcServer) HistoricStatus(context.Context, *proto.StatusRequest) (*proto.StatusResponse, error) {
	return nil, nil
}

func (gs grpcServer) FindEvents(context.Context, *proto.EventRequest) (*proto.EventResponse, error) {
	return nil, nil
}

func (gs grpcServer) EventsStatus(context.Context, *proto.StatusRequest) (*proto.StatusResponse, error) {
	return nil, nil
}

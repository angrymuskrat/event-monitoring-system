package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/auth/basic"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Start(confPath string) {
	conf, err := readConfig(confPath)
	if err != nil {
		return
	}
	logger := setupLog(conf.LogPath)
	var svc BackendService
	svc = &backendService{
		storageConn: MockConnector{},
	}
	svc = &loggingMiddleware{logger, svc}
	r := mux.NewRouter()
	r.Methods("GET").Path("/heatmap/{city}/{topLeft}/{botRight}/{hour}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeHeatmapEndpoint(svc)),
		decodeHeatmapRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/timeline/{city}/{start}/{finish}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeTimelineEndpoint(svc)),
		decodeTimelineRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/events/{city}/{topLeft}/{botRight}/{hour}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeEventsEndpoint(svc)),
		decodeEventsRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/search/{city}/{tags}/{start}/{finish}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeEventsSearchEndpoint(svc)),
		decodeSearchRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	http.Handle("/", r)
	err = http.ListenAndServe(conf.Address, nil)
	if err != nil {
		unilog.Logger().Error("error in service handler", zap.Error(err))
	}
}

func decodeHeatmapRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := HeatmapRequest{}
	city, ok := vars["city"]
	if !ok {
		return nil, errors.New("unable to get city name")
	}
	req.City = city
	topLeftRaw, ok := vars["topLeft"]
	if !ok {
		return nil, errors.New("unable to get top left coordinates")
	}
	coords := strings.Split(topLeftRaw, ",")
	if len(coords) != 2 {
		return nil, errors.New("incorrect format of top left coordinates")
	}
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, errors.New("unable to parse latitude of top left")
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, errors.New("unable to parse longitude of top left")
	}
	req.TopLeft = data.Point{
		Lat: lat,
		Lon: lon,
	}
	botRightRaw, ok := vars["botRight"]
	if !ok {
		return nil, errors.New("unable to get bottom right coordinates")
	}
	coords = strings.Split(botRightRaw, ",")
	if len(coords) != 2 {
		return nil, errors.New("incorrect format of bottom right coordinates")
	}
	lat, err = strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, errors.New("unable to parse latitude of bottom right")
	}
	lon, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, errors.New("unable to parse longitude of bottom right")
	}
	req.BottomRight = data.Point{
		Lat: lat,
		Lon: lon,
	}
	hourRaw, ok := vars["hour"]
	if !ok {
		return nil, errors.New("unable to get hour")
	}
	hour, err := strconv.ParseInt(hourRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of hour")
	}
	req.Hour = hour
	return req, nil
}

func decodeTimelineRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := TimelineRequest{}
	city, ok := vars["city"]
	if !ok {
		return nil, errors.New("unable to get city name")
	}
	req.City = city
	startRaw, ok := vars["start"]
	if !ok {
		return nil, errors.New("unable to get start")
	}
	start, err := strconv.ParseInt(startRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of start")
	}
	req.Start = start
	finishRaw, ok := vars["finish"]
	if !ok {
		return nil, errors.New("unable to get finish")
	}
	finish, err := strconv.ParseInt(finishRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of finish")
	}
	req.Finish = finish
	return req, nil
}

func decodeEventsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := EventsRequest{}
	city, ok := vars["city"]
	if !ok {
		return nil, errors.New("unable to get city name")
	}
	req.City = city
	topLeftRaw, ok := vars["topLeft"]
	if !ok {
		return nil, errors.New("unable to get top left coordinates")
	}
	coords := strings.Split(topLeftRaw, ",")
	if len(coords) != 2 {
		return nil, errors.New("incorrect format of top left coordinates")
	}
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, errors.New("unable to parse latitude of top left")
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, errors.New("unable to parse longitude of top left")
	}
	req.TopLeft = data.Point{
		Lat: lat,
		Lon: lon,
	}
	botRightRaw, ok := vars["botRight"]
	if !ok {
		return nil, errors.New("unable to get bottom right coordinates")
	}
	coords = strings.Split(botRightRaw, ",")
	if len(coords) != 2 {
		return nil, errors.New("incorrect format of bottom right coordinates")
	}
	lat, err = strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, errors.New("unable to parse latitude of bottom right")
	}
	lon, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, errors.New("unable to parse longitude of bottom right")
	}
	req.BottomRight = data.Point{
		Lat: lat,
		Lon: lon,
	}
	hourRaw, ok := vars["hour"]
	if !ok {
		return nil, errors.New("unable to get hour")
	}
	hour, err := strconv.ParseInt(hourRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of hour")
	}
	req.Hour = hour
	return req, nil
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := SearchRequest{}
	city, ok := vars["city"]
	if !ok {
		return nil, errors.New("unable to get city name")
	}
	req.City = city
	startRaw, ok := vars["start"]
	if !ok {
		return nil, errors.New("unable to get start")
	}
	start, err := strconv.ParseInt(startRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of start")
	}
	req.Start = start
	finishRaw, ok := vars["finish"]
	if !ok {
		return nil, errors.New("unable to get finish")
	}
	finish, err := strconv.ParseInt(finishRaw, 10, 64)
	if err != nil {
		return nil, errors.New("incorrect format of finish")
	}
	req.Finish = finish
	tagsRaw, ok := vars["tags"]
	tags := strings.Split(tagsRaw, ",")
	req.Keytags = tags
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func setupLog(path string) *zap.Logger {
	conf := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	if len(path) > 0 {
		conf.OutputPaths = append(conf.OutputPaths, path)
		conf.ErrorOutputPaths = append(conf.ErrorOutputPaths, path)
	}
	log, err := conf.Build()
	if err != nil {
		fmt.Println("unable to initialize log")
		fmt.Println(err)
		log = defaultLog()
	}
	return log
}

func defaultLog() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	return zap.New(core)
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/go-kit/kit/auth/basic"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(confPath string, cr *crawler.Crawler) {
	conf, err := readConfig(confPath)
	if err != nil {
		return
	}
	logger := setupLog(conf.LogPath)
	var svc CrawlerService
	svc = &crawlerService{
		crawler: cr,
	}
	svc = &loggingMiddleware{logger, svc}
	r := mux.NewRouter()
	r.Methods("POST").Path("/new").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeNewEndpoint(svc)),
		decodeNewRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/status/{id}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeStatusEndpoint(svc)),
		decodeStatusRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/entities/{id}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeEntitiesEndpoint(svc)),
		decodeStatusRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/posts/{id}/{cursor}/{num}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makePostsEndpoint(svc)),
		decodePostsRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("POST").Path("/stop").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeStopEndpoint(svc)),
		decodeStopRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	http.Handle("/", r)
	err = http.ListenAndServe(conf.Address, nil)
	if err != nil {
		unilog.Logger().Error("error in service handler", zap.Error(err))
	}
}

func decodeNewRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var p crawler.Parameters
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("unable to get session ID")
	}
	return idEpRequest{
		ID: id,
	}, nil
}

func decodePostsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("unable to get session ID")
	}
	cursor, ok := vars["cursor"]
	if !ok {
		return nil, errors.New("unable to get cursor")
	}
	if cursor == "none" {
		cursor = ""
	}
	numRaw, ok := vars["num"]
	if !ok {
		return nil, errors.New("unable to get number of posts")
	}
	num, err := strconv.Atoi(numRaw)
	if err != nil {
		return nil, err
	}
	return postsEpRequest{
		ID:     id,
		Offset: cursor,
		Num:    num,
	}, nil
}

func decodeStopRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req idEpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
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

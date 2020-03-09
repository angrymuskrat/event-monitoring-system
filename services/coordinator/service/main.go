package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	storage "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/go-kit/kit/auth/basic"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net/http"
	"os"
)

var Storage storage.GrpcService

func Start(confPath string) {
	conf, err := readConfig(confPath)
	if err != nil {
		return
	}
	conn, err := grpc.Dial(conf.DataStorage.Address, grpc.WithInsecure(), grpc.WithMaxMsgSize(storage.MaxMsgSize))
	if err != nil {
		unilog.Logger().Error("do not be able to connect to data-storage", zap.Error(err))
		return
	}
	Storage = storage.NewGRPCClient(conn)

	endpoints := ServiceEndpoints{
		Crawler:        conf.Crawler,
		EventDetection: conf.EventDetection,
	}
	//logger := setupLog(conf.LogPath)
	var svc CoordinatorService
	svc = &coordinatorService{
		endpoints: endpoints,
	}
	//svc = &loggingMiddleware{logger, svc} //TODO: add logging middleware
	r := mux.NewRouter()
	r.Methods("POST").Path("/new-session").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeNewSessionEndpoint(svc)),
		decodeNewSessionRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/status/{id}").Handler(httptransport.NewServer(
		basic.AuthMiddleware(conf.User, conf.Password, "realm")(makeStatusEndpoint(svc)),
		decodeStatusRequest,
		encodeResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	http.Handle("/", accessControl(r))
	unilog.Logger().Info("successfully started")
	err = http.ListenAndServe(conf.Address, nil)
	if err != nil {
		unilog.Logger().Error("error in service handler", zap.Error(err))
	}
}

func decodeNewSessionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SessionParameters
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		unilog.Logger().Error("unable to decode request", zap.Error(err))
		return nil, err
	}
	return req, nil
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("unable to get city name")
	}
	return id, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
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
		conf.OutputPaths = []string{path}
		conf.ErrorOutputPaths = []string{path}
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

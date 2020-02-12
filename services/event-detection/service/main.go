package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, confPath string) {
	conf, err := readConfig(confPath)
	if err != nil {
		return
	}
	//logger := setupLog(conf.LogPath)
	var svc Service
	svc = newEventService()
	// TODO: implement logging middleware
	//svc = &loggingMiddleware{logger, svc}
	grpcServer := Server(svc)

	var g group.Group

	grpcListener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		unilog.Logger().Error("error in transport gRPC Listener", zap.Error(err))
		os.Exit(1)
	}
	g.Add(func() error {
		unilog.Logger().Info("start gRPC transport", zap.String("url", conf.Address))
		baseServer := grpc.NewServer(
			grpc.UnaryInterceptor(kitgrpc.Interceptor),
			grpc.MaxRecvMsgSize(MaxMsgSize),
		)
		proto.RegisterEventDetectionServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)

		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})

	g.Run()
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

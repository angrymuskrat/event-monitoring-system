package main

import (
	"context"
	"flag"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/storage"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

func main() {
	aLog := flag.String("al", "app.log", "path to application log file")
	serviceConfig := flag.String("sc", "service.toml", "path to service configuration file")
	connectorConfig := flag.String("cc", "storage.toml", "path to db storage configuration file")
	flag.Parse()

	logCfg := unilog.DefaultConfig()
	logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if len(*aLog) > 0 {
		logCfg.OutputPaths = []string{*aLog}
		logCfg.ErrorOutputPaths = []string{*aLog}
	}
	unilog.InitLog(logCfg)

	dbConnector, err := storage.New(context.Background(), *connectorConfig)
	if err != nil {
		fmt.Print(err)
		return
	}

	storagesvc.Start(context.Background(), *serviceConfig, dbConnector)
}

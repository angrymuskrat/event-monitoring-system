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
	aLog := flag.String("al", "services/data-storage/app.log", "path to application log file")
	serviceConfig := flag.String("sc", "services/data-storage/cmd/service/service.toml", "path to service configuration file")
	connectorConfig := flag.String("cc", "services/data-storage/cmd/service/storage.toml", "path to db storage configuration file")
	flag.Parse()

	logCfg := unilog.DefaultConfig()
	logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if len(*aLog) > 0 {
		logCfg.OutputPaths = append(logCfg.OutputPaths, *aLog)
		logCfg.ErrorOutputPaths = append(logCfg.ErrorOutputPaths, *aLog)
	}
	unilog.InitLog(logCfg)

	dbConnector, err := storage.NewStorage(context.Background(), *connectorConfig)
	if err != nil {
		fmt.Print(err)
		return
	}

	storagesvc.Start(context.Background(), *serviceConfig, dbConnector)
}

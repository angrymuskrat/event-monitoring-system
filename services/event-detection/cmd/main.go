package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

func main() {
	lp := flag.String("log", "./log.txt", "path to the log file")
	cp := flag.String("config", "./config.toml", "path to the config file")
	flag.Parse()

	logCfg := unilog.DefaultConfig()
	logCfg.OutputPaths = append(logCfg.OutputPaths, *lp)
	logCfg.ErrorOutputPaths = append(logCfg.ErrorOutputPaths, *lp)
	unilog.InitLog(logCfg)

	var cfg service.Config
	_, err := toml.DecodeFile(*cp, &cfg)
	if err != nil {
		unilog.Logger().Error("unable to read config file", zap.String("path", *cp), zap.Error(err))
		panic(err)
	}

	service.ServerStart(cfg)
}

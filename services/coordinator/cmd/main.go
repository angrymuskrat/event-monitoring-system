package main

import (
	"flag"

	service "github.com/angrymuskrat/event-monitoring-system/services/coordinator/service"
	"github.com/visheratin/unilog"
)

func main() {
	lp := flag.String("log", "./log.txt", "path to the log file")
	cp := flag.String("config", "./config.toml", "path to the config file")
	flag.Parse()

	logCfg := unilog.DefaultConfig()
	logCfg.OutputPaths = append(logCfg.OutputPaths, *lp)
	logCfg.ErrorOutputPaths = append(logCfg.ErrorOutputPaths, *lp)
	unilog.InitLog(logCfg)

	service.Start(*cp)
}

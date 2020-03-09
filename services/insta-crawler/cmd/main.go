package main

import (
	"flag"

	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/service"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

func main() {
	aLog := flag.String("al", "app.log", "path to application log file")
	serviceConfig := flag.String("sc", "service.toml", "path to service configuration file")
	crawlerConfig := flag.String("cc", "crawler.toml", "path to crawler configuration file")
	flag.Parse()

	logCfg := unilog.DefaultConfig()
	logCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if len(*aLog) > 0 {
		logCfg.OutputPaths = []string{*aLog}
		logCfg.ErrorOutputPaths = []string{*aLog}
	}
	unilog.InitLog(logCfg)
	cr, err := crawler.NewCrawler(*crawlerConfig)
	if err != nil {
		return
	}
	service.Start(*serviceConfig, cr)
}

package main

import (
	"flag"
	"github.com/angrymuskrat/event-monitoring-system/services/backend/service"
)

func main() {
	serviceConfig := flag.String("sc", "config.toml", "path to service configuration file")
	flag.Parse()
	service.Start(*serviceConfig)
}

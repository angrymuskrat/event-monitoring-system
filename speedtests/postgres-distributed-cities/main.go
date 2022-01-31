package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"
)

func main() {
	configPath := flag.String("sc", "config.toml", "path to configuration file")
	config, err := readConfig(*configPath)
	if err != nil {
		return
	}
	//cities := []string{"moscow", "spb", "nyc", "london", "moscow", "spb", "nyc", "london"}
	dbNames := []string{"moscow", "spb", "nyc", "london"}
	sort.Strings(dbNames)
	labels := dbNames
	connectors, err := NewListConnectors(context.Background(), &config, dbNames, labels)

	generator := NewGenerator(config.Seed)
	times := generator.genHours(config.TimesCount, config.MinTime, config.MaxTime)
	hub := make(chan CityExperiment)
	var testsCount int
	var outputs []string

	fmt.Println("Parallel count posts")
	for _, conn := range connectors {
		go conn.runPostCountHour(hub, times, config.Restarts)
	}
	testsCount = 0
	outputs = []string{}
	for test := range hub {
		outputs = append(outputs, test.String())
		testsCount += 1
		if testsCount == len(connectors) {
			break
		}
	}
	sort.Strings(outputs)
	fmt.Println(strings.Join(outputs, "\n"))

	fmt.Println("Sequence count posts")
	for _, conn := range connectors {
		//go conn.runSpeedTestPostHour(hub, times, restarts)
		go conn.runPostCountHour(hub, times, config.Restarts)
		test := <-hub
		fmt.Println(test.String())
	}
	fmt.Println("")

	fmt.Println("Parallel data posts")
	for _, conn := range connectors {
		go conn.runPostHour(hub, times, config.Restarts)
	}
	testsCount = 0
	outputs = []string{}
	for test := range hub {
		outputs = append(outputs, test.String())
		testsCount += 1
		if testsCount == len(connectors) {
			break
		}
	}
	sort.Strings(outputs)
	fmt.Println(strings.Join(outputs, "\n"))

	fmt.Println("Sequence data posts")
	for _, conn := range connectors {
		go conn.runPostHour(hub, times, config.Restarts)
		test := <-hub
		fmt.Println(test.String())
	}
	fmt.Println("")

	fmt.Println("Parallel count events")
	for _, conn := range connectors {
		go conn.runEventCountDay(hub, times, config.Restarts)
	}
	testsCount = 0
	outputs = []string{}
	for test := range hub {
		outputs = append(outputs, test.String())
		testsCount += 1
		if testsCount == len(connectors) {
			break
		}
	}
	sort.Strings(outputs)
	fmt.Println(strings.Join(outputs, "\n"))

	fmt.Println("Sequence count events")
	for _, conn := range connectors {
		go conn.runEventCountDay(hub, times, config.Restarts)
		test := <-hub
		fmt.Println(test.String())
	}
	fmt.Println("")

	fmt.Println("Parallel data events")
	for _, conn := range connectors {
		go conn.runEventDay(hub, times, config.Restarts)
	}
	testsCount = 0
	outputs = []string{}
	for test := range hub {
		outputs = append(outputs, test.String())
		testsCount += 1
		if testsCount == len(connectors) {
			break
		}
	}
	sort.Strings(outputs)
	fmt.Println(strings.Join(outputs, "\n"))

	fmt.Println("Sequence data events")
	for _, conn := range connectors {
		go conn.runEventDay(hub, times, config.Restarts)
		test := <-hub
		fmt.Println(test.String())
	}
	fmt.Println("")
}

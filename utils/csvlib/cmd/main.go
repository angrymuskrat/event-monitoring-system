package main

import (
	"context"
	"flag"
	"github.com/angrymuskrat/event-monitoring-system/utils/csvlib"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"os"
)

func main() {
	connectorConfig := flag.String("config", "connection.toml", "path to configuration file")
	cityCode := flag.String("city", "nyc", "code of the city to the database of which the request is sent")
	requestType := flag.String("type", "LoadEventPosts", "type of request")
	rootPath := flag.String("root", "./", "working dir path")
	additionalParams := map[string]*string{}
	additionalParams["EventTableName"] = flag.String("EventTableName", "events_6", "name of events table for request connected with events")
	additionalParams["OutputFile"] = flag.String("OutputFile", "output.csv", "name of file for saving data")
	additionalParams["InputFile"] = flag.String("InputFile", "input.csv", "name of input files InputFile")
	intParams := map[string]*int64{}
	intParams["Start"] = flag.Int64("Start", 0, "Start timestamp for query, 0 - None and default value")
	intParams["Finish"] = flag.Int64("Finish", 0, "Finish timestamp for query, 0 - None and default value")
	flag.Parse()
	/*fmt.Printf("cityCode: %v\n", *cityCode)
	fmt.Printf("requestType: %v\n", *requestType)
	fmt.Printf("rootPath: %v\n", *rootPath)
	fmt.Printf("additionalParams: %v\n", *additionalParams["EventTableName"])*/

	connector, err := csvlib.NewConnector(context.Background(), connectorConfig, *cityCode)
	if err != nil {
		unilog.Logger().Error("don't be able to build connector", zap.Error(err))
		os.Exit(0)
	}
	err = connector.ExecuteRequest(context.Background(), *requestType, *rootPath, additionalParams, intParams)
	if err != nil {
		unilog.Logger().Error("don't be able exec request", zap.Error(err))
		os.Exit(0)
	}
}

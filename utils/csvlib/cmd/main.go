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
	additionalParams["EventPostsOutput"] = flag.String("EventsPostsFile", "events_posts.csv", "name of file for saving events")
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
	err = connector.ExecuteRequest(context.Background(), *requestType, *rootPath, additionalParams)
	if err != nil {
		unilog.Logger().Error("don't be able exec request", zap.Error(err))
		os.Exit(0)
	}
}

package main

import (
	"context"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"gocloud.dev/blob/fileblob"
	"google.golang.org/grpc"
	"os"
)



func main() {
	var (
		svc storagesvc.Service
		err error
	)
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithMaxMsgSize(storagesvc.MaxMsgSize))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	svc = storagesvc.NewGRPCClient(conn)
	testAll(svc)
	//test("pullLocations", svc)
}

var AllMethods = []string{"pushPosts", "selectPosts", "selectAggrPosts", "pullTimeline", "pushEvents", "pullEvents", "pushGrid", "pullGrid", "pushLocations", "pullLocations"}
func testAll(svc storagesvc.Service) {
	for _, method := range AllMethods {
		test(method, svc)
	}
}

func test (method string, svc storagesvc.Service) {
	switch method {
	case "pushPosts":
		testPosts := GeneratePosts(1000)
		res, err := svc.PushPosts(context.Background(), "nyc", testPosts)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))

	case "selectPosts":
		res, err := svc.SelectPosts(context.Background(), "nyc", data.SpatioTemporalInterval{MinTime: 0, MaxTime: 100000000000,
			TopLeft:  data.Point{Lat: -100, Lon: 100},
			BotRight: data.Point{Lat: 100, Lon: -100}})
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))

	case "selectAggrPosts":
		res, err := svc.SelectAggrPosts(context.Background(), "nyc", data.SpatioHourInterval{Hour: 1579622400,
			TopLeft:  data.Point{Lat: -100, Lon: -100},
			BotRight: data.Point{Lat: 100, Lon: 100}})
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))

	case "pullTimeline":
		res, err := svc.PullTimeline(context.Background(), "nyc", 0, 3600)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))

	case "pushGrid":
		bucket, err := fileblob.OpenBucket("tests/", nil)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		defer bucket.Close()
		ctx := context.Background()
		b, err := bucket.ReadAll(ctx, "foo100000000.blob")
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		id := RandString(20)
		err = svc.PushGrid(context.Background(), "nyc", id, b)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right", method)

	case "pullGrid":
		res, err := svc.PullGrid(context.Background(), "nyc", "asasasas")
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))

	case "pushEvents" :
		events := []data.Event{
			{Title: "test", Start: 10000, Finish: 100, Center: data.Point{Lat: 10, Lon: 10}, PostCodes: []string{"ffdsfs", "sfdsdfsf", "sfsfsf"}, Tags: []string{"#tag1", "#tag2"}},
			{Title: "test2", Start: 1, Finish: 100, Center: data.Point{Lat: 9, Lon: 9}, PostCodes: []string{"ffdsfs", "sfdsdfsf", "sfsfsf"}, Tags: []string{"#tag1", "#tag2"}},
		}
		err := svc.PushEvents(context.Background(), "nyc", events)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right", method)

	case "pullEvents" :
		res, err := svc.PullEvents(context.Background(), "nyc", data.SpatioHourInterval{Hour: 9000,
			TopLeft:  data.Point{Lat: 1, Lon: 1},
			BotRight: data.Point{Lat: 100, Lon: 100}})
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))
	case "pushLocations":
		//city := data.City{Title:"New York", ID:"nyc"}
		locations := []data.Location{
			{Title:"loc1", ID:"65", Slug:"slug1", Position: data.Point{Lat:1, Lon:1}},
			{Title:"loc2", ID:"2", Slug:"slug2", Position: data.Point{Lat:2, Lon:2}},
			{Title:"loc3", ID:"4", Slug:"slug4", Position: data.Point{Lat:3, Lon:3}},
		}
		err := svc.PushLocations(context.Background(),  "nyc",  locations)
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right", method)
	case "pullLocations":
		res, err := svc.PullLocations(context.Background(), "nyc")
		if err != nil {
			fmt.Printf("\n%v: error: %v", method, err)
			return
		}
		fmt.Printf( "\n%v: It is all right, len(res):%v", method, len(res))
	default:
		fmt.Printf("\nerror: invalid method %v", method)
		return
	}
}
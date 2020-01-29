package main

import (
	"context"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"gocloud.dev/blob/fileblob"
	"google.golang.org/grpc"
	"log"
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

	method := "pullTimeline"

	switch method {
	case "pushPosts":
		testPosts := GeneratePosts(1000)
		res, err := svc.PushPosts(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right\n%v", len(res))

	case "selectPosts":
		res, err := svc.SelectPosts(context.Background(), data.SpatioTemporalInterval{MinTime: 0, MaxTime: 100000000000,
			TopLeft:  &data.Point{Lat: -100, Lon: 100},
			BotRight: &data.Point{Lat: 100, Lon: -100}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", len(res))

	case "selectAggrPosts":
		res, err := svc.SelectAggrPosts(context.Background(), data.SpatioHourInterval{Hour: 1579622400,
			TopLeft:  &data.Point{Lat: -100, Lon: -100},
			BotRight: &data.Point{Lat: 100, Lon: 100}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", len(res))

	case "pullTimeline":
		res, err := svc.PullTimeline(context.Background(), "cityId", 0, 3600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", res)

	case "pushGrid":
		bucket, err := fileblob.OpenBucket("tests/", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer bucket.Close()
		ctx := context.Background()
		b, err := bucket.ReadAll(ctx, "foo100000000.blob")
		if err != nil {
			log.Fatal(err)
		}
		id := RandString(20)
		err = svc.PushGrid(context.Background(), id, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right, id: %v", id)

	case "pullGrid":
		res, err := svc.PullGrid(context.Background(), "asasasas")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right, len of res - \n%v", len(res))

	case "pushEvents" :
		events := []data.Event{
			{Title: "test", Start: 10000, Finish: 100, Center: data.Point{Lat: 10, Lon: 10}, PostCodes: []string{"ffdsfs", "sfdsdfsf", "sfsfsf"}, Tags: []string{"#tag1", "#tag2"}},
			{Title: "test2", Start: 1, Finish: 100, Center: data.Point{Lat: 9, Lon: 9}, PostCodes: []string{"ffdsfs", "sfdsdfsf", "sfsfsf"}, Tags: []string{"#tag1", "#tag2"}},
		}
		err := svc.PushEvents(context.Background(), events)
		fmt.Print(err)
	case "pullEvents" :
		res, err := svc.PullEvents(context.Background(), data.SpatioHourInterval{Hour: 9000,
			TopLeft:  &data.Point{Lat: 1, Lon: 1},
			BotRight: &data.Point{Lat: 100, Lon: 100}})
		fmt.Printf("res: %v, err: %v", res, err)
	case "pushLocations":
		city := data.City{Title:"New York", ID:"nyc"}
		locations := []data.Location{
			{Title:"loc1", ID:"5", Slug:"slug1", Position: &data.Point{Lat:1, Lon:1}},
			//{Title:"loc2", ID:"2", Slug:"slug2", Position: &data.Point{Lat:2, Lon:2}},
			//{Title:"loc4", ID:"4", Slug:"slug4", Position: &data.Point{Lat:3, Lon:3}},
		}
		err := svc.PushLocations(context.Background(),  city, locations)
		fmt.Print(err)
	case "pullLocations":
		res, err := svc.PullLocations(context.Background(), "nyc")
		fmt.Printf("res: %v, err: %v", res, err)
	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}

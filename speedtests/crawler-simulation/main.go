package main

import (
	"context"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	prand "github.com/angrymuskrat/event-monitoring-system/utils/rand/positional"
	"google.golang.org/grpc"
	"os"
	"time"
)

func main() {
	var (
		svc        storagesvc.Service
		postAmount int
		err        error
	)
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(storagesvc.MaxMsgSize)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	svc = storagesvc.NewGRPCClient(conn)
	postAmount = 300000
	randomizer := prand.New(data.Point{Lat: 60.115617, Lon: 30.103768}, data.Point{Lat: 59.738057, Lon: 30.637967})
	posts := randomizer.Posts(postAmount, 1514764800, 1577836800)
	fmt.Println("posts generated!")

	start := time.Now()
	err = svc.PushPosts(context.Background(), "spb_empty_test", posts)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("spb_empty_test completed!")
	timePerPostEmpty := time.Since(start).Microseconds() / int64(postAmount)

	start = time.Now()
	err = svc.PushPosts(context.Background(), "spb_test", posts)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("spb_test completed!")
	timePerPost := time.Since(start).Microseconds() / int64(postAmount)
	fmt.Printf("avg time per post for: \n    empty db: %v ms \n    db with real posts: %v ms", timePerPostEmpty, timePerPost)
}

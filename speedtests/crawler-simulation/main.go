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
	postAmount = 5
	randomizer := prand.New(data.Point{Lat: 60.115617, Lon: 30.103768}, data.Point{Lat: 59.738057, Lon: 30.637967})
	nExps := 50
	ans := int64(0)
	for i := 0; i < nExps; i++ {
		fmt.Printf("\n\nIteration: %v\n", i)

		posts := randomizer.Posts(postAmount, 1578836800, 1578836800+(3600*24))
		fmt.Println("posts generated!")

		start := time.Now()
		err = svc.PushPosts(context.Background(), "spb_test", posts)
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Println("spb_test completed!")
		timePerPost := time.Since(start).Microseconds() / int64(postAmount)
		fmt.Printf("avg time per post for: \n %v mcs", timePerPost)
		ans += timePerPost
	}
	fmt.Printf("avg time per post for: \n %v mcs", ans/int64(nExps))
}

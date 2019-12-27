package main

import (
	"context"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"google.golang.org/grpc"
	"os"
	"time"
)

func main() {
	var (
		svc storagesvc.Service
		err error
	)
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	svc = storagesvc.NewGRPCClient(conn)

	testPosts := GeneratePosts(1)
	testPosts[0].ID = "dHirNwnQr"
	method := "pullGrid"

	switch method {
	case "pushPosts":
		res, err := svc.PushPosts(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right\n%v", res)

	case "select":
		res, err := svc.SelectPosts(context.Background(), data.SpatioTemporalInterval{ 0, 1000000, 5,
			5, 30, 30, struct{}{}, nil, 0 })
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", len(res))

	case "pushGrid":
		err := svc.PushGrid(context.Background(), "adc", []byte("tesjhfgjakgiqywyyuyyuyyyyyyuyulyuyyl,hst"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right")

	case "pullGrid":
		res, err := svc.PullGrid(context.Background(), "abc")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right\n%v", string(res))

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}


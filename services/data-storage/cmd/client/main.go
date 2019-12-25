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
	/*fs := flag.NewFlagSet("dbcient", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", ":8082", "gRPC address of addsvc")
		method   = fs.String("method", "push", "push, select")
	)
	var (
		svc storagesvc.Service
		err error
	)
	if *grpcAddr != "" {
		//conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		svc = storagesvc.NewGRPCClient(conn)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}*/
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
	method := "push"

	switch method {
	case "push":

		res, err := svc.Push(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right\n%v", res)

	case "select":
		res, err := svc.Select(context.Background(), data.SpatioTemporalInterval{ 0, 1000000, 5,
			5, 30, 30, struct{}{}, nil, 0 })
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", len(res))

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}


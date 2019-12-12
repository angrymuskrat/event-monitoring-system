package main

import (
	"context"
	"flag"
	"fmt"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"google.golang.org/grpc"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("dbcient", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", ":8082", "gRPC address of addsvc")
		method   = fs.String("method", "push", "push, select")
	)
	var (
		svc storagesvc.Service
		err error
	)
	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
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
	}

	testPosts := GeneratePosts(10)

	switch *method {
	case "push":

		res, err := svc.Push(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right\n%v", res)

	case "select":
		res, err := svc.Select(context.Background(), data.SpatioTemporalInterval{ 0, 1000000, 5,
			5, 10, 10, struct{}{}, nil, 0 })
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", len(res))

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"os"
	"text/tabwriter"
	"time"

	"google.golang.org/grpc"

	//"sourcegraph.com/sourcegraph/appdash"
	//appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"

	"github.com/go-kit/kit/log"
)

func main() {
	// The dbclient presumes no service discovery system, and expects users to
	// provide the direct address of an dbsvc. This presumption is reflected in
	// the dbcient binary and the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.
	fs := flag.NewFlagSet("dbcient", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", ":8082", "gRPC address of addsvc")
		method   = fs.String("method", "select", "push, select")
	)

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.
	var (
		svc dbsvc.Service
		err error
	)
	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		svc = dbsvc.NewGRPCClient(conn, log.NewNopLogger())
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	testPosts := GeneratePosts(4)

	switch *method {
	case "push":

		err := svc.Push(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right")

	case "select":
		res, err := svc.Select(context.Background(), data.SpatioTemporalInterval{ 0, 1000, 0,
			0, 10, 10, struct{}{}, nil, 0 })
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok %v", res)

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

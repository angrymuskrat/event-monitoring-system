package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"google.golang.org/grpc"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbtransport"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"

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
		grpcAddr     = fs.String("grpc-addr", ":8082", "gRPC address of addsvc")
		//zipkinURL    = fs.String("zipkin-url", ":9090", "Enable Zipkin tracing via HTTP reporter URL e.g. http://localhost:9411/api/v2/spans")
		//zipkinBridge = fs.Bool("zipkin-ot-bridge", false, "Use Zipkin OpenTracing bridge instead of native implementation")
		//lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		//appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
		method = fs.String("method", "select", "push, select")
	)

	/*fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])
	fmt.Print(fs.Args(), len(fs.Args()))
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}*/

	// This is a demonstration of the native Zipkin tracing client. If using
	// Zipkin this is the more idiomatic client over OpenTracing.
	/*var zipkinTracer *zipkin.Tracer
	{
		if *zipkinURL != "" {
			var (
				err         error
				hostPort    = "" // if host:port is unknown we can keep this empty
				serviceName = "dbsvc-cli"
				reporter    = zipkinhttp.NewReporter(*zipkinURL)
			)
			defer reporter.Close()
			zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to create zipkin tracer: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}*/

	// This is a demonstration client, which supports multiple tracers.
	// Your clients will probably just use one tracer.
	/*var otTracer stdopentracing.Tracer
	{
		if *zipkinBridge && zipkinTracer != nil {
			otTracer = zipkinot.Wrap(zipkinTracer)
			zipkinTracer = nil // do not instrument with both native and ot bridge
		} else {
			otTracer = stdopentracing.GlobalTracer() // no-op
		}
	}*/

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.
	var (
		svc dbservice.Service
		err error
	)
	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		svc = dbtransport.NewGRPCClient(conn/*, otTracer, zipkinTracer*/, log.NewNopLogger())
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	testPosts := []pb.Post{pb.Post{"13131", "23", "-", false, "sasda", 5, 5, 5, false, "1313", "12313", 12.4, 12.5}}

	switch *method {
	case "push":

		err := svc.Push(context.Background(), testPosts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "It is all right")

	case "select":
		_, err := svc.Select(context.Background(), pb.SpatialTemporalInterval{1, 2, 1, 2, 1, 2})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Ok")

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


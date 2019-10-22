package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	//lightstep "github.com/lightstep/lightstep-tracer-go"
	"github.com/oklog/oklog/pkg/group"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	//zipkin "github.com/openzipkin/zipkin-go"
	//zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbendpoint"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbtransport"
	dbpb "github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
)

func main() {
	// Define our flags. Your service probably won't need to bind listeners for
	// *all* supported transports, or support both Zipkin and LightStep, and so
	// on, but we do it here for demonstration purposes.
	fs := flag.NewFlagSet("dbsvc", flag.ExitOnError)
	var (
		//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
		//zipkinURL      = fs.String("zipkin-url", ":9090", "Enable Zipkin tracing via HTTP reporter URL e.g. http://localhost:9411/api/v2/spans")
		//zipkinBridge   = fs.Bool("zipkin-ot-bridge", false, "Use Zipkin OpenTracing bridge instead of native implementation")
		//lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		//appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	/*var zipkinTracer *zipkin.Tracer
	{
		if *zipkinURL != "" {
			var (
				err         error
				hostPort    = "localhost:80"
				serviceName = "dbsvc"
				reporter    = zipkinhttp.NewReporter(*zipkinURL)
			)
			defer reporter.Close()
			zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			if !(*zipkinBridge) {
				logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
			}
		}
	}*/

	// Determine which OpenTracing tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency.
	/*var tracer stdopentracing.Tracer
	{
		if *zipkinBridge && zipkinTracer != nil {
			logger.Log("tracer", "Zipkin", "type", "OpenTracing", "URL", *zipkinURL)
			tracer = zipkinot.Wrap(zipkinTracer)
			zipkinTracer = nil // do not instrument with both native tracer and opentracing bridge
		} else {
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}*/

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	var pushed, selected metrics.Counter
	{
		// Business-level metrics.
		pushed = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "event_monitoring_system",
			Subsystem: "dbsvc",
			Name:      "push_query",
			Help:      "Total count of pushed posts to DB",
		}, []string{})
		selected = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "event_monitoring_system",
			Subsystem: "dbsvc",
			Name:      "select_query",
			Help:      "Total count of Selected posts in DB",
		}, []string{})
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "event_monitoring_system",
			Subsystem: "dbsvc",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// Build the layers of the service "onion" from the inside out. First, the
	// business logic service; then, the set of endpoints that wrap the service;
	// and finally, a series of concrete transport adapters. The adapters, like
	// the HTTP handler or the gRPC server, are the bridge between Go kit and
	// the interfaces that the transports expect. Note that we're not binding
	// them to ports or anything yet; we'll do that next.
	var (
		service        = dbservice.New(logger, pushed, selected)
		endpoints      = dbendpoint.New(service, logger, duration/*, tracer, zipkinTracer*/)
		grpcServer     = dbtransport.NewGRPCServer(endpoints/*, tracer, zipkinTracer*/, logger)
	)

	// Now we're to the part of the func main where we want to start actually
	// running things, like servers bound to listeners to receive connections.
	//
	// The method is the same for each component: add a new actor to the group
	// struct, which is a combination of 2 anonymous functions: the first
	// function actually runs the component, and the second function should
	// interrupt the first function and cause it to return. It's in these
	// functions that we actually bind the Go kit server/handler structs to the
	// concrete transports and run them.
	//
	// Putting each component into its own block is mostly for aesthetics: it
	// clearly demarcates the scope in which each listener/socket may be used.
	var g group.Group
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", *grpcAddr)
			// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by
			// the here demonstrated zipkin tracing middleware.
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			dbpb.RegisterDbsvcServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
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
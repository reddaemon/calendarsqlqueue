package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"

	"github.com/reddaemon/calendarsqlqueue/internal/database/postgres"

	"flag"

	"context"

	grpc_prometh "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/reddaemon/calendarsqlqueue/config"
	db2 "github.com/reddaemon/calendarsqlqueue/db"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/server"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/service"
	"github.com/reddaemon/calendarsqlqueue/logger"
	eventpb "github.com/reddaemon/calendarsqlqueue/protofiles/protofiles/api"
	"google.golang.org/grpc"
)

var (
	grpcRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "calendar_grpc_requests_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"method", "code"})

	latencies = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "calendar_grpc_request_latency_seconds",
		Help:       "The temperature of the frog pond.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.97: 0.001, 0.99: 0.001},
	}, []string{"method", "code"})
)

func metricsInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	handlingBeginsAt := time.Now()
	code := codes.OK

	res, err := handler(ctx, request)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			code = st.Code()
		} else {
			code = codes.Unknown
		}
	}

	labels := prometheus.Labels{
		"method": info.FullMethod,
		"code":   code.String(),
	}
	latencies.With(labels).Observe(time.Since(handlingBeginsAt).Seconds())
	grpcRequests.With(labels).Inc()

	return res, err
}

func init() {
	prometheus.MustRegister(grpcRequests)
	prometheus.MustRegister(latencies)
}

func main() {
	go runPrometheusServer()

	configPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()
	c, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	l, err := logger.GetLogger(c)
	if err != nil {
		log.Fatalf("unable to load logger: %v", err)
	}

	db, err := db2.GetDb(c)

	if err != nil {
		log.Fatalf("unable to load db: %v", err)
	}

	lis, err := net.Listen("tcp", c.Host+":"+c.Port)
	if err != nil {
		l.Fatal(fmt.Sprintf("failed to listen %v", err))
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metricsInterceptor),
	)

	var contextTimeout time.Duration
	contextTimeout = time.Millisecond * 500
	repo := postgres.NewPsqlRepository(db)
	ucs := service.NewEventUsecase(repo, contextTimeout)
	eventpb.RegisterEventServiceServer(grpcServer, server.NewServer(ucs, l))

	grpc_prometh.Register(grpcServer)
	grpc_prometh.EnableHandlingTimeHistogram()
	err = grpcServer.Serve(lis)

	if err != nil {
		l.Fatal(err.Error())
	}
}

func runPrometheusServer() {
	log.Println("run prometheus exporter server")
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":9187", nil)
}

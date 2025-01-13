package main

import (
	"addsrv/pb"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	httpAddr = flag.String("http-addr", ":8080", "HTTP listen address")
	grpcAddr = flag.String("grpc-addr", ":8972", "gRPC listen address")
)

/*
go-kit分为三层
1. 传输层(Transport layer)
2. 端点层(Endpoint layer)
3. 服务层(Service layer)
请求在第一层进入服务，向下流到第三层，响应则相反
*/

func main() {

	// instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	logger := log.NewLogfmtLogger(os.Stderr)
	bs := NewService()
	bs = NewLogMiddleware(logger, bs)
	bs = instrumentingMiddleware{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		countResult:    countResult,
		next:           bs,
	}
	
	var g errgroup.Group

	// HTTP服务
	g.Go(func() error {
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			fmt.Printf("http: net.Listen(tcp, %s) failed, err:%v\n", *httpAddr, err)
			return err
		}
		defer httpListener.Close()
		logger := log.NewLogfmtLogger(os.Stderr)
		httpHandler := NewHTTPServer(bs, logger)
		return http.Serve(httpListener, httpHandler)
	})

	// GRPC服务
	g.Go(func() error {
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			fmt.Printf("grpc: net.Listen(tcp, %s) faield, err:%v\n", *grpcAddr, err)
			return err
		}
		defer grpcListener.Close()
		s := grpc.NewServer()
		pb.RegisterAddServer(s, NewGRPCServer(bs))
		return s.Serve(grpcListener)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("server exit with err:%v\n", err)
	}
}

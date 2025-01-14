package main

import (
	"addsrv/pb"
	"context"
	"encoding/json"
	"net/http"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

// ----------------------------Transports-------------------------------
// 传输域绑定到具体的传输协议，如 HTTP 或 gRPC。在一个微服务可能支持一个或多个传输协议的世界中
// 这是非常强大的：你可以在单个微服务中支持原有的 HTTP API 和新增的 RPC 服务

// gRPC的请求与响应
// decodeGRPCSumRequest 将Sum方法的gRPC请求参数转为内部的SumRequest
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return SumRequest{A: int(req.A), B: int(req.B)}, nil
}

// decodeGRPCConcatRequest 将Concat方法的gRPC请求参数转为内部的ConcatRequest
func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ConcatRequest)
	return ConcatRequest{A: req.A, B: req.B}, nil
}

// encodeGRPCSumResponse 封装Sum的gRPC响应
func encodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(SumResponse)
	return &pb.SumResponse{V: int64(resp.V), Err: resp.Err}, nil
}

// encodeGRPCConcatResponse 封装Concat的gRPC响应
func encodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ConcatResponse)
	return &pb.ConcatResponse{V: resp.V, Err: resp.Err}, nil
}

// gRPC
type grpcServer struct {
	pb.UnimplementedAddServer

	sum    grpctransport.Handler
	concat grpctransport.Handler
}

func (s grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	_, resp, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SumResponse), nil
}

func (s grpcServer) Concat(ctx context.Context, req *pb.ConcatRequest) (*pb.ConcatResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ConcatResponse), nil
}

// NewGRPCServer 构造函数
func NewGRPCServer(svc AddService) pb.AddServer {
	return &grpcServer{
		sum: grpctransport.NewServer(
			makeSumEndpoint(svc), // endpoint
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
		),
		concat: grpctransport.NewServer(
			makeConcatEndpoint(svc),
			decodeGRPCConcatRequest,
			encodeGRPCConcatResponse,
		),
	}
}

// HTTP
func decodeSumRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeConcatRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// encodeTrimRequest 将内部结构体转为protobuf中的结构体
// 对外发起gRPC请求
func encodeTrimRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(trimRequest)
	return &pb.TrimRequest{S: req.s}, nil
}

// decodeTrimResponse 将收到的gRPC响应转为内部的响应结构体
func decodeTrimResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.TrimResponse)
	return trimResponse{s: resp.S}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// transport层日志
// HTTP Server
func NewHTTPServer(svc AddService, logger log.Logger) http.Handler {
	// 加入日志中间件
	sum := makeSumEndpoint(svc)
	sum = loggingMiddleware(log.With(logger, "method", "sum"))(sum)
	sum = rateMiddleware(rate.NewLimiter(1, 1))(sum)
	sumHandler := httptransport.NewServer(
		sum,
		decodeSumRequest,
		encodeResponse,
	)

	concat := makeConcatEndpoint(svc)
	concat = loggingMiddleware(log.With(logger, "method", "concat"))(concat)
	concatHandler := httptransport.NewServer(
		concat,
		decodeConcatRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/sum", sumHandler).Methods("POST")
	r.Handle("/concat", concatHandler).Methods("POST")
	http.Handle("/metrics", promhttp.Handler())
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	return r
}

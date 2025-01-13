package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// --------------------------endpoint-----------------------------
// 端点就像控制器上的动作/处理程序; 它是安全性和抗脆弱性逻辑的所在。
// 如果实现两种传输(HTTP 和 gRPC) ，则可能有两种将请求发送到同一端点的方法
/*
	type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
	它表示单个 RPC。也就是说，我们的服务接口中只有一个方法。我们将编写简单的 '适配器' 来将服务的每个方法转换为一个端点。
	每个适配器接受一个 AddService，并返回与其中一个方法对应的端点
*/
type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumResponse struct {
	V   int    `json:"v"`
	Err string `json:"err,omitempty"`
}

type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type ConcatResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

func makeSumEndpoint(srv AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		v, err := srv.Sum(ctx, req.A, req.B) // 方法调用
		if err != nil {
			return SumResponse{V: v, Err: err.Error()}, nil
		}
		return SumResponse{V: v}, nil
	}
}

func makeConcatEndpoint(srv AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatRequest)
		v, err := srv.Concat(ctx, req.A, req.B) // 方法调用
		if err != nil {
			return ConcatResponse{V: v, Err: err.Error()}, nil
		}
		return ConcatResponse{V: v}, nil
	}
}

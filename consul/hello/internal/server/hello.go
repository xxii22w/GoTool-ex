package server

import (
	"context"
	"fmt"
	"hello-ex/internal/service"
	"net/http"
)

type HelloSrv struct {
	service.UnimplementedHelloServer
}

func (h *HelloSrv) SayHello(_ context.Context,request *service.HelloRequest) (*service.HelloResponse,error) {
	// tokenString := request.Username
	// claims,err := token.ParseToken(tokenString)
	// if err != nil {
	// 	return &service.HelloResponse{
	// 		Code: http.StatusInternalServerError,
	// 	},nil
	// }
	username := request.Username
	return &service.HelloResponse{
		Code: http.StatusOK,
		Msg: "hello",
		// Data: fmt.Sprintf("hello %v.",claims.Username),
		Data: fmt.Sprintf("hello %s.",username),
	},nil
}
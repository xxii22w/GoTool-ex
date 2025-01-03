package api

import (
	"context"
	"gateway/center"
	"gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Manager) RouteHello() {
	m.handler.POST("/hello", m.hello)
}

type helloMsg struct {
	Username string `json:"username"`
}

func (m *Manager) hello(ctx *gin.Context) {
	var msg helloMsg
	// 获取请求中的json数据
	if err := ctx.ShouldBindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// center.Resolver() 参数为调用的服务名
	// 该函数会进行自动负载均衡并返回一个*grpc.ClientConn
	conn, err := center.Resolver("hello")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	c := service.NewHelloClient(conn)
	response, _ := c.SayHello(context.Background(), &service.HelloRequest{
		Username: msg.Username,
	})

	ctx.JSON(int(response.Code), gin.H{"msg": response.Msg, "data": map[string]any{
		"data": response.Data,
	}})
}
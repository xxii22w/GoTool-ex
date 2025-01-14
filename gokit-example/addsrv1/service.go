package main

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	sdconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// ---------------------------Services---------------------------------
// Go kit 服务层应该努力遵守整洁架构或六边形架构
// 服务（指Go kit中的service层）是实现所有业务逻辑的地方。服务层通常将多个端点粘合在一起。
// 在 Go kit 中，服务层通常被抽象为接口，这些接口的实现包含业务逻辑。
// AddService 把两个东西加在一起

type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

// addService 一个AddService接口的具体实现
// 它的内部可以按需添加各种字段
type addService struct {
	// db db.Conn
	// logger zap.Logger
}

var (
	// ErrEmptyString 两个参数都是空字符串的错误
	ErrEmptyString = errors.New("两个参数都是空字符串")
	ErrRateLimit   = errors.New("request rate limit")
)

// Sum 返回两个数的和
func (addService) Sum(_ context.Context, a, b int) (int, error) {
	// 业务逻辑
	return a + b, nil
}

// Concat 拼接两个字符串
func (addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrEmptyString
	}
	return a + b, nil
}

// NewService 创建一个add service
func NewService() AddService {
	return &addService{}
}

// consul
// 从注册中心获取trim服务的地址
// 基于consul实现对trim service的服务发现
func getTrimServiceFromConsul(consulAddr string, logger log.Logger, srvName string, tags []string) (endpoint.Endpoint, error) {
	// 1. 连consul
	cfg := consulapi.DefaultConfig()
	cfg.Address = consulAddr
	cc, err := consulapi.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// 2. 使用go kit 提供的适配器
	sdClient := sdconsul.NewClient(cc)

	instancer := sdconsul.NewInstancer(sdClient, logger, srvName, tags, true)
	// 3. Endpointer
	// Go kit 为不同的服务发现系统（eureka、zookeeper、consul、etcd等）提供适配器
	// Endpointer负责监听服务发现系统，并根据需要生成一组相同的端点
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	// 4. Balancer
	balancer := lb.NewRoundRobin(endpointer)
	// 5. retry
	// 重试策略包装负载均衡器，并返回可用的端点。重试策略将重试失败的请求，直到达到最大尝试或超时为止
	retry := lb.Retry(3, time.Second, balancer)
	return retry, nil
}

// 将实例字符串(例如host:port)转换为特定端点的函数。提供多个端点的实例需要多个工厂函数
// 工厂函数还返回一个当实例消失并需要清理时调用的io.Closer
func factory(instance string) (endpoint.Endpoint, io.Closer, error) {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	e := makeTrimEndpoint(conn)
	return e, conn, err
}
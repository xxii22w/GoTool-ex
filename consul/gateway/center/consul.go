package center

import (
	"fmt"
	"log"

	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var (
	addr = "127.0.0.1:8500"
	client *consul.Client
	dc     = "dc1"           // 数据中心，根据实际情况修改
)

func init() {
	var err error 
	config := consul.DefaultConfig()
	config.Address = addr
	client,err = consul.NewClient(config)
	if err != nil {
		log.Fatalf("failed to init consul: %v", err)
	}
}

// register consul注册服务
func Register(reg consul.AgentServiceRegistration) error {
	agent := client.Agent()
	return agent.ServiceRegister(&reg)
}

// resolver对服务进行负载均衡
func Resolver(name string) (*grpc.ClientConn,error) {
	 // 使用Consul服务发现
	 query := &consul.QueryOptions{
        Datacenter: dc,
    }
    entries, _, err := client.Health().Service(name, "", true, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query Consul for service %s: %w", name, err)
    }
    if len(entries) == 0 {
        return nil, fmt.Errorf("no instances of service %s found", name)
    }

    // 构建gRPC连接地址列表
    var serverAddresses []string
    for _, entry := range entries {
        serverAddresses = append(serverAddresses, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
    }

    // 创建轮询负载均衡的gRPC连接
    conn, err := grpc.Dial(
        serverAddresses[0], // 选择第一个服务实例进行连接
        grpc.WithInsecure(), // 禁用TLS，根据实际情况决定是否需要启用
        grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
    )
    if err != nil {
        return nil, fmt.Errorf("did not connect: %w", err)
    }
    return conn, nil
}
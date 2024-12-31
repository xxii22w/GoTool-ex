package center

import (
	"log"

	consul "github.com/hashicorp/consul/api"
)

var (
	addr = "127.0.0.1:8500"
	client *consul.Client
)

func init() {
	var err error 
	config := consul.DefaultConfig()
	config.Address = addr
	client,err = consul.NewClient(config)
	if err != nil {
		log.Fatalf("failed to init consul: %v",err)
	}
}

func Register(reg consul.AgentServiceRegistration) error {
	agent := client.Agent()
	return agent.ServiceRegister(&reg)
}
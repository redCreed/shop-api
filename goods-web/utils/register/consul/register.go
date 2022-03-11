package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type RegistryClient interface {
	//注册服务
	Register(address string, port int, name string, tag []string, id string)
	//注销服务
	DeRegister(id string) error
}

type Registry struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewRegistry(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(address string, port int, name string, tag []string, id string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	check := &api.AgentServiceCheck{
		Interval:                       "60s",
		Timeout:                        "5s",
		HTTP:                           fmt.Sprintf(`http://%s:%d/health`, address, port),
		DeregisterCriticalServiceAfter: "15s",
	}
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tag,
		Port:    port,
		Address: address,
		Check:   check,
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		panic(err)
	}
}

func (r *Registry) DeRegister(id string) error{
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	//id为注册服务id
	return client.Agent().ServiceDeregister(id)
}

package rpc

import (
	consulAPI "github.com/hashicorp/consul/api"
	"mxshop/app/mxshop/api/internal/data"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/registry/consul"
)

func NewDiscovery(opts *options.RegistryOptions) registry.Discovery {
	//客户端 服务发现
	c := consulAPI.DefaultConfig()
	c.Address = opts.Address
	c.Scheme = opts.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(any(err))
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

// rpc的连接，基于服务发现
func GetDataFactoryOr(options *options.RegistryOptions) (data.UserData, error) {
	//这里负责依赖的所有的rpc连接
	discovery := NewDiscovery(options)
	userClient := NewUserServiceClient(discovery)
	return NewUsers(userClient), nil

}

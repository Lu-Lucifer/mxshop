package rpc

import (
	"fmt"
	consulAPI "github.com/hashicorp/consul/api"
	gpb "mxshop/api/goods/v1"
	upb "mxshop/api/user/v1"
	"mxshop/app/mxshop/api/internal/data"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/registry/consul"
	"mxshop/pkg/errors"
	"sync"
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

var (
	dbFactory data.DataFactory
	once      sync.Once
)

// rpc的连接，基于服务发现
func GetDataFactoryOr(options *options.RegistryOptions) (data.DataFactory, error) {
	//这里负责依赖的所有的rpc连接
	if options == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get grpc store factory")
	}
	once.Do(func() {
		discovery := NewDiscovery(options)
		userClient := NewUserServiceClient(discovery)
		goodsClient := NewGoodsServiceClient(discovery)
		dbFactory = &grpcData{
			gc: goodsClient,
			uc: userClient,
		}
	})
	if dbFactory == nil {
		return nil, errors.WithCode(code.ErrConnectGRPC, "failed to get grpc factory")
	}
	return dbFactory, nil
}

type grpcData struct {
	gc gpb.GoodsClient
	uc upb.UserClient
}

func (g *grpcData) Goods() gpb.GoodsClient {
	return g.gc
}

func (g *grpcData) User() data.UserData {

	return NewUsers(g.uc)
}

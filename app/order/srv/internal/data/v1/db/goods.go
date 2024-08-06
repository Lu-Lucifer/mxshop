package db

import (
	"context"
	consulAPI "github.com/hashicorp/consul/api"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/registry/consul"

	gpbv1 "mxshop/api/goods/v1"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/gmicro/server/rpcserver/clientinterceptors"
)

const goodsServiceName = "discovery:///mxshop-goods-srv"

// 目前是基于consul实现的  以后想换成nocos etcd等  可以直接在这换
func NewDiscovery(opts *options.RegistryOptions) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = opts.Address
	c.Scheme = opts.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}
func GetGoodsClient(opts *options.RegistryOptions) gpbv1.GoodsClient {
	discovery := NewDiscovery(opts)
	//这里负责依赖所有的rpc连接
	goodsClient := NewGoodsServiceClient(discovery)
	return goodsClient
}
func NewGoodsServiceClient(r registry.Discovery) gpbv1.GoodsClient {

	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(goodsServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := gpbv1.NewGoodsClient(conn)
	return c
}

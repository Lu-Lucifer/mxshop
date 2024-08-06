package db

import (
	"context"
	proto "mxshop/api/inventory/v1"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/gmicro/server/rpcserver/clientinterceptors"
)

const invServiceName = "discovery:///mxshop-inventory-srv"

func GetInventoryClient(opts *options.RegistryOptions) proto.InventoryClient {
	discovery := NewDiscovery(opts)
	//这里负责依赖所有的rpc连接
	inventoryClient := NewInventoryServiceClient(discovery)
	return inventoryClient
}

func NewInventoryServiceClient(r registry.Discovery) proto.InventoryClient {

	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(invServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := proto.NewInventoryClient(conn)
	return c
}

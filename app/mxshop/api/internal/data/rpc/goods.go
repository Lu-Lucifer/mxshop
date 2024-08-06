package rpc

import (
	"context"
	gpb "mxshop/api/goods/v1"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/gmicro/server/rpcserver/clientinterceptors"
)

const (
	goodsServiceName = "discovery:///mxshop-goods-srv"
)

func NewGoodsServiceClient(r registry.Discovery) gpb.GoodsClient {
	conn, err := rpcserver.DialInsecure(context.Background(),
		rpcserver.WithEndpoint(goodsServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(any(err))
	}
	c := gpb.NewGoodsClient(conn)
	return c
}

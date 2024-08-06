package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	v1 "mxshop/api/goods/v1"
	"mxshop/gmicro/registry/consul"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/gmicro/server/rpcserver/selector"
	"mxshop/gmicro/server/rpcserver/selector/random"
)

func main() {
	//设置全局负载均衡算法
	selector.SetGlobalSelector(random.NewBuilder())
	rpcserver.InitBuilder()
	//客户端 服务发现
	c := api.DefaultConfig()
	c.Address = "127.0.0.1:8500"
	c.Scheme = "http"
	cli, err := api.NewClient(c)
	if err != nil {
		panic(any(err))
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	conn, err := rpcserver.DialInsecure(context.Background(),
		rpcserver.WithEndpoint("discovery:///mxshop-goods-srv"),
		//rpcserver.WithEndpoint("192.168.0.101:8078"),
		rpcserver.WithDiscovery(r),
		rpcserver.WithBalancerName("selector"),
	)
	if err != nil {
		panic(any(err))
	}
	defer conn.Close()

	uc := v1.NewGoodsClient(conn)

	//rsp, _ := uc.GoodsList(context.Background(), &v1.GoodsFilterRequest{
	//	KeyWords: "猕猴桃",
	//})
	rsp, _ := uc.BatchGetGoods(context.Background(), &v1.BatchGoodsIdInfo{
		Id: []int32{421, 422},
	})
	fmt.Println(rsp)

}

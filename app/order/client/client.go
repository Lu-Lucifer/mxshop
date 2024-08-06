package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
	v1 "mxshop/api/order/v1"
	"mxshop/gmicro/registry/consul"
	"mxshop/gmicro/server/rpcserver"
	// _ "mxshop/gmicro/server/rpcserver/resolver/direct" // 这个是直接连接的 下面已经实现watcher长轮询了  弃用
	"mxshop/gmicro/server/rpcserver/selector"
	"mxshop/gmicro/server/rpcserver/selector/random"
	"time"
)

func generateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}
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
		rpcserver.WithEndpoint("discovery:///mxshop-order-srv"),
		rpcserver.WithDiscovery(r),
		rpcserver.WithBalancerName("selector"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	uc := v1.NewOrderClient(conn)
	_, err = uc.SubmitOrder(context.Background(), &v1.OrderRequest{
		UserId:  1,
		Address: "深圳新东方",
		OrderSn: generateOrderSn(12),
		Name:    "haha",
		Post:    "快点",
		Mobile:  "15963659889",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("订单新建成功")
}

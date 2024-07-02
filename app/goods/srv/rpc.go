package srv

import (
	"fmt"
	gpb "mxshop/api/goods/v1"
	"mxshop/app/goods/srv/config"
	v1 "mxshop/app/goods/srv/internal/controller/v1"
	"mxshop/app/goods/srv/internal/data/v1/db"
	"mxshop/app/goods/srv/internal/data_search/v1/es"
	srv1 "mxshop/app/goods/srv/internal/service/v1"
	"mxshop/gmicro/core/trace"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/pkg/log"
)

func NewGoodsRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化opentelemetry的exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})
	//初始化db，esClient
	dbFactory, err := db.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	seatchFactory, err := es.GetSeatchFactoryOr(cfg.EsOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	//构造一个goodsServer结构体对象（十分繁琐，使用工厂模式来解决；或ioc依赖注入）
	//data := db.NewGoods(gormDB)
	//categoryData := db.NewCategory(gormDB)
	//searchData := es.NewGoods(seatchClient)
	//brandData := db.NewBrands(gormDB)
	//srv := srv1.NewGoodsService(data, categoryData, searchData, brandData)
	//goodsServer := v1.NewGoodsServer(srv)

	//使用工厂模式来解决
	//srv := srv1.NewGoodsService(dbFactory, seatchFactory)
	srvFactory := srv1.NewService(dbFactory, seatchFactory)
	goodsServer := v1.NewGoodsServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	gpb.RegisterGoodsServer(urpcServer.Server, goodsServer)
	return urpcServer, nil
}

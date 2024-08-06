package srv

import (
	"fmt"
	gpb "mxshop/api/inventory/v1"
	"mxshop/app/inventory/srv/config"
	v1 "mxshop/app/inventory/srv/internal/controller/v1"
	"mxshop/app/inventory/srv/internal/data/v1/db"
	srv1 "mxshop/app/inventory/srv/internal/service/v1"
	"mxshop/gmicro/core/trace"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/pkg/log"
)

func NewInventoryRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化opentelemetry的exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})
	//初始化db
	dbFactory, err := mysql.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	srvFactory := srv1.NewService(dbFactory, cfg.RedisOptions)
	invServer := v1.NewInventoryServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	gpb.RegisterInventoryServer(grpcServer.Server, invServer)
	return grpcServer, nil
}

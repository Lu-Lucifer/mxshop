package srv

import (
	"fmt"
	upb "mxshop/api/user/v1"
	"mxshop/app/user/srv/config"
	"mxshop/app/user/srv/controller/user"
	"mxshop/app/user/srv/data/v1/db"
	srv1 "mxshop/app/user/srv/service/v1"
	"mxshop/gmicro/core/trace"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/pkg/log"
)

func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化opentelemetry的exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})
	//初始化db
	gormDB, err := db.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	//构造一个userServer结构体对象
	data := db.NewUsers(gormDB)
	srv := srv1.NewUserService(data)
	userver := user.NewUserServer(srv)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	upb.RegisterUserServer(urpcServer.Server, userver)
	return urpcServer, nil
}

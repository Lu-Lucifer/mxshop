package api

import (
	"mxshop/app/mxshop/api/config"
	"mxshop/gmicro/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	aRestServer := restserver.NewServer(restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddlewares(cfg.Server.Middlewares),
	)

	//配置路由
	initRouter(aRestServer, cfg)
	return aRestServer, nil

}

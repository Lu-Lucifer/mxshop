package srv

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"mxshop/app/inventory/srv/config"
	"mxshop/app/pkg/options"
	gapp "mxshop/gmicro/app"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/registry/consul"
	"mxshop/pkg/app"
	"mxshop/pkg/log"
	"mxshop/pkg/storage"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("inventory", "mxshop",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		//app.WithNoConfig(), //设置不从配置文件中读取
	)
	return appl
}

//func run(basename string) error {
//	fmt.Println("Starting")
//	return nil
//}

// 上面run函数进一步包装，传入参数config，打印项目启动后的日志等级
func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		fmt.Println(basename)
		inventoryApp, err := NewInventoryApp(cfg)
		if err != nil {
			return err
		}
		//启动
		if err := inventoryApp.Run(); err != nil {
			log.Errorf("run inventory app error: %s", err)
		}
		return nil
	}
}

// 创建启动user app, 使用启动grpc服务
func NewInventoryApp(cfg *config.Config) (*gapp.App, error) {
	//初始化log
	log.Init(cfg.Log)
	defer log.Flush()

	//服务注册
	register := NewRegistrar(cfg.Registry)

	//连接redis
	redisConfig := &storage.Config{
		Host:                  cfg.RedisOptions.Host,
		Port:                  cfg.RedisOptions.Port,
		Addrs:                 cfg.RedisOptions.Addrs,
		MasterName:            cfg.RedisOptions.MasterName,
		Username:              cfg.RedisOptions.Username,
		Password:              cfg.RedisOptions.Password,
		Database:              cfg.RedisOptions.Database,
		MaxIdle:               cfg.RedisOptions.MaxIdle,
		MaxActive:             cfg.RedisOptions.MaxActive,
		Timeout:               cfg.RedisOptions.Timeout,
		EnableCluster:         cfg.RedisOptions.EnableCluster,
		UseSSL:                cfg.RedisOptions.UseSSL,
		SSLInsecureSkipVerify: cfg.RedisOptions.SSLInsecureSkipVerify,
		EnableTracing:         cfg.RedisOptions.EnableTracing,
	}
	go storage.ConnectToRedis(context.Background(), redisConfig)

	//生成rpc服务
	rpcServer, err := NewInventoryRPCServer(cfg)
	if err != nil {
		return nil, err
	}
	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRPCServer(rpcServer),
		gapp.WithRegistrar(register),
	), nil

}

func NewRegistrar(registry *options.RegistryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(any(err))
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

package admin

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"mxshop/app/pkg/options"
	"mxshop/app/user/srv/config"
	gapp "mxshop/gmicro/app"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/registry/consul"
	"mxshop/pkg/app"
	"mxshop/pkg/log"
)

// 启动rest server项目示例
func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("user", "mxshop",
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
		fmt.Println(cfg.Log.Level)
		userApp, err := NewUserApp(cfg)
		if err != nil {
			return err
		}
		//启动
		if err := userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
		}
		return nil
	}
}

// 创建启动user app, 使用启动rest服务
func NewUserApp(cfg *config.Config) (*gapp.App, error) {
	//初始化log
	log.Init(cfg.Log)
	defer log.Flush()

	//服务注册
	register := NewRegistrar(cfg.Registry)
	//生成rpc服务
	restServer, err := NewUserHTTPServer(cfg)
	if err != nil {
		return nil, err
	}
	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRestServer(restServer),
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

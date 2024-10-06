package srv

import (
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"go.opentelemetry.io/otel"
	upb "mxshop/api/user/v1"
	"mxshop/app/pkg/options"
	"mxshop/app/user/srv/config"
	"mxshop/app/user/srv/controller/user"
	"mxshop/app/user/srv/data/v1/db"
	srv1 "mxshop/app/user/srv/service/v1"
	"mxshop/gmicro/core/trace"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/pkg/log"
	"time"
)

func NewNacosDataSource(opts *options.NacosOptions) (*nacos.NacosDataSource, error) {
	sc := []constant.ServerConfig{
		{
			ContextPath: "/nacos",
			IpAddr:      opts.Host,
			Port:        opts.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId: opts.NamespaceId,
		TimeoutMs:   5000,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, err
	}
	//注册流控规则handler
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	//创建nacosdatasource数据源
	nds, err := nacos.NewNacosDataSource(client, opts.Group, opts.DataId, h)
	if err != nil {
		return nil, err
	}
	return nds, nil
}

func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化opentelemetry的exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	//telemetry使用
	tr := otel.Tracer("mxshop-user-service")
	spanCtx, span := tr.Start(context.Background(), "GetDBFactoryOr")
	//初始化db

	gormDB, err := db.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	span.End()

	//构造一个userServer结构体对象
	data := db.NewUsers(gormDB)
	srv := srv1.NewUserService(data)
	userver := user.NewUserServer(srv)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	var opts []rpcserver.ServerOption
	opts = append(opts, rpcserver.WithAddress(rpcAddr))
	//是否开启限流，动态组装
	if cfg.Server.EnableLimit {
		opts = append(opts, rpcserver.WithUnaryInterceptor(grpc.NewUnaryServerInterceptor()))
		nds, err := NewNacosDataSource(cfg.Nacos)
		if err != nil {
			return nil, err
		}
		err = nds.Initialize()
		if err != nil {
			return nil, err
		}
	}
	//telemetry使用
	_, newServerspan := tr.Start(spanCtx, "NewServer")
	urpcServer := rpcserver.NewServer(opts...)
	upb.RegisterUserServer(urpcServer.Server, userver)
	newServerspan.End()

	//telemetry使用
	_, mockServerSpan := tr.Start(spanCtx, "mockServer")
	time.Sleep(time.Millisecond * 300)
	mockServerSpan.End()
	return urpcServer, nil
}

package rpcserver

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	apimd "mxshop/api/metadata"
	srvintc "mxshop/gmicro/server/rpcserver/serverinterceptors"
	"mxshop/pkg/host"
	"mxshop/pkg/log"
	"net"
	"net/url"
	"time"
)

type Server struct {
	*grpc.Server
	address    string
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	lis        net.Listener

	health   *health.Server
	endpoint *url.URL
	metadata *apimd.Server
	timeout  time.Duration //服务端超时
}

func (s *Server) Address() string {
	return s.address
}
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0",
		health:  health.NewServer(),
		//timeout:  time.Second * 1,
	}
	for _, opt := range opts {
		opt(srv)
	}
	//TODO 用户不设置拦截器时，默认加上一些拦截器 crash，tracing
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryRecoverInterceptor,
		//srvintc.UnaryTimeoutInterceptor(srv.timeout),
		//opentelemetry
		otelgrpc.UnaryServerInterceptor(),
	}
	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	//把我们传入的拦截器转换成grpc的ServerOption
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...)}
	//把用户自己传入的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	//srv.Server为内嵌的grpc.server
	srv.Server = grpc.NewServer(grpcOpts...)
	//注册metadata的Server
	srv.metadata = apimd.NewServer(srv.Server)
	//自动解析address和设置listener
	err := srv.listenAndEndpoint()
	if err != nil {
		return nil
	}

	//注册health
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	//可以支持用户直接通过grpc的一个接口查看当前支持的所有rpc服务(使用kratos中的api
	//中的metadata)
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)
	return srv
}

type ServerOption func(s *Server)

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func WithLis(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithUnaryInterceptor(ints ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = ints
	}
}

func WithStreamInterceptor(ints ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = ints
	}
}

func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// 完成ip和端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

func (s *Server) Start(ctx context.Context) error {
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	//改grpc核心变量 状态
	//只有.Resume()之后，请求才能进来
	//s.health.Shutdown()相反
	s.health.Resume()
	return s.Server.Serve(s.lis)

}
func (s *Server) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving 防止接受新的请求
	s.health.Shutdown()
	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}

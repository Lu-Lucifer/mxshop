package rpcserver

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	"log"
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/rpcserver/clientinterceptors"
	"mxshop/gmicro/server/rpcserver/resolver/discovery"
	"time"
)

type clientOptions struct {
	endpoint      string
	timeout       time.Duration
	discovery     registry.Discovery
	unaryInts     []grpc.UnaryClientInterceptor
	streamInts    []grpc.StreamClientInterceptor
	rpcOpts       []grpc.DialOption
	balancerName  string
	logger        *log.Logger
	enableTracing bool
}

type ClientOption func(c *clientOptions)

func WithEndpoint(endpoint string) ClientOption {
	return func(c *clientOptions) {
		c.endpoint = endpoint
	}
}

func WithEnableTracing(ok bool) ClientOption {
	return func(c *clientOptions) {
		c.enableTracing = ok
	}
}

func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.timeout = timeout
	}
}

func WithDiscovery(discovery registry.Discovery) ClientOption {
	return func(c *clientOptions) {
		c.discovery = discovery
	}
}

func WithClientUnaryInterceptor(ints ...grpc.UnaryClientInterceptor) ClientOption {
	return func(c *clientOptions) {
		c.unaryInts = ints
	}
}

func WithClientStreamInterceptor(ints ...grpc.StreamClientInterceptor) ClientOption {
	return func(c *clientOptions) {
		c.streamInts = ints
	}
}

// 设置grpc client的拨号选项
func WithClientOptions(rpcOpts []grpc.DialOption) ClientOption {
	return func(c *clientOptions) {
		c.rpcOpts = rpcOpts
	}
}

// 设置负载均衡器
func WithBalancerName(name string) ClientOption {
	return func(c *clientOptions) {
		c.balancerName = name
	}
}

func WithLogger(logger *log.Logger) ClientOption {
	return func(c *clientOptions) {
		c.logger = logger
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// 设置拨号连接
func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin",
		enableTracing: true,
	}
	for _, opt := range opts {
		opt(&options)
	}
	// todo 设置默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
		//opentelemetry
		//otelgrpc.UnaryClientInterceptor(), //配置为可选项
	}
	if options.enableTracing {
		ints = append(ints, otelgrpc.UnaryClientInterceptor())
	}
	streamInts := []grpc.StreamClientInterceptor{}
	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}
	//todo 服务发现的选项
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			discovery.NewBuilder(options.discovery,
				discovery.WithInsecure(insecure),
			),
		))
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
	//return grpc.NewClient(options.endpoint,grpcOpts...)

}

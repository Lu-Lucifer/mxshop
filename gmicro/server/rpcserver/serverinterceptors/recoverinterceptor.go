package serverinterceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop/pkg/log"
	"runtime/debug"
)

// StreamRecoverInterceptor catches panics in processing stream requests and recovers.
func StreamRecoverInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(svr, stream)
}

// UnaryRecoverInterceptor catches panics in processing unary requests and recovers.
func UnaryRecoverInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(ctx, req)
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	log.Errorf("%+v\n\n%s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}

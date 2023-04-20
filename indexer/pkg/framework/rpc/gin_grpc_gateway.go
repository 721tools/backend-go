package rpc

import (
	"context"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcServerRegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type GrpcServerRegisterOption struct {
	f        GrpcServerRegisterFunc
	endpoint string
}

func NewGrpcServerRegisterOption(f GrpcServerRegisterFunc, endpoint string) GrpcServerRegisterOption {
	return GrpcServerRegisterOption{f, endpoint}
}

var grpcDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithUnaryInterceptor(
		grpc_opentracing.UnaryClientInterceptor(
			grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
		),
	),
}

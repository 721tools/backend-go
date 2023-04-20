package rpc

import (
	"context"
	"fmt"

	"github.com/721tools/backend-go/indexer/pkg/framework/third_party"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			stackErr := third_party.Stack(3)
			log.Error(fmt.Sprintf("[Recovery] panic recovered: %s\n%s", err, stackErr))
		}
	}()
	resp, err = handler(ctx, req)
	return resp, err
}

func OpenTracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	rootCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(metautils.ExtractIncoming(ctx)))
	if err == nil {
		parentSpan := opentracing.GlobalTracer().StartSpan(
			info.FullMethod,
			opentracing.ChildOf(rootCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "grpc"})
		defer parentSpan.Finish()
		withCtx := opentracing.ContextWithSpan(ctx, parentSpan)
		return handler(withCtx, req)
	}
	return handler(ctx, req)
}

func RouterLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Info(fmt.Sprintf("gRPC Method: %s, req: %v", info.FullMethod, req))
	resp, err = handler(ctx, req)
	log.Info(fmt.Sprintf("gRPC method: %s, resp: %v", info.FullMethod, resp))
	return resp, err
}

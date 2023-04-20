package middleware

import (
	"context"
	"net/http"

	"github.com/721tools/backend-go/index/pkg/framework/third_party"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

var log = log15.New("module", "middleware.opentracing")

var ginTag = opentracing.Tag{Key: string(ext.Component), Value: "ginTag"}

func GinOpentracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		tracer, closer := third_party.NewTracer()
		defer closer.Close()
		// 直接从 c.Request.Header 中提取 span, 如果没有就新建一个
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(
				c.Request.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "http"})
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				ginTag,
			)
			defer parentSpan.Finish()
		}
		// 然后存到 g.ctx 中 供后续使用
		parentCtx := opentracing.ContextWithSpan(context.Background(), parentSpan)
		c.Set("tracer", tracer)
		c.Set("ctx", parentCtx)

		// Inject header
		err = opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			log.Error("inject to http header failed", "err", err)
		}
		c.Next()
	}
}

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func GrpcGatewayOpentracing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 http header中找到parent span context
		var parentSpan opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))

		if err == nil || err == opentracing.ErrSpanContextNotFound {
			// 找到了 parent span context
			parentSpan = opentracing.StartSpan(
				r.URL.Path,
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				opentracing.ChildOf(spCtx),
				grpcGatewayTag,
			)
			defer parentSpan.Finish()

			// Inject header to req
			err = opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Error("inject to http header failed", "err", err)
			}

			tracer, ok := parentSpan.Context().(jaeger.SpanContext)
			if ok {
				w.Header().Set(jaeger.TraceContextHeaderName, tracer.TraceID().String())
			}
		}
		h.ServeHTTP(w, r)
	})
}

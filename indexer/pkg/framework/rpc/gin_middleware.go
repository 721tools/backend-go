package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/721tools/backend-go/index/pkg/framework/third_party"
	"github.com/721tools/backend-go/index/pkg/utils/log16"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

var log = log16.NewLogger("module", "rpc")

// RecoveryMiddleware 捕获异常
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				stackErr := third_party.Stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")
				if brokenPipe {
					log.Error(fmt.Sprintf("[brokenPipe]%s\n%s", err, headersToStr))
				} else {
					log.Error(fmt.Sprintf("[RecoveryMiddleware] panic recovered: %s\n%s", err, stackErr))
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}

// RouterLoggerMiddleware http router log middleware
func RouterLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		if raw != "" {
			path = path + "?" + raw
		}
		log.Info(fmt.Sprintf("%d %s %s", c.Writer.Status(), c.Request.Method, path),
			"method", c.Request.Method,
			"request_start", start,
			"request_latency", time.Now().Sub(start),
			"request_path", path,
			"status", c.Writer.Status(),
			"err_message", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"body_size", c.Writer.Size(),
			"client_ip", c.ClientIP())
	}
}

// GrpcGatewayOpentracingMiddleware grpc gateway opentracing middleware
func GrpcGatewayOpentracingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opentracing.SpanFromContext(r.Context()) != nil {
			gatewaySpan, ctx := opentracing.StartSpanFromContext(r.Context(), "grpc-gateway")
			defer gatewaySpan.Finish()
			tracer, ok := gatewaySpan.Context().(jaeger.SpanContext)
			if ok {
				r = r.WithContext(ctx)
				w.Header().Set(jaeger.TraceContextHeaderName, tracer.TraceID().String())
			}
		}
		h.ServeHTTP(w, r)
	})
}

// GinOpentracingMiddleware gin开启opentracing中间件
func GinOpentracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		// 直接从 c.Request.Header 中提取 span, 如果没有就新建一个
		rootCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			// err != nil，那么从说明rootCtx是不存在的
			parentSpan = opentracing.GlobalTracer().StartSpan(
				c.Request.URL.Path,
				opentracing.Tag{Key: string(ext.Component), Value: "gin"})
			defer parentSpan.Finish()
		} else {
			// rootCtx存在，生成child span
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(rootCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "gin"},
			)
			defer parentSpan.Finish()
		}
		// 然后存到 g.ctx 中 供后续使用
		parentCtx := opentracing.ContextWithSpan(context.Background(), parentSpan)
		c.Set("ctx", parentCtx)
		c.Request = c.Request.WithContext(parentCtx)

		// Inject header
		err = opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			log.Error("inject to http header failed", "err", err)
		}
		c.Next()
	}
}

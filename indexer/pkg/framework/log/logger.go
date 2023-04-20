package log

import (
	"context"
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"runtime"
	"strings"
)

var l = log15.New("module", "log")

func WithCtx(ctx context.Context) log15.Logger {
	var logger = l
	if _, file, line, isOk := runtime.Caller(1); isOk {
		logger = logger.New("caller", fmt.Sprintf("%s:%d", strings.ReplaceAll(strings.ReplaceAll(file, "\\", "/"), "", ""), line))
	}
	if tracer, ok := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext); ok {
		logger = logger.New(jaeger.TraceContextHeaderName, tracer.TraceID().String())
	}
	return logger
}

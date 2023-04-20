package third_party

import (
	"fmt"

	"github.com/opentracing/opentracing-go"

	//"github.com/721tools/backend-go/index/configs"
	"io"

	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func init() {
	tracerMap = make(map[string]Handler)
	tracerMap["jaeger"] = Jaeger
}

var tracerMap map[string]Handler

type Handler func() (opentracing.Tracer, io.Closer)

func NewTracer() (opentracing.Tracer, io.Closer) {
	if handler, ok := tracerMap[""]; ok {
		return handler()
	}
	return Jaeger()
}

func Jaeger() (opentracing.Tracer, io.Closer) {
	cfg := jaegerConfig.Configuration{
		ServiceName: "",
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "",
		},
	}
	// 不传递 logger 就不会打印日志
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

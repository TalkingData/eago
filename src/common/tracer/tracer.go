package tracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"io"
)

const (
	ctxTracerKey = "Opentracing-Context"
)

// NewTracer 创建一个jaeger Tracer
func NewTracer(srvName, jaegerHostPort string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, //固定采样
			Param: 1,                       //1=全采样、0=不采样
		},

		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: srvName,
	}

	sender, err := jaeger.NewUDPTransport(jaegerHostPort, 0)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegerCfg.Reporter(reporter),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: new tracer: %v\n", err))
	}

	return tracer, closer
}

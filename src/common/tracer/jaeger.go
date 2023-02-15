package tracer

import (
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"time"
)

// NewJaegerTracer 创建一个JaegerTracer
func NewJaegerTracer(options ...Option) (Tracer, error) {
	opts := newOptions(options...)

	cfg := &jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, //固定采样
			Param: opts.SamplerRate,        //1=全采样、0=不采样
		},

		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  opts.JaegerHostPort,
		},

		ServiceName: opts.RegisterKey,
	}

	sender, err := jaeger.NewUDPTransport(opts.JaegerHostPort, 0)
	if err != nil {
		return nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	t, closer, err := cfg.NewTracer(
		jaegerCfg.Reporter(reporter),
	)
	if err != nil {
		return nil, err
	}

	return &tracer{
		tracer: t,
		closer: closer,
	}, nil
}

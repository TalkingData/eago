package tracer

import (
	"github.com/opentracing/opentracing-go"
	"io"
)

type Tracer interface {
	GetTracer() opentracing.Tracer
	Close()
}

type tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func (t *tracer) GetTracer() opentracing.Tracer {
	return t.tracer
}

func (t *tracer) Close() {
	if t.closer != nil {
		_ = t.closer.Close()
	}
}

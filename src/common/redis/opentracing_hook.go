package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
)

type OpentracingHook struct {
	globalTracer opentracing.Tracer
}

func NewOpentracingHook() *OpentracingHook {
	return &OpentracingHook{
		globalTracer: opentracing.GlobalTracer(),
	}
}

func (oh *OpentracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}

	span, ctxWithSpan := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("redis.%s", cmd.FullName()))

	ext.DBType.Set(span, "redis")
	ext.DBStatement.Set(span, cmd.String())

	return ctxWithSpan, nil
}

func (oh *OpentracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	defer span.Finish()

	if err := cmd.Err(); err != nil {
		ext.Error.Set(span, true)

		span.LogFields(tracerLog.Error(err))
	}

	return nil
}

func (oh *OpentracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if ctx == nil {
		return ctx, nil
	}

	span, ctxWithSpan := opentracing.StartSpanFromContext(ctx, "redis.pipeline")

	ext.DBType.Set(span, "redis")
	span.SetTag("db.redis.num_cmd", len(cmds))
	for idx, cmd := range cmds {
		span.SetTag(fmt.Sprintf("db.statement[%d]", idx), cmd.String())
	}

	return ctxWithSpan, nil

}
func (oh *OpentracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	defer span.Finish()

	for idx, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			ext.Error.Set(span, true)
			span.LogFields(tracerLog.Error(fmt.Errorf("db.statement[%d]: %w", idx, err)))
		}
	}

	return nil
}

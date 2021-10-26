package tracer

import (
	"context"
	"eago/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// ExtractTraceContext 提取Context
func ExtractTraceContext(c *gin.Context) context.Context {
	v, exist := c.Get(ctxTracerKey)
	if exist == false {
		return context.Background()
	}

	ctx, ok := v.(context.Context)
	if ok {
		return ctx
	}
	return context.Background()
}

// StartSpanFromContext 从Context创建新span
func StartSpanFromContext(ctx context.Context) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, utils.GetFuncName(3))
}

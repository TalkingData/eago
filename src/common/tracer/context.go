package tracer

import (
	"context"
	"eago/common/global"
	"github.com/gin-gonic/gin"
)

// ExtractTraceCtxFromGin 从gin.context中提取带有tracer的Context
func ExtractTraceCtxFromGin(c *gin.Context) context.Context {
	if v, exist := c.Get(global.OpentracingCtxKey); exist {
		if ctx, ok := v.(context.Context); ok {
			return ctx
		}
	}

	return context.Background()
}

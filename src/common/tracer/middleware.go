package tracer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

// TracerMiddleware tracer中间件
func TracerMiddleware(c *gin.Context) {
	//if we need to log res body
	gbw := ginBodyWriter{buffer: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = gbw

	tracer := opentracing.GlobalTracer()
	// 创建新span
	span := tracer.StartSpan(c.Request.URL.Path)
	md := make(map[string]string)
	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	if err == nil {
		span = tracer.StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		tracer = span.Tracer()
	}
	defer span.Finish()

	err = tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		fmt.Println(err)
	}

	// 设置context写入gin.Context中
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)
	c.Set(ctxTracerKey, ctx)

	c.Next()

	// 取得http返回码
	statusCode := c.Writer.Status()
	// 将http返回码写入span
	ext.HTTPStatusCode.Set(span, uint16(statusCode))
	// 将http request method写入span
	ext.HTTPMethod.Set(span, c.Request.Method)
	// 将http request url写入span
	ext.HTTPUrl.Set(span, c.Request.URL.EscapedPath())
	//  将http status大于等于400的返回码的span均标记为错误
	if statusCode >= http.StatusBadRequest {
		ext.Error.Set(span, true)
	}

	// 劫持Gin response查看code，如果code!=0说明有错误发生
	resp := gbw.GetMapString()
	code, ok := resp["code"]
	if ok && code.(float64) > 0 {
		ext.Error.Set(span, true)
		span.LogKV("response.code", code)
		m, ok := resp["message"]
		if ok {
			span.LogKV("response.message", m)
		}
	}
}

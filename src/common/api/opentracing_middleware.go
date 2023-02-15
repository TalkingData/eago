package api

import (
	"bytes"
	"context"
	perm "eago/common/api/permission"
	"eago/common/global"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

// OpentracingMiddleware opentracing的gin中间件
func OpentracingMiddleware(c *gin.Context) {
	// if we need to log res body
	gbw := &ginBodyWriter{buffer: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = gbw

	// 创建新span
	operationName := fmt.Sprintf("%s_%s", c.Request.Method, c.Request.URL.Path)
	span, ctx := opentracing.StartSpanFromContext(context.Background(), operationName)
	defer span.Finish()

	// 将tracer id填入response header
	if jSpanCtx, ok := span.Context().(jaeger.SpanContext); ok {
		tId := jSpanCtx.TraceID().String()
		c.Header(global.GinRespHeaderTracerIdKey, tId)
		c.Set(global.GinCtxTracerIdKey, tId)
	}

	md := map[string]string{}
	_ = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))

	// 设置context写入gin.Context中
	c.Set(global.OpentracingCtxKey, metadata.NewContext(opentracing.ContextWithSpan(ctx, span), md))

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

	if tc, ok := perm.GetTokenContent(c); ok {
		span.SetTag("token_content.username", tc.Username)
		span.SetTag("token_content.user_id", tc.UserId)
	}

	// 劫持Gin response查看code，如果code!=0说明有错误发生
	resp := gbw.GetMapString()
	if code, ok := resp["code"]; ok {
		if code.(float64) > 0 {
			ext.Error.Set(span, true)
		}
		span.SetTag("response.code", code)
	}
	if msg, ok := resp["message"]; ok {
		span.SetTag("response.message", msg)
	}
}

// ginBodyWriter gin返回Writer，暂时放在这里
type ginBodyWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

// Write
func (w ginBodyWriter) Write(b []byte) (int, error) {
	// memory copy here!
	w.buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ginBodyWriter) GetBodyString() string {
	return w.buffer.String()
}

func (w *ginBodyWriter) GetMapString() map[string]interface{} {
	ret := make(map[string]interface{})
	_ = json.Unmarshal([]byte(w.GetBodyString()), &ret)
	return ret
}

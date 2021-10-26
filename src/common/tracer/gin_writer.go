package tracer

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// ginBodyWriter gin返回Writer，暂时放在这里
type ginBodyWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

// Write
func (w ginBodyWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

// GetBodyString
func (w *ginBodyWriter) GetBodyString() string {
	return w.buffer.String()
}

func (w *ginBodyWriter) GetMapString() map[string]interface{} {
	ret := make(map[string]interface{})
	_ = json.Unmarshal([]byte(w.GetBodyString()), &ret)
	return ret
}

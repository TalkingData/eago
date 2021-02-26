package middleware

import (
	"eago-common/api-suite/handler"
	"github.com/gin-gonic/gin"
	"strings"
)

func Dispatch(realPatch string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, realPatch) {
			handler.PageNoFound(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

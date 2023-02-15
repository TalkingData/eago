package api

import (
	"eago/common/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func GinCustomLogger(lg *logger.Logger) gin.HandlerFunc {
	var skip map[string]struct{}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			ts := time.Now()

			if raw != "" {
				path = path + "?" + raw
			}
			lg.DebugWithFields(logger.Fields{
				"status_code":   c.Writer.Status(),
				"latency":       ts.Sub(start),
				"client_ip":     c.ClientIP(),
				"method":        c.Request.Method,
				"path":          path,
				"body_size":     c.Writer.Size(),
				"error_message": c.Errors.ByType(gin.ErrorTypePrivate).String(),
			}, "Got gin request.")
		}
	}
}

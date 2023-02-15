package api

import (
	"eago/common/api/ext"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PageNotFound 找不到页面时使用的Handler
func PageNotFound(c *gin.Context) {
	ext.WriteAny(c, http.StatusNotFound, "Page not found.")
}

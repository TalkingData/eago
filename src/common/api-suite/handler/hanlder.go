package handler

import (
	w "eago/common/api-suite/writter"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PageNoFound 找不到页面时使用的Handler
func PageNoFound(c *gin.Context) {
	w.WriteAny(c, http.StatusNotFound, "Page not found.")
}

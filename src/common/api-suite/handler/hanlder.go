package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// PageNoFound 找不到页面时使用的Handler
func PageNoFound(c *gin.Context) {
	resp := gin.H{
		"code":    http.StatusNotFound,
		"message": "Page not found.",
	}
	c.JSON(http.StatusOK, resp)
}

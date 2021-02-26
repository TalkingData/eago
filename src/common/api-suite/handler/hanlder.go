package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageNoFound(c *gin.Context) {
	resp := gin.H{
		"code":    http.StatusNotFound,
		"message": "Page not found.",
	}
	c.JSON(http.StatusOK, resp)
}

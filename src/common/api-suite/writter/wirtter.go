package writter

import (
	"eago/common/api-suite/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

// WriteSuccessPayload 写带有数据的成功消息到RestfulApi
func WriteSuccessPayload(c *gin.Context, key string, pld interface{}) {
	rsp := make(gin.H)
	rsp["code"] = 0
	rsp["message"] = "Success"
	defer c.JSON(http.StatusOK, rsp)

	if pld == nil {
		return
	}

	switch pldType := pld.(type) {
	case *pagination.Paginator:
		rsp["page"] = pldType.Page
		rsp["pages"] = pldType.Pages
		rsp["page_size"] = pldType.PageSize
		rsp["total"] = pldType.Total
		rsp[key] = pldType.Data

	default:
		rsp[key] = pld
	}
}

// WriteSuccess 写成功消息到RestfulApi
func WriteSuccess(c *gin.Context) {
	rsp := make(gin.H)
	rsp["code"] = 0
	rsp["message"] = "Success"

	c.JSON(http.StatusOK, rsp)
}

// WriteAny 自由写消息RestfulApi
func WriteAny(c *gin.Context, code int, message string) {
	rsp := make(gin.H)
	rsp["code"] = code
	rsp["message"] = message

	c.JSON(http.StatusOK, rsp)
}

// WriteAnyAndAbort 自由写消息RestfulApi
func WriteAnyAndAbort(c *gin.Context, code int, message string) {
	rsp := make(gin.H)
	rsp["code"] = code
	rsp["message"] = message

	c.AbortWithStatusJSON(http.StatusOK, rsp)
}

// WriteAnyAndAbortWithError 自由写消息RestfulApi
func WriteAnyAndAbortWithError(c *gin.Context, code int, message string, err interface{}) {
	rsp := make(gin.H)
	rsp["code"] = code
	rsp["message"] = message
	rsp["error"] = err

	c.AbortWithStatusJSON(http.StatusOK, rsp)
}

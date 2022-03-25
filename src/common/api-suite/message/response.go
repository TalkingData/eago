package message

import (
	"eago/common/api-suite/pagination"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	code    int
	message string
	payload gin.H
}

// Deprecated: SetPayload 装载消息的Payload
func (resp *response) SetPayload(key string, v interface{}) *response {
	resp.payload[key] = v
	return resp
}

// Deprecated: SetPagedPayload 装载分页过的Payload
func (resp *response) SetPagedPayload(paged *pagination.Paginator, pldKey string) *response {
	resp.payload["page"] = paged.Page
	resp.payload["pages"] = paged.Pages
	resp.payload["page_size"] = paged.PageSize
	resp.payload["total"] = paged.Total
	resp.payload[pldKey] = paged.Data

	return resp
}

// Deprecated: Write
func (resp *response) Write(to *gin.Context) {
	resp.payload["code"] = resp.code
	resp.payload["message"] = resp.message
	to.JSON(http.StatusOK, resp.payload)
}

// Deprecated: WriteAndAbort
func (resp *response) WriteAndAbort(to *gin.Context) {
	resp.payload["code"] = resp.code
	resp.payload["message"] = resp.message
	to.AbortWithStatusJSON(http.StatusOK, resp.payload)
}

// Deprecated: String 返回消息的字符串
func (resp *response) String() string {
	return resp.message + "."
}

// Deprecated: Error 以Error实例返回
func (resp *response) Error() error {
	return errors.New(resp.message)
}

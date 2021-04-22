package message

import (
	"eago-common/api-suite/pagination"
	"github.com/gin-gonic/gin"
)

type Msg struct {
	BaseMsg
	payload *gin.H
}

type BaseMsg struct {
	Code    int
	Message string
}

// NewMsg 新生成一条消息
func (bm *BaseMsg) NewMsg(detail ...string) *Msg {
	var m = Msg{}
	m.Code = bm.Code
	m.Message = bm.Message

	for _, d := range detail {
		m.Message = m.Message + " " + d
	}

	return &m
}

// SetPayload 装载消息的Payload
func (m *Msg) SetPayload(payload *gin.H) *Msg {
	m.payload = payload
	return m
}

// SetPagedPayload 装载分页过的Payload
func (m *Msg) SetPagedPayload(paged *pagination.Paginator, objRename string) *Msg {
	m.payload = &gin.H{
		"page":      paged.Page,
		"pages":     paged.Pages,
		"page_size": paged.PageSize,
		"total":     paged.Total,
	}

	pld := *m.payload
	pld[objRename] = paged.Data
	return m
}

// GinH 返回GinH
func (m *Msg) GinH() *gin.H {
	var h = gin.H{}
	if m.payload != nil {
		for k, v := range *m.payload {
			h[k] = v
		}
	}

	h["code"] = m.Code
	h["message"] = m.Message
	return &h
}

// String 返回消息的字符串
func (m *Msg) String() string {
	return m.Message
}

package message

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Message struct {
	Code    int
	Message string
}

// SetDetail 设置消息详情
func (m *Message) SetDetail(detail ...string) *Message {
	var newM = Message{}
	if len(detail) >= 1 {
		newM.Message = m.Message + " " + strings.Join(detail, " ")
	} else {
		newM.Message = m.Message
	}

	return &newM
}

// String 返回消息的字符串
func (m *Message) String() string {
	return m.Message
}

// GenResponse 新生成一条响应体
func (m *Message) GenResponse(detail ...string) *response {
	var resp = response{}
	resp.code = m.Code
	if len(detail) >= 1 {
		resp.message = m.Message + " " + strings.Join(detail, " ")
	} else {
		resp.message = m.Message
	}
	resp.payload = make(gin.H)

	return &resp
}

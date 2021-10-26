package message

import (
	"eago/common/log"
	"eago/common/utils"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const RPC_ERROR_SPLITOR = ", "

type Message struct {
	code int
	msg  string
	err  interface{}
}

func NewMessage(code int, msg string) *Message {
	return &Message{
		code: code,
		msg:  msg,
	}
}

// SetDetail 设置消息详情
func (m *Message) SetDetail(detail ...string) *Message {
	newM := *m
	if len(detail) > 0 {
		newM.msg = newM.msg + " " + strings.Join(detail, " ")
	}

	return &newM
}

// SetError 设置成功消息的Error
func (m *Message) SetError(err interface{}, detail ...string) *Message {
	newM := *m
	newM.err = err
	if len(detail) > 0 {
		newM.msg = newM.msg + " " + strings.Join(detail, " ")
	}

	return &newM
}

// WriteRest 写成息到RestfulApi
func (m *Message) WriteRest(c *gin.Context) {
	rsp := make(gin.H)
	rsp["code"] = m.code
	rsp["message"] = m.msg
	defer c.JSON(http.StatusOK, rsp)

	if m.err == nil {
		return
	}

	rsp["errors"] = formatError(m.err)

}

// LogFields 转换为日志字段类型
func (m *Message) LogFields() (log.Fields, interface{}) {
	f := make(log.Fields)
	f["code"] = m.code
	if m.err != nil {
		f["errors"] = formatError(m.err)
	}

	return f, m.String()
}

// RpcError 转换为Rpc error类型
func (m *Message) RpcError() error {
	return fmt.Errorf("%d%s%s", m.code, RPC_ERROR_SPLITOR, m.msg)
}

// String 转换为string类型
func (m *Message) String() string {
	return m.msg
}

// Code 返回Message的code值
func (m *Message) Code() int {
	return m.code
}

func formatError(err interface{}) interface{} {
	switch errType := err.(type) {
	case error:
		return errType.Error()

	case []*validation.Error:
		errMap := make(map[string]string)
		for _, e := range errType {
			errMap[utils.Camel2Case(e.Field)] = e.Message
		}
		return errMap
	default:
		return err
	}
}

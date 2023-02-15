package code_msg

import (
	"eago/common/logger"
	"eago/common/utils"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

type CodeMsg struct {
	code int
	msg  string
	err  interface{}
}

func NewCodeMsg(code int, msg string) *CodeMsg {
	return &CodeMsg{
		code: code,
		msg:  msg,
	}
}

// ToMicroErr 将CodeMsg类型转为ToMicroErr
func (cm *CodeMsg) ToMicroErr() error {
	return errors.New(defaultMicroErrId, cm.msg, int32(cm.code))
}

// ToErr 转换为error类型
func (cm *CodeMsg) ToErr() error {
	return fmt.Errorf("%d%s%s", cm.code, defaultSeparator, cm.msg)
}

// ToRpcErr 将CodeMsg类型转为RpcError
func (cm *CodeMsg) ToRpcErr() error {
	return status.Error(codes.Code(cm.code), cm.msg)
}

// Write2GinCtx 将CodeMsg写到Gin.Context中
func (cm *CodeMsg) Write2GinCtx(c *gin.Context) {
	ginRsp := make(gin.H)
	ginRsp["code"] = cm.code
	ginRsp["message"] = cm.msg
	defer c.JSON(http.StatusOK, ginRsp)

	if cm.err != nil {
		ginRsp["errors"] = fmtErr(cm.err)
	}
}

// ToLoggerFields 转换为日志字段类型
func (cm *CodeMsg) ToLoggerFields() logger.Fields {
	f := make(logger.Fields)
	f["code"] = cm.code
	if cm.err != nil {
		f["errors"] = fmtErr(cm.err)
	}

	return f
}

// GetMsg 获取Msg
func (cm *CodeMsg) GetMsg() string {
	return cm.msg
}

// GetCode 获取Code
func (cm *CodeMsg) GetCode() int {
	return cm.code
}

// SetDetail 设置消息详情
func (cm *CodeMsg) SetDetail(details ...string) *CodeMsg {
	if len(details) < 1 {
		return cm
	}

	newCm := *cm
	newCm.msg += strings.Join(details, defaultSeparator)
	return &newCm
}

// SetError 设置成功消息的Error
func (cm *CodeMsg) SetError(err interface{}, details ...string) *CodeMsg {
	newCm := *cm
	newCm.err = err
	if len(details) > 0 {
		newCm.msg += strings.Join(details, defaultSeparator)
	}

	return &newCm
}

func fmtErr(err interface{}) interface{} {
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

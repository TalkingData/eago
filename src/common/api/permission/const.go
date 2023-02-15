package permission

import (
	cMsg "eago/common/code_msg"
	"net/http"
)

const (
	defaultTokenKey              = "token"
	defaultTokenContentGinCtxKey = "__eago_token_content"
)

var (
	msgAuthClientErr = cMsg.NewCodeMsg(http.StatusInternalServerError, "访问Auth的Srv服务时发生意外，请先尝试重试，若无效请联系管理员")
)

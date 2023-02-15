package code_msg

import (
	"net/http"
)

var (
	MsgInvalidUriFailed   = NewCodeMsg(http.StatusBadRequest, "请求错误，请确保Params中包含以下字段:")
	MsgSerializeFailed    = NewCodeMsg(http.StatusBadRequest, "无法读取提交的数据，请确保发送正确的数据格式")
	MsgValidateFailed     = NewCodeMsg(http.StatusBadRequest, "验证失败，提交的内容不合法")
	MsgInvalidTokenFailed = NewCodeMsg(http.StatusUnauthorized, "无效的Token")
	MsgUserNotRoleFailed  = NewCodeMsg(http.StatusUnauthorized, "无权限，需要当前用户属于以下角色中的任意一个:")
	MsgNotFoundFailed     = NewCodeMsg(http.StatusNotFound, "操作失败，找不到指定对象")

	MsgCheckRoleErr = NewCodeMsg(http.StatusInternalServerError, "检查用户角色时出错，请联系管理员")
	MsgUndefinedErr = NewCodeMsg(http.StatusInternalServerError, "遇到未知错误失败，请先尝试重试，若无效请联系管理员")
)

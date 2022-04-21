package permission

import (
	m "eago/common/message"
	"net/http"
)

var (
	InvalidParams = m.NewMessage(http.StatusBadRequest, "请求错误，请确保Params中包含以下字段:")

	InvalidToken = m.NewMessage(http.StatusForbidden, "Token不正确，请确认已经登录")
	UserNotRole  = m.NewMessage(http.StatusForbidden, "无权限，需要当前用户属于以下角色中的任意一个:")

	CheckRoleError = m.NewMessage(http.StatusInternalServerError, "检查用户角色时出错，请联系管理员")
)

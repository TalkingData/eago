package msg

import (
	m "eago/common/message"
	"net/http"
)

var (
	// Default
	InvalidUriFailed = m.NewMessage(http.StatusBadRequest, "请求错误，请确保Params中包含以下字段:")
	SerializeFailed  = m.NewMessage(http.StatusBadRequest, "无法读取提交的数据，请确保发送正确的数据格式")
	ValidateFailed   = m.NewMessage(http.StatusBadRequest, "验证失败，提交的内容不合法")
	NotFoundFailed   = m.NewMessage(http.StatusNotFound, "操作失败，找不到指定对象")
	UndefinedError   = m.NewMessage(http.StatusInternalServerError, "操作失败，请先尝试重试，若无效请联系管理员")

	// Login 1000xx
	LoginInactiveCrowdUserFailed = m.NewMessage(100001, "登录失败，当前用户在Crowd中是禁用状态")
	LoginDisabledUserFailed      = m.NewMessage(100002, "登录失败，用户处于禁用状态")
	LoginNoPasswordUserFailed    = m.NewMessage(100003, "登录失败，当前用户没有设置密码")
	LoginAuthenticationFailed    = m.NewMessage(100004, "登录失败，请检查用户名密码是否正确")
	LoginNewTokenFailed          = m.NewMessage(100005, "创建Token失败，请联系管理员")
	LoginUnknownFailed           = m.NewMessage(100099, "登录失败，请联系管理员")

	// User 1001xx
	HandoverUnknownError = m.NewMessage(100101, "用户交接处理失败，请联系管理员")

	// Role 1002xx
	AssociatedRoleFailed = m.NewMessage(100201, "无法执行操作，仍有用户与该角色关联")

	// Role 1003xx
	AssociatedProductFailed = m.NewMessage(100301, "无法执行操作，仍有用户与该产品线关联")

	// Role 1004xx
	AssociatedGroupFailed = m.NewMessage(100401, "无法执行操作，仍有用户与该组关联")

	// Role 1005xx
	AssociatedDepartmentFailed     = m.NewMessage(100501, "无法执行操作，仍有子部门与该部门关联")
	AssociatedDepartmentUserFailed = m.NewMessage(100502, "无法执行操作，仍有用户与该部门关联")
)

package msg

import (
	cMsg "eago/common/code_msg"
)

var (
	// Login 1000xx
	MsgLoginInactiveCrowdUserFailed = cMsg.NewCodeMsg(100001, "登录失败，当前用户在Crowd中是禁用状态")
	MsgLoginDisabledUserFailed      = cMsg.NewCodeMsg(100002, "登录失败，用户处于禁用状态")
	MsgLoginNoPasswordUserFailed    = cMsg.NewCodeMsg(100003, "登录失败，当前用户没有设置密码")
	MsgLoginAuthenticationFailed    = cMsg.NewCodeMsg(100004, "登录失败，请检查用户名密码是否正确")
	MsgLoginNewTokenFailed          = cMsg.NewCodeMsg(100005, "创建Token失败，请联系管理员")
	MsgGetTokenContentFailed        = cMsg.NewCodeMsg(100006, "获取TokenContent失败")
	MsgLoginUnknownFailed           = cMsg.NewCodeMsg(100099, "登录失败，请联系管理员")

	// Department 1001xx
	MsgAssociatedDepartmentFailed     = cMsg.NewCodeMsg(100115, "无法执行操作，仍有子部门与该部门关联")
	MsgAssociatedDepartmentUserFailed = cMsg.NewCodeMsg(100120, "无法执行操作，仍有用户与该部门关联")

	// Group 1002xx
	MsgAssociatedGroupFailed = cMsg.NewCodeMsg(100220, "无法执行操作，仍有用户与该组关联")

	// Product 1003xx
	MsgAssociatedProductFailed = cMsg.NewCodeMsg(100320, "无法执行操作，仍有用户与该产品线关联")

	// Role 1004xx
	MsgAssociatedRoleFailed = cMsg.NewCodeMsg(100420, "无法执行操作，仍有用户与该角色关联")

	// User 1005xx
	MsgUserHandoverFailed = cMsg.NewCodeMsg(100500, "用户交接操作失败")

	// Others
	MsgAuthDaoErr   = cMsg.NewCodeMsg(109900, "Auth服务的DAO层发生意外，请先尝试重试，若无效请联系管理员")
	MsgAuthCacheErr = cMsg.NewCodeMsg(109901, "Auth服务的Cache发生意外，请先尝试重试，若无效请联系管理员")
)

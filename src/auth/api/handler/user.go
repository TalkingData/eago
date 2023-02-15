package handler

import (
	"eago/auth/api/form"
	"eago/auth/conf/msg"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

// SetUser 更新用户
func (ah *AuthHandler) SetUser(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err)
		logF := m.ToLoggerFields()
		logF["user_id"] = userId
		ah.logger.WarnWithFields(logF, m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := new(form.SetUserForm)
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	u, err := ah.dao.SetUser(ctx, userId, frm.Email, frm.Phone)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "user", u)
}

// PagedListUsers 列出所有用户-分页
func (ah *AuthHandler) PagedListUsers(c *gin.Context) {
	// 设置查询filter
	pFrm := new(form.PagedListUsersParamsForm)
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.PagedListUsers, skipped it.")
	}

	paged, err := ah.dao.PagedListUsers(
		tracer.ExtractTraceCtxFromGin(c),
		pFrm.GenQuery(),
		c.GetInt(global.GinCtxPageKey),
		c.GetInt(global.GinCtxPageSizeKey),
		c.GetStringSlice(global.GinCtxOrderByKey)...,
	)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "users", paged)
}

// ListUsersRoles 列出用户所属角色
func (ah *AuthHandler) ListUsersRoles(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	roles, err := ah.dao.ListUsersRoles(tracer.ExtractTraceCtxFromGin(c), userId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "roles", roles)
}

// ListUsersProducts 列出用户所在产品线
func (ah *AuthHandler) ListUsersProducts(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	prods, err := ah.dao.ListUsersProducts(tracer.ExtractTraceCtxFromGin(c), userId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "products", prods)
}

// ListUsersGroups 列出用户所在组
func (ah *AuthHandler) ListUsersGroups(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	groups, err := ah.dao.ListUsersGroups(tracer.ExtractTraceCtxFromGin(c), userId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "groups", groups)
}

// GetUsersDepartment 获得指定用户所在部门
func (ah *AuthHandler) GetUsersDepartment(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	dept, err := ah.dao.GetUsersDepartment(tracer.ExtractTraceCtxFromGin(c), userId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if dept != nil {
		ext.WriteSuccessPayload(c, "department", dept)
		return
	}
	ext.WriteSuccessPayload(c, "department", "{}")
}

// GetUsersDepartmentChain 获得指定用户所在部门链，包含所有层级情况
func (ah *AuthHandler) GetUsersDepartmentChain(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(
		c, "department", ah.dao.GetUsersDepartmentChain(tracer.ExtractTraceCtxFromGin(c), userId),
	)
}

// MakeUserHandover 用户交接
func (ah *AuthHandler) MakeUserHandover(c *gin.Context) {
	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	tgtUserId, err := ext.ParamUint32(c, "target_user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "target_user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := new(form.MakeUserHandoverForm)
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, userId, tgtUserId); m != nil {
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 操作用户交接
	if err = ah.biz.MakeUserHandover(ctx, userId, tgtUserId); err != nil {
		m := msg.MsgUserHandoverFailed.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

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

// NewRole 新建角
func (ah *AuthHandler) NewRole(c *gin.Context) {
	frm := form.NewRoleForm{}
	// 序列化request body
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	r, err := ah.dao.NewRole(ctx, frm.Name, *frm.Description)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "role", r)
}

// RemoveRole 删除角色
func (ah *AuthHandler) RemoveRole(c *gin.Context) {
	roleId, err := ext.ParamUint32(c, "role_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	frm := form.RemoveRoleForm{}
	if m := frm.Validate(ctx, ah.dao, roleId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 删除
	if err = ah.dao.RemoveRole(ctx, roleId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetRole 更新角色
func (ah *AuthHandler) SetRole(c *gin.Context) {
	roleId, err := ext.ParamUint32(c, "role_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetRoleForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, roleId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	role, err := ah.dao.SetRole(ctx, roleId, frm.Name, *frm.Description)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "role", role)
}

// PagedListRoles 列出所有角色-分页
func (ah *AuthHandler) PagedListRoles(c *gin.Context) {
	// 设置查询filter
	pFrm := new(form.PagedListRolesParamsForm)
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.PagedListRoles, skipped it.")
	}

	paged, err := ah.dao.PagedListRoles(
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

	ext.WriteSuccessPayload(c, "roles", paged)
}

// AddUser2Role 添加用户至角色
func (ah *AuthHandler) AddUser2Role(c *gin.Context) {
	roleId, err := ext.ParamUint32(c, "role_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.AddUser2RoleForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, roleId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.AddUser2Role(ctx, roleId, frm.UserId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// RemoveRolesUser 移除角色中用户
func (ah *AuthHandler) RemoveRolesUser(c *gin.Context) {
	roleId, err := ext.ParamUint32(c, "role_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "user_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveRolesUserForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, roleId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveRolesUser(ctx, roleId, userId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// ListRolesUsers 列出角色所有用户
func (ah *AuthHandler) ListRolesUsers(c *gin.Context) {
	roleId, err := ext.ParamUint32(c, "role_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	u, err := ah.dao.ListRolesUsers(tracer.ExtractTraceCtxFromGin(c), roleId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "users", u)
}

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

// NewGroup 新建组
func (ah *AuthHandler) NewGroup(c *gin.Context) {
	frm := form.NewGroupForm{}
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

	group, err := ah.dao.NewGroup(ctx, frm.Name, *frm.Description)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "group", group)
}

// RemoveGroup 删除组
func (ah *AuthHandler) RemoveGroup(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "group_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveGroupForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, gId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveGroup(ctx, gId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetGroup 更新组
func (ah *AuthHandler) SetGroup(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "group_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetGroupForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, gId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	group, err := ah.dao.SetGroup(ctx, gId, frm.Name, *frm.Description)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "group", group)
}

// PagedListGroups 列出所有组-分页
func (ah *AuthHandler) PagedListGroups(c *gin.Context) {
	// 设置查询filter
	pFrm := form.PagedListGroupsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.PagedListGroups, skipped it.")
	}

	paged, err := ah.dao.PagedListGroups(
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

	ext.WriteSuccessPayload(c, "groups", paged)
}

// AddUser2Group 添加用户至组
func (ah *AuthHandler) AddUser2Group(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "group_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.AddUser2GroupForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, gId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.AddUser2Group(ctx, gId, frm.UserId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// RemoveGroupsUser 移除组中用户
func (ah *AuthHandler) RemoveGroupsUser(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
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

	frm := form.RemoveGroupsUserForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, gId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveGroupsUser(ctx, gId, userId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetGroupsOwner 设置用户是否是组Owner
func (ah *AuthHandler) SetGroupsOwner(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
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

	frm := form.SetGroupsOwnerForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, gId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.SetGroupsOwner(ctx, gId, userId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// ListGroupsUsers 列出角色所有用户
func (ah *AuthHandler) ListGroupsUsers(c *gin.Context) {
	gId, err := ext.ParamUint32(c, "group_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	pFrm := form.ListGroupsUsersParamsForm{}
	// 序列化request body
	if err = c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.ListGroupsUsers, skipped it.")
	}

	if m := pFrm.Validate(ctx, ah.dao, gId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	u, err := ah.dao.ListGroupsUsers(ctx, gId, pFrm.GenQuery())
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "users", u)
}

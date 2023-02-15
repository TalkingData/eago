package handler

import (
	"eago/auth/api/form"
	"eago/auth/conf/msg"
	"eago/auth/dto"
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

// NewDepartment 新建部门
func (ah *AuthHandler) NewDepartment(c *gin.Context) {
	frm := form.NewDepartmentForm{}
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

	// 新建
	dept, err := ah.dao.NewDepartment(ctx, frm.Name, frm.ParentId)
	if dept == nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "department", dept)
}

// RemoveDepartment 删除部门
func (ah *AuthHandler) RemoveDepartment(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveDepartmentForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, deptId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveDepartment(ctx, deptId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetDepartment 更新部门
func (ah *AuthHandler) SetDepartment(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetDepartmentForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, deptId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	dept, err := ah.dao.SetDepartment(ctx, deptId, frm.Name, frm.ParentId)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "department", dept)
}

// PagedListDepartments 列出所有部门-分页
func (ah *AuthHandler) PagedListDepartments(c *gin.Context) {
	// 设置查询filter
	pFrm := form.PagedListDepartmentsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.PagedListDepartments, skipped it.")
	}

	paged, err := ah.dao.PagedListDepartments(
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

	ext.WriteSuccessPayload(c, "departments", paged)
}

// ListDepartmentFullTree 以树结构列出所有部门
func (ah *AuthHandler) ListDepartmentFullTree(c *gin.Context) {
	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 查找根部门
	dept, err := ah.dao.GetDepartment(ctx, orm.Query{"parent_id": nil})
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		ext.WriteSuccessPayload(c, "tree", make(map[string]interface{}))
		return
	}

	// 列出所有部门
	deptList, err := ah.dao.ListDepartments(ctx, orm.Query{})
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 将根部门转化为树结构
	root := dto.TransDepartment2Tree(dept)
	ah.dao.ListDepartmentTree(root, deptList)

	ext.WriteSuccessPayload(c, "tree", root)
}

// ListDepartmentTree 列出指定部门子树
func (ah *AuthHandler) ListDepartmentTree(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 查找指定部门
	dept, err := ah.dao.GetDepartment(ctx, orm.Query{"id=?": deptId})
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		ext.WriteSuccessPayload(c, "tree", make(map[string]interface{}))
		return
	}

	// 列出所有部门
	deptList, err := ah.dao.ListDepartments(ctx, orm.Query{})
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 将根部门转化为树结构
	root := dto.TransDepartment2Tree(dept)
	ah.dao.ListDepartmentTree(root, deptList)

	ext.WriteSuccessPayload(c, "tree", root)
}

// AddUser2Department 添加用户至部门
func (ah *AuthHandler) AddUser2Department(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.AddUser2DepartmentForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, deptId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.AddUser2Department(ctx, frm.UserId, deptId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// RemoveDepartmentsUser 移除部门中用户
func (ah *AuthHandler) RemoveDepartmentsUser(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
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

	frm := form.RemoveDepartmentsUserForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, deptId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveDepartmentsUser(ctx, userId, deptId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetDepartmentsOwner 设置用户是否是部门Owner
func (ah *AuthHandler) SetDepartmentsOwner(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
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

	frm := form.SetDepartmentsOwnerForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, deptId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.SetDepartmentsOwner(ctx, deptId, userId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// ListDepartmentsUsers 列出部门中所有用户
func (ah *AuthHandler) ListDepartmentsUsers(c *gin.Context) {
	deptId, err := ext.ParamUint32(c, "department_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "department_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	pFrm := form.ListDepartmentsUsersParamsForm{}
	// 序列化request body
	if err = c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.ListDepartmentsUsers, skipped it.")
	}

	if m := pFrm.Validate(ctx, ah.dao, deptId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	u, err := ah.dao.ListDepartmentsUsers(ctx, deptId, pFrm.GenQuery())
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "users", u)
}

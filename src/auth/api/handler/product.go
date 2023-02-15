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

// NewProduct 新建产品线
func (ah *AuthHandler) NewProduct(c *gin.Context) {
	frm := form.NewProductForm{}
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

	p, err := ah.dao.NewProduct(ctx, frm.Name, frm.Alias, *frm.Description, frm.Disabled)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "product", p)
}

// RemoveProduct 删除产品线
func (ah *AuthHandler) RemoveProduct(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveProductForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, prdId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveProduct(ctx, prdId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetProduct 更新产品线
func (ah *AuthHandler) SetProduct(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetProductForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, prdId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	prod, err := ah.dao.SetProduct(ctx, prdId, frm.Name, frm.Alias, *frm.Description, *frm.Disabled)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "product", prod)
}

// PagedListProducts 列出所有产品线-分页
func (ah *AuthHandler) PagedListProducts(c *gin.Context) {
	// 设置查询filter
	pFrm := form.PagedListProductsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.PagedListProducts, skipped it.")
	}

	paged, err := ah.dao.PagedListProducts(
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

	ext.WriteSuccessPayload(c, "products", paged)
}

// AddUser2Product 添加用户至产品线
func (ah *AuthHandler) AddUser2Product(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.AddUser2ProductForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, prdId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.AddUser2Product(ctx, prdId, frm.UserId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// RemoveProductsUser 移除产品线中用户
func (ah *AuthHandler) RemoveProductsUser(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveProductsUserForm{}
	// 验证数据
	if m := frm.Validate(ctx, ah.dao, prdId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.RemoveProductsUser(ctx, prdId, userId); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetProductsOwner 设置产品线Owner
func (ah *AuthHandler) SetProductsOwner(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	userId, err := ext.ParamUint32(c, "user_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "role_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetProductsOwnerForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, ah.dao, prdId, userId); m != nil {
		// 数据验证未通过
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = ah.dao.SetProductsOwner(ctx, prdId, userId, frm.IsOwner); err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// ListProductsUsers 列出产品线所有用户
func (ah *AuthHandler) ListProductsUsers(c *gin.Context) {
	prdId, err := ext.ParamUint32(c, "product_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "product_id")
		ah.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	pFrm := form.ListProductsUsersParamsForm{}
	// 序列化request body
	if err = c.ShouldBindQuery(&pFrm); err != nil {
		ah.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in AuthHandler.ListProductsUsers, skipped it.")
	}

	// 设置查询filter
	query := pFrm.GenQuery()
	u, err := ah.dao.ListProductsUsers(tracer.ExtractTraceCtxFromGin(c), prdId, query)
	if err != nil {
		m := msg.MsgAuthDaoErr.SetError(err)
		ah.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "users", u)
}

package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/tracer"
	"eago/flow/api/form"
	"eago/flow/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewCategory 新建类别
func (h *FlowHandler) NewCategory(c *gin.Context) {
	frm := form.NewCategoryForm{}
	// 序列化request body
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 新建
	cat, err := h.dao.NewCategory(ctx, frm.Name, perm.MustGetTokenContent(c).Username)
	// 新建失败
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "category", cat)
}

// RemoveCategory 删除类别
func (h *FlowHandler) RemoveCategory(c *gin.Context) {
	catId, err := ext.ParamUint32(c, "category_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "category_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveCategoryForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, catId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.RemoveCategory(ctx, catId); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetCategory 更新类别
func (h *FlowHandler) SetCategory(c *gin.Context) {
	catId, err := ext.ParamUint32(c, "category_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "category_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetCategoryForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao, catId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	cat, err := h.dao.SetCategory(ctx, catId, frm.Name, perm.MustGetTokenContent(c).Username)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "category", cat)
}

// ListCategories 列出所有类别
func (h *FlowHandler) ListCategories(c *gin.Context) {
	pFrm := form.ListCategoriesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	cats, err := h.dao.ListCategories(tracer.ExtractTraceCtxFromGin(c), pFrm.GenQuery())
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "categories", cats)
}

// ListCategoryFlows 列出指定类别中所关联流程
func (h *FlowHandler) ListCategoryFlows(c *gin.Context) {
	catId, err := ext.ParamUint32(c, "category_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "category_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListCategoriesRelations{}
	// 设置查询filter
	if m := frm.Validate(ctx, h.dao, catId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = c.ShouldBindQuery(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	query := frm.GenQuery()
	query["categories_id=?"] = catId

	flows, err := h.dao.ListFlows(ctx, query)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "flows", flows)
}

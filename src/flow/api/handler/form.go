package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/flow/api/form"
	"eago/flow/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewForm 新建表单
func (h *FlowHandler) NewForm(c *gin.Context) {
	frm := form.NewFormForm{}
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
	frmObj, err := h.dao.NewForm(ctx, frm.Name, *frm.Disabled, *frm.Description, *frm.Body, perm.MustGetTokenContent(c).Username)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "form", frmObj)
}

// SetForm 更新表单
func (h *FlowHandler) SetForm(c *gin.Context) {
	frmId, err := ext.ParamUint32(c, "form_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "form_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.SetFormForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao, frmId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frmObj, err := h.dao.SetForm(ctx, frmId, frm.Name, *frm.Disabled, *frm.Description, perm.MustGetTokenContent(c).Username)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "form", frmObj)
}

// GetForm 获取指定表单
func (h *FlowHandler) GetForm(c *gin.Context) {
	frmId, err := ext.ParamUint32(c, "form_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "form_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.GetFormForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, frmId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frmObk, err := h.dao.GetForm(ctx, orm.Query{"id=?": frmId})
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "form", frmObk)
}

// PagedListForms 列出所有表单-分页
func (h *FlowHandler) PagedListForms(c *gin.Context) {
	pFrm := form.PagedListFormsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	paged, err := h.dao.PagedListForms(
		tracer.ExtractTraceCtxFromGin(c),
		pFrm.GenQuery(),
		c.GetInt(global.GinCtxPageKey),
		c.GetInt(global.GinCtxPageSizeKey),
		c.GetStringSlice(global.GinCtxOrderByKey)...,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "forms", paged)
}

// ListFormFlows 列出表单所关联流程
func (h *FlowHandler) ListFormFlows(c *gin.Context) {
	frmId, err := ext.ParamUint32(c, "form_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "form_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListFormRelationsForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, frmId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	flows, err := h.dao.ListFlows(ctx, orm.Query{"form_id=?": frmId})
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "flows", flows)
}

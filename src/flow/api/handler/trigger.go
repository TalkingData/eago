package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/tracer"
	"eago/flow/api/form"
	"eago/flow/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewTrigger 新建触发器
func (h *FlowHandler) NewTrigger(c *gin.Context) {
	frm := form.NewTriggerForm{}

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

	tri, err := h.dao.NewTrigger(
		ctx, frm.Name, *frm.Description, frm.TaskCodename, frm.Arguments, perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "trigger", tri)
}

// RemoveTrigger 删除触发器
func (h *FlowHandler) RemoveTrigger(c *gin.Context) {
	triId, err := ext.ParamUint32(c, "trigger_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "trigger_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveTriggerForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, triId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.RemoveTrigger(ctx, triId); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetTrigger 更新触发器
func (h *FlowHandler) SetTrigger(c *gin.Context) {
	triId, err := ext.ParamUint32(c, "trigger_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "trigger_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.SetTriggerForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, triId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	tri, err := h.dao.SetTrigger(
		ctx, triId, frm.Name, *frm.Description, frm.TaskCodename, frm.Arguments, perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "trigger", tri)
}

// PagedListTriggers 列出所有触发器-分页
func (h *FlowHandler) PagedListTriggers(c *gin.Context) {
	pFrm := form.PagedListTriggersParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	paged, err := h.dao.PagedListTriggers(
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

	ext.WriteSuccessPayload(c, "triggers", paged)
}

// ListTriggersNodes 列出触发器所关联节点
func (h *FlowHandler) ListTriggersNodes(c *gin.Context) {
	triId, err := ext.ParamUint32(c, "trigger_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "trigger_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListTriggersNodesForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, triId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	nodes, err := h.dao.ListTriggersNodes(ctx, triId)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "nodes", nodes)
}

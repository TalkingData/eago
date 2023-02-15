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

// HandleInstance 处理流程实例
func (h *FlowHandler) HandleInstance(c *gin.Context) {
	instId, err := ext.ParamUint32(c, "instance_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "instance_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.HandleInstanceForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, instId, perm.MustGetTokenContent(c).Username); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 执行实际流程实例处理
	if err = h.biz.HandleInstance(ctx, frm.Instance, frm.CreatedBy, *frm.Result, frm.Content); err != nil {
		m := msg.MsgHandleInstanceErr
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "instance_id", instId)
}

// PagedListInstances 列出所有流程实例-分页
func (h *FlowHandler) PagedListInstances(c *gin.Context) {
	pFrm := form.PagedListInstancesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	h.pagedListInstances(c, pFrm.GenDefaultQuery(perm.MustGetTokenContent(c).Username))
}

// PagedListMyInstances 列出我发起的流程实例-分页
func (h *FlowHandler) PagedListMyInstances(c *gin.Context) {
	pFrm := form.PagedListInstancesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	h.pagedListInstances(c, pFrm.GenMyInstancesQuery(perm.MustGetTokenContent(c).Username))
}

// PagedListTodoInstances 列出我代办的流程实例-分页
func (h *FlowHandler) PagedListTodoInstances(c *gin.Context) {
	pFrm := form.PagedListInstancesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	h.pagedListInstances(c, pFrm.GenTodoInstancesQuery(perm.MustGetTokenContent(c).Username))
}

// PagedListDoneInstances 列出我已办的流程实例-分页
func (h *FlowHandler) PagedListDoneInstances(c *gin.Context) {
	pFrm := form.PagedListInstancesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	h.pagedListInstances(c, pFrm.GenDoneInstancesQuery(perm.MustGetTokenContent(c).Username))
}

// pagedListInstances 列出所有流程实例-分页
func (h *FlowHandler) pagedListInstances(c *gin.Context, query orm.Query) {
	paged, err := h.dao.PagedListInstances(
		tracer.ExtractTraceCtxFromGin(c),
		query,
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

	ext.WriteSuccessPayload(c, "instances", paged)
}

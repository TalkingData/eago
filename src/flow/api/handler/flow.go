package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/tracer"
	"eago/flow/api/form"
	"eago/flow/conf/msg"
	"eago/flow/dto"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// InstantiateFlow 发起流程
func (h *FlowHandler) InstantiateFlow(c *gin.Context) {
	flowId, err := ext.ParamUint32(c, "flow_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "flow_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.InstantiateFlowForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao, flowId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 反序列化FormData
	fData := make(map[string]interface{})
	if err = json.Unmarshal([]byte(*frm.FormData), &fData); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err, "无法反序列化表单数据内容")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	tc := perm.MustGetTokenContent(c)
	// 将创建人信息存入FormData
	fData[dto.InitiatorKeyUserId] = tc.UserId
	fData[dto.InitiatorKeyUsernameKey] = tc.Username
	fData[dto.InitiatorKeyPhone] = tc.Phone

	// 序列化FormData
	fDataStr, err := json.Marshal(fData)
	if err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err, "无法序列化表单数据内容")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 调用流程服务，发起流程，返回流程实例ID
	insId, err := h.biz.InstantiateFlow(ctx, frm.FormId, string(fDataStr), perm.MustGetTokenContent(c).Username)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "instance_id", insId)
}

// NewFlow 新建流程
func (h *FlowHandler) NewFlow(c *gin.Context) {
	frm := form.NewFlowForm{}
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
	flow, err := h.dao.NewFlow(
		ctx,
		frm.Name,
		frm.InstanceTitle,
		frm.CategoriesId,
		*frm.Description,
		*frm.Disabled,
		frm.FormId,
		frm.FirstNodeId,
		perm.MustGetTokenContent(c).Username,
	)
	// 新建失败
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "flow", flow)
}

// RemoveFlow 删除流程
func (h *FlowHandler) RemoveFlow(c *gin.Context) {
	flowId, err := ext.ParamUint32(c, "flow_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "flow_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveFlowForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, flowId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.RemoveFlow(ctx, flowId); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetFlow 更新流程
func (h *FlowHandler) SetFlow(c *gin.Context) {
	flowId, err := ext.ParamUint32(c, "flow_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "flow_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.SetFlowForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, flowId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	flow, err := h.dao.SetFlow(
		ctx,
		flowId,
		frm.Name,
		frm.InstanceTitle,
		frm.CategoriesId,
		*frm.Description,
		*frm.Disabled,
		frm.FormId,
		frm.FirstNodeId,
		perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "flow", flow)
}

// PagedListFlows 列出所有流程-分页
func (h *FlowHandler) PagedListFlows(c *gin.Context) {
	pFrm := form.PagedListFlowsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	paged, err := h.dao.PagedListFlows(
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

	ext.WriteSuccessPayload(c, "flows", paged)
}

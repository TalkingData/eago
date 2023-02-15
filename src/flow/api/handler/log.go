package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/flow/api/form"
	"eago/flow/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewLog 新建指定流程实例的审批日志
func (h *FlowHandler) NewLog(c *gin.Context) {
	// 获得流程实例ID
	instId, err := ext.ParamUint32(c, "instance_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "instance_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.NewLogForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao, instId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	log, err := h.dao.NewLog(ctx, instId, frm.Result, *frm.Content, perm.MustGetTokenContent(c).Username)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "log", log)
}

// ListLogs 列出指定流程实例的所有审批日志
func (h *FlowHandler) ListLogs(c *gin.Context) {
	// 获得流程实例ID
	instId, err := ext.ParamUint32(c, "instance_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "instance_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListLogForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, instId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	logs, err := h.dao.ListLogs(ctx, orm.Query{"instance_id=?": instId})
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "logs", logs)
}

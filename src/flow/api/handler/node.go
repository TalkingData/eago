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

// NewNode 新建节点
func (h *FlowHandler) NewNode(c *gin.Context) {
	frm := form.NewNodeForm{}
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
	node, err := h.dao.NewNode(
		ctx,
		frm.Name,
		frm.ParentId,
		frm.Category,
		frm.EntryCondition,
		frm.AssigneeCondition,
		frm.VisibleFields,
		frm.EditableFields,
		perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "node", node)
}

// RemoveNode 删除节点
func (h *FlowHandler) RemoveNode(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveNodeForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.RemoveNode(ctx, nodeId); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetNode 更新节点
func (h *FlowHandler) SetNode(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.SetNodeForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	node, err := h.dao.SetNode(
		ctx,
		nodeId,
		frm.Name,
		frm.ParentId,
		frm.Category,
		frm.EntryCondition,
		frm.AssigneeCondition,
		frm.VisibleFields,
		frm.EditableFields,
		perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "node", node)
}

// GetNodeChain 列出指定节点链
func (h *FlowHandler) GetNodeChain(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 查找根节点
	node, err := h.dao.GetNode(ctx, orm.Query{"id=?": nodeId})
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 找不到根节点则直接返回空
	if node == nil || node.Id < 1 {
		ext.WriteSuccessPayload(c, "chain", make(map[string]interface{}))
		return
	}

	// 将根节点转化为链结构
	root := h.dao.Node2Chain(ctx, node)
	if err = h.dao.GetNodeChain(ctx, root); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "chain", root)
}

// PagedListNodes 列出所有节点-分页
func (h *FlowHandler) PagedListNodes(c *gin.Context) {
	pFrm := form.PagedListNodesParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	paged, err := h.dao.PagedListNodes(
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

	ext.WriteSuccessPayload(c, "nodes", paged)
}

// AddTrigger2Node 添加触发器至节点
func (h *FlowHandler) AddTrigger2Node(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	frm := form.AddTrigger2NodeForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.AddNodesTrigger(ctx, nodeId, frm.TriggerId, perm.MustGetTokenContent(c).Username); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// RemoveNodesTrigger 移除节点中触发器
func (h *FlowHandler) RemoveNodesTrigger(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	triId, err := ext.ParamUint32(c, "trigger_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "trigger_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveNodesTriggerForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId, triId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = h.dao.RemoveNodesTrigger(ctx, nodeId, triId); err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// ListNodesTriggers 列出节点中所有触发器
func (h *FlowHandler) ListNodesTriggers(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListNodesRelationsForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	tris, err := h.dao.ListNodesTriggers(ctx, nodeId)
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "triggers", tris)
}

// ListNodesFlows 列出节点所关联流程
func (h *FlowHandler) ListNodesFlows(c *gin.Context) {
	nodeId, err := ext.ParamUint32(c, "node_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "node_id")
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.ListNodesRelationsForm{}
	// 验证数据
	if m := frm.Validate(ctx, h.dao, nodeId); m != nil {
		// 数据验证未通过
		h.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	flows, err := h.dao.ListFlows(ctx, orm.Query{"first_node_id=?": nodeId})
	if err != nil {
		m := msg.MsgFlowDaoErr.SetError(err)
		h.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "flows", flows)
}

package msg

import (
	cMsg "eago/common/code_msg"
)

var (
	// Trigger 1400xx
	MsgAssociatedTriggerNodeFailed = cMsg.NewCodeMsg(140000, "无法执行操作，仍有节点与该触发器关联")

	// Form 1401xx

	// Node 1402xx
	MsgAssociatedNodeFailed                = cMsg.NewCodeMsg(140200, "无法执行操作，仍有子节点与该节点关联")
	MsgAssociatedNodeTriggerFailed         = cMsg.NewCodeMsg(140201, "无法执行操作，仍有触发器与该节点关联")
	MsgAssociatedNodeFlowFailed            = cMsg.NewCodeMsg(140202, "无法执行操作，仍有流程与该节点关联")
	MsgAssociatedParentNodeNotFoundFailed  = cMsg.NewCodeMsg(140203, "无法执行操作，该父节点不存在")
	MsgAssociatedParentNodeSelfFailed      = cMsg.NewCodeMsg(140204, "无法执行操作，不能将自身节点设置为父节点")
	MsgAssociatedParentNodeDuplicateFailed = cMsg.NewCodeMsg(140205, "无法执行操作，已经有其他节点关联了该父节点")

	// Flow 1403xx

	// Category 1404xx
	MsgAssociatedCategoryFlowFailed = cMsg.NewCodeMsg(140400, "无法执行操作，仍有流程与该类别关联")

	// Instance 1405xx
	MsgHandleInstancePermDenyErr = cMsg.NewCodeMsg(140500, "没有审批权限")
	MsgHandleInstanceErr         = cMsg.NewCodeMsg(140599, "处理流程失败")

	// Others
	MsgFlowDaoErr   = cMsg.NewCodeMsg(149900, "Flow服务的DAO层发生意外，请先尝试重试，若无效请联系管理员")
	MsgFlowCacheErr = cMsg.NewCodeMsg(149901, "Flow服务的Cache发生意外，请先尝试重试，若无效请联系管理员")
)

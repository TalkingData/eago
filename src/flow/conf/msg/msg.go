package msg

import (
	m "eago/common/message"
	"net/http"
)

var (
	// Default
	InvalidUriFailed = m.NewMessage(http.StatusBadRequest, "请求错误，请确保Params中包含以下字段:")
	SerializeFailed  = m.NewMessage(http.StatusBadRequest, "无法读取提交的数据，请确保发送正确的数据格式")
	ValidateFailed   = m.NewMessage(http.StatusBadRequest, "验证失败，提交的内容不合法")
	NotFoundFailed   = m.NewMessage(http.StatusNotFound, "操作失败，找不到指定对象")
	UndefinedError   = m.NewMessage(http.StatusInternalServerError, "操作失败，请先尝试重试，若无效请联系管理员")

	// Trigger 1200xx
	AssociatedTriggerNodeFailed = m.NewMessage(120000, "无法执行操作，仍有节点与该触发器关联")

	// Form 1201xx

	// Node 1202xx
	AssociatedNodeFailed                = m.NewMessage(120200, "无法执行操作，仍有子节点与该节点关联")
	AssociatedNodeTriggerFailed         = m.NewMessage(120201, "无法执行操作，仍有触发器与该节点关联")
	AssociatedNodeFlowFailed            = m.NewMessage(120202, "无法执行操作，仍有流程与该节点关联")
	AssociatedParentNodeNotFoundFailed  = m.NewMessage(120203, "无法执行操作，该父节点不存在")
	AssociatedParentNodeSelfFailed      = m.NewMessage(120204, "无法执行操作，不能将自身节点设置为父节点")
	AssociatedParentNodeDuplicateFailed = m.NewMessage(120205, "无法执行操作，已经有其他节点关联了该父节点")

	// Flow 1203xx

	// Category 1204xx
	AssociatedCategoryFlowFailed = m.NewMessage(120400, "无法执行操作，仍有流程与该类别关联")

	// Instance 1205xx
	HandleInstancePermDenyError = m.NewMessage(120500, "没有审批权限")
)

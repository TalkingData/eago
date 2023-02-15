package dto

// 流程实例状态
const (
	InstanceStatusPanicEnd    = -200 // 系统异常
	InstanceStatusRejectedEnd = -1   // 被驳回
	InstanceStatusApprovedEnd = 0    // 审批通过
	InstanceStatusPending     = 1    // 系统处理中
	InstanceStatusRunning     = 2    // 流转中
)

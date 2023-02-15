package msg

import (
	cMsg "eago/common/code_msg"
)

var (
	// Task 1300xx
	MsgCallTaskFailed           = cMsg.NewCodeMsg(130000, "调用任务失败，请先尝试重试，若无效请联系管理员")
	MsgTaskUniqueIdDecodeFailed = cMsg.NewCodeMsg(130001, "将任务唯一Id解码为任务结果Id和分区时失败")
	MsgAssociatedScheduleFailed = cMsg.NewCodeMsg(130302, "无法执行操作，仍有计划任务与该任务关联")

	// Result 1301xx
	MsgSetResultStatusTaskEndFailed       = cMsg.NewCodeMsg(130100, "结束任务失败，任务已经结束")
	MsgSetResultStatusInvalidStatusFailed = cMsg.NewCodeMsg(130101, "变更任务结果状态失败，状态值不合法")
	MsgKillTaskPartitionNotFoundFailed    = cMsg.NewCodeMsg(130102, "结束任务失败，没有找到对应的分区")
	MsgKillTaskFailed                     = cMsg.NewCodeMsg(130103, "结束任务失败，请先尝试重试，若无效请联系管理员")
	MsgListResultsPartNotFoundFailed      = cMsg.NewCodeMsg(130104, "无法列出任务结果，没有找到对应的分区")

	// Log 1302xx
	MsgBizNewLogFailed            = cMsg.NewCodeMsg(130200, "新增任务日志失败")
	MsgNewLogStreamSendFailed     = cMsg.NewCodeMsg(130201, "新增任务日志时，返回请求结果给客户端失败")
	MsgNewLogStreamCloseFailed    = cMsg.NewCodeMsg(130202, "新增任务日志时，关闭客户端通讯流失败")
	MsgListLogsPartNotFoundFailed = cMsg.NewCodeMsg(130203, "无法列出任务结果日志，没有找到对应的分区")

	// Others
	MsgTaskDaoErr   = cMsg.NewCodeMsg(139900, "Task服务的DAO层发生意外，请先尝试重试，若无效请联系管理员")
	MsgTaskCacheErr = cMsg.NewCodeMsg(139901, "Task服务的Cache发生意外，请先尝试重试，若无效请联系管理员")
)

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

	// Task 1300xx
	CallTaskFailed           = m.NewMessage(130000, "调用任务失败，请先尝试重试，若无效请联系管理员")
	AssociatedScheduleFailed = m.NewMessage(130300, "无法执行操作，仍有计划任务与该任务关联")

	// Result 1301xx
	KillTaskPartitionNotFoundFailed = m.NewMessage(130100, "结束任务失败，没有找到对应的分区")
	KillTaskFailed                  = m.NewMessage(130101, "结束任务失败，请先尝试重试，若无效请联系管理员")
	ListResultsPartNotFoundFailed   = m.NewMessage(130102, "无法列出任务结果，没有找到对应的分区")

	// Log 1302xx
	ListLogsPartNotFoundFailed = m.NewMessage(130200, "无法列出任务结果日志，没有找到对应的分区")
)

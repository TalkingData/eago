package dto

const (
	TaskResultStatusPanicEnd                 = -255 // 任务异常
	TaskResultStatusCallErrEnd               = -202 // 调用错误
	TaskResultStatusNoWorkerErrEnd           = -201 // 找不到执行器
	TaskResultStatusWorkerTaskNotFoundErrEnd = -200 // 找不到任务
	TaskResultStatusFailedEnd                = -3   // 任务失败
	TaskResultStatusTimeoutEnd               = -2   // 任务超时
	TaskResultStatusManualEnd                = -1   // 手动结束
	TaskResultStatusSuccessEnd               = 0    // 任务结束
	TaskResultStatusInitialization           = 1    // 初始化
	TaskResultStatusPending                  = 2    // 等待中
	TaskResultStatusRunning                  = 3    // 运行中
)

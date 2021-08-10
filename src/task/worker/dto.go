package worker

const (
	TASK_PANIC_END_STATUS                       = -255 // 任务异常
	TASK_CALL_ERROR_END_STATUS                  = -202 // 调用错误
	TASK_NO_WORKER_ERROR_END_STATUS             = -201 // 找不到Worker
	TASK_WORKER_TASK_NOT_FOUND_ERROR_END_STATUS = -200 // 找不到任务
	TASK_FAILED_END_STATUS                      = -3   // 任务失败
	TASK_TIMEOUT_END_STATUS                     = -2   // 任务超时
	TASK_MANUAL_END_STATUS                      = -1   // 手动结束
	TASK_SUCCESS_END_STATUS                     = 0    // 任务结束
	TASK_INITIALIZATION_STATUS                  = 1    // 初始化
	TASK_PENDING_STATUS                         = 2    // 等待中
	TASK_RUNNING_STATUS                         = 3    // 运行中
)

const (
	WORKER_LOGGER_BUFFER_SIZE = 100
	WORKER_REGISTER_TTL       = 3
)

type CallTaskReq struct {
	TaskCodename string `json:"task_codename"`

	TaskUniqueId string `json:"task_unique_id"`
	Arguments    string `json:"arguments"`
	Timeout      int64  `json:"timeout"`
	Caller       string `json:"caller"`

	Timestamp int64 `json:"timestamp"`
}

type KillTaskReq struct {
	TaskUniqueId string `json:"task_unique_id"`
	Timestamp    int64  `json:"timestamp"`
}

type WorkerInfo struct {
	Modular   string `json:"modular"`
	Address   string `json:"address"`
	WorkerId  string `json:"worker_id"`
	StartTime string `json:"start_time"`
}

type WorkerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

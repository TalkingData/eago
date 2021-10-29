package conf

const (
	SERVICE_NAME     = "eago-task"
	RPC_REGISTER_KEY = "eago.srv.task"
	API_REGISTER_KEY = "eago.api.task"

	AUTH_RPC_REGISTER_KEY = "eago.srv.auth"

	ADMIN_ROLE_NAME = "task_admin"

	TIMESTAMP_FORMAT                = "2006-01-02 15:04:05"
	TASK_PARTITION_TIMESTAMP_FORMAT = "200601"

	CONFIG_FILE_PATHNAME = "../conf/eago_task.conf"

	TASK_UNIQUE_ID_SEPARATOR = "::"
)

// 任务类别
const (
	BUTILIN_TASK_CATEGORY = 1   // 内置任务
	BASH_TASK_CATEGORY    = 100 // Bash任务
	PYTHON_TASK_CATEGORY  = 101 // Python任务
)

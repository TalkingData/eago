package conf

import (
	"eago/common/global"
	"time"
)

type constConf struct {
	ServiceName    string
	RpcRegisterKey string
	ApiRegisterKey string

	AdminRole string

	TaskCategoryBuiltin int
	TaskCategoryBash    int
	TaskCategoryPython  int

	TaskResultPartitionTsFormat string
	TaskUniqueIdSeparator       string

	TaskLogRefreshIntervalMs time.Duration
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.TaskServiceName,
		RpcRegisterKey: global.TaskRpcRegisterKey,
		ApiRegisterKey: global.TaskApiRegisterKey,

		AdminRole: "task_admin",

		TaskCategoryBuiltin: 1,
		TaskCategoryBash:    100,
		TaskCategoryPython:  101,

		TaskResultPartitionTsFormat: "2006",
		TaskUniqueIdSeparator:       "::",

		TaskLogRefreshIntervalMs: 3000,
	}
}

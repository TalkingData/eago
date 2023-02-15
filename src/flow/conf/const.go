package conf

import (
	"eago/common/global"
)

type constConf struct {
	ServiceName    string
	RpcRegisterKey string
	ApiRegisterKey string

	AdminRole string

	InstanceNameMaxLength int
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.FlowServiceName,
		RpcRegisterKey: global.FlowRpcRegisterKey,
		ApiRegisterKey: global.FlowApiRegisterKey,

		AdminRole: "flow_admin",

		InstanceNameMaxLength: 200,
	}
}

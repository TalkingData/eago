package conf

import (
	"eago/common/global"
)

type constConf struct {
	ServiceName    string
	RpcRegisterKey string
	ApiRegisterKey string

	AdminRole string

	EagleTokenKeyPrefix string
}

func newConstConf() *constConf {
	return &constConf{
		ServiceName:    global.AuthServiceName,
		RpcRegisterKey: global.AuthRpcRegisterKey,
		ApiRegisterKey: global.AuthApiRegisterKey,

		AdminRole: "auth_admin",

		EagleTokenKeyPrefix: "eagle_auth_token",
	}
}

package cli

import (
	"eago-auth/config"
	"eago-auth/srv/proto"
	"eago-common/etcd"
	"github.com/micro/go-micro/v2"
)

var AuthClient auth.AuthService

// 启动RPC客户端
func InitCli() {
	cli := micro.NewService(
		micro.Name(config.Config.RpcServiceName),
		micro.Registry(etcd.EtcdReg),
	)

	cli.Init()

	AuthClient = auth.NewAuthService(config.Config.RpcServiceName, cli.Client())
}

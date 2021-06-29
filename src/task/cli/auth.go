package cli

import (
	"eago/task/conf"
	"eago/task/srv/proto/auth"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

var AuthClient auth.AuthService

// InitAuthCli 启动Auth RPC客户端
func InitAuthCli() {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	AuthClient = auth.NewAuthService(conf.AUTH_RPC_SERVICE_NAME, cli.Client())
}

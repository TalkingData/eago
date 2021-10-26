package cli

import (
	"eago/auth/conf"
	auth "eago/auth/srv/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

var AuthClient auth.AuthService

// InitAuthCli 启动Auth RPC客户端
func InitAuthCli() {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Conf.EtcdAddresses...),
		etcdv3.Auth(conf.Conf.EtcdUsername, conf.Conf.EtcdPassword),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	AuthClient = auth.NewAuthService(conf.RPC_REGISTER_KEY, cli.Client())
}

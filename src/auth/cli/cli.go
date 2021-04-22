package cli

import (
	"eago-auth/conf"
	"eago-auth/srv/proto"
	"eago-common/etcd"
	"github.com/micro/go-micro/v2"
	"sync"
)

var AuthClient auth.AuthService

// InitCli 启动RPC客户端
func InitCli(wg *sync.WaitGroup) {
	defer wg.Done()

	cli := micro.NewService(
		micro.Name(conf.RPC_SERVICE_NAME),
		micro.Registry(etcd.EtcdRegistry),
	)

	AuthClient = auth.NewAuthService(conf.RPC_SERVICE_NAME, cli.Client())
}

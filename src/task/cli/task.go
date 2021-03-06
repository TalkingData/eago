package cli

import (
	"eago/task/conf"
	"eago/task/srv/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

var TaskClient task.TaskService

// InitTaskCli 启动Task RPC客户端
func InitTaskCli() {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	TaskClient = task.NewTaskService(conf.RPC_SERVICE_NAME, cli.Client())
}

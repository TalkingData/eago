package cli

import (
	"eago/common/global"
	taskpb "eago/task/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

// NewTaskClient 创建Task客户端
func NewTaskClient(etcdUname, etcdPasswd string, etcdAddrs []string, cliOpt ...client.Option) taskpb.TaskService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
	)

	_ = cli.Client().Init(cliOpt...)

	return taskpb.NewTaskService(global.TaskRpcRegisterKey, cli.Client())
}

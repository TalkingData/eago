package cli

import (
	authpb "eago/auth/proto"
	"eago/common/global"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
)

// NewAuthClient 创建Auth客户端
func NewAuthClient(etcdUname, etcdPasswd string, etcdAddrs []string, cliOpt ...client.Option) authpb.AuthService {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(etcdAddrs...),
		etcdv3.Auth(etcdUname, etcdPasswd),
	)
	cli := micro.NewService(
		micro.Registry(etcdReg),
	)

	_ = cli.Client().Init(cliOpt...)

	return authpb.NewAuthService(global.AuthRpcRegisterKey, cli.Client())
}

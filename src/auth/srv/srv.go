package srv

import (
	"eago-auth/config"
	"eago-auth/srv/proto"
	"eago-common/etcd"
	"eago-common/log"
	"github.com/micro/go-micro/v2"
)

type AuthService struct{}

// 启动RPC服务
func InitSrv() {
	srv := micro.NewService(
		micro.Name(config.Config.RpcServiceName),
		micro.Registry(etcd.EtcdReg),
	)

	srv.Init()
	auth.RegisterAuthHandler(srv.Server(), &AuthService{})

	if err := srv.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

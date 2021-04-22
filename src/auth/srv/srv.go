package srv

import (
	"eago-auth/conf"
	"eago-auth/srv/proto"
	"eago-common/etcd"
	"eago-common/log"
	"github.com/micro/go-micro/v2"
	"sync"
)

type AuthService struct{}

// InitSrv 启动RPC服务
func InitSrv(wg *sync.WaitGroup) {
	defer wg.Done()

	srv := micro.NewService(
		micro.Name(conf.RPC_SERVICE_NAME),
		micro.Registry(etcd.EtcdRegistry),
	)

	_ = auth.RegisterAuthHandler(srv.Server(), &AuthService{})

	if err := srv.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

package api

import (
	"eago-auth/api/router"
	"eago-auth/config"
	"eago-auth/util/sso"
	"eago-common/etcd"
	"eago-common/log"
	"github.com/micro/go-micro/v2/web"
)

// 启用API服务
func InitApi() {
	// 初始化Crowd
	if err := sso.InitCrowd(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	router.InitEngine()

	apiV1 := web.NewService(
		web.Name(config.Config.ApiV1ServiceName),
		web.Registry(etcd.EtcdReg),
		web.Handler(router.Engine),
	)

	apiV1.Init()

	if err := apiV1.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

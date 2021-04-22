package api

import (
	"eago-auth/api/router"
	"eago-auth/conf"
	"eago-auth/util/sso"
	"eago-common/etcd"
	"eago-common/log"
	"github.com/micro/go-micro/v2/web"
	"sync"
)

// InitApi 启用API服务
func InitApi(wg *sync.WaitGroup) {
	defer wg.Done()

	// 初始化Crowd
	if err := sso.InitCrowd(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	router.InitEngine()

	apiV1 := web.NewService(
		web.Name(conf.API_SERVICE_NAME),
		web.Registry(etcd.EtcdRegistry),
		web.Handler(router.Engine),
	)

	if err := apiV1.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

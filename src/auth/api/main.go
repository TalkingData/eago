package main

import (
	"eago/auth/api/router"
	"eago/auth/conf"
	"eago/auth/model"
	"eago/auth/util/sso"
	"eago/common/log"
	"eago/common/redis"
	"fmt"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/sirupsen/logrus"
	"runtime"
)

func main() {
	// 初始化DAO
	if err := model.InitDb(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// 初始化Redis
	redis.InitRedis(
		conf.Config.RedisAddress,
		conf.Config.RedisPassword,
		conf.MODULAR_NAME,
		conf.Config.RedisDb,
	)

	// 初始化Crowd
	if err := sso.InitCrowd(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	router.InitEngine()

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	apiV1 := web.NewService(
		web.Name(conf.API_SERVICE_NAME),
		web.Registry(etcdReg),
		web.Handler(router.Engine),
		web.Version("v1"),
	)

	if err := apiV1.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logLvl, err := logrus.ParseLevel(conf.Config.LogLevel)
	if err != nil {
		panic(err)
	}
	// 加载日志设置
	err = log.InitLog(
		conf.Config.LogPath,
		conf.MODULAR_NAME,
		conf.TIMESTAMP_FORMAT,
		logLvl,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}

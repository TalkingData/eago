package main

import (
	"eago-auth/api"
	"eago-auth/cli"
	"eago-auth/conf"
	db "eago-auth/database"
	"eago-auth/srv"
	"eago-common/etcd"
	"eago-common/log"
	"eago-common/redis"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
)

func main() {
	// 初始化关系型数据库
	if err := db.InitDb(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// 初始化Redis
	redis.InitRedis(
		conf.Config.RedisAddress,
		conf.Config.RedisPassword,
		conf.APP_NAME,
		conf.Config.RedisDb,
	)

	// 初始化Etcd
	etcd.InitEtcd(
		conf.Config.EtcdAddress,
		conf.Config.EtcdUsername,
		conf.Config.EtcdPassword,
	)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go srv.InitSrv(&wg)

	wg.Add(1)
	go cli.InitCli(&wg)

	wg.Add(1)
	go api.InitApi(&wg)

	wg.Wait()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置文件
	if err := conf.InitConfig(); err != nil {
		fmt.Println("Failed to init config, error:", err.Error())
		panic(err)
	}

	logLvl, err := logrus.ParseLevel(conf.Config.LogLevel)
	if err != nil {
		panic(err)
	}
	// 加载日志设置
	err = log.InitLog(
		conf.Config.LogPath,
		conf.APP_NAME,
		conf.TIMESTAMP_FORMAT,
		logLvl,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}

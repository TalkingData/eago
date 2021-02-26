package main

import (
	"eago-auth/api"
	"eago-auth/cli"
	"eago-auth/config"
	db "eago-auth/database"
	"eago-auth/srv"
	"eago-common/etcd"
	"eago-common/log"
	"eago-common/redis"
	"fmt"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置文件
	if err := config.InitConfig(); err != nil {
		fmt.Println("Failed to init config, error:", err.Error())
		panic(err)
	}

	// 加载日志设置
	if err := log.InitLog(config.Config.LogPath, config.Config.ServiceName); err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}

func main() {
	// 初始化关系型数据库
	if err := db.InitDb(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// 初始化Redis
	redis.InitRedis(
		config.Config.RedisAddress,
		config.Config.RedisPassword,
		config.Config.RedisDb,
		config.Config.ServiceName,
	)

	// 初始化Etcd
	etcd.InitEtcd(
		config.Config.EtcdAddress,
		config.Config.EtcdUsername,
		config.Config.EtcdPassword,
	)

	go srv.InitSrv()
	go cli.InitCli()
	go api.InitApi()

	select {}
}

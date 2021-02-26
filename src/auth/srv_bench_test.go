package main

import (
	"eago-auth/config"
	db "eago-auth/database"
	"eago-auth/srv"
	"eago-common/log"
	"eago-common/redis"
	"fmt"
	"runtime"
	"testing"
)

// 测试创建Token性能
func Benchmark_NewToken(b *testing.B) {
	db.UserModel.New("bench_test", "bench_test", true)
	userObj, _ := db.UserModel.Get(&db.Query{"username=?": "bench_test"})
	tokens := make([]string, 0)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tokens = append(tokens, srv.NewToken(userObj))
	}
	b.StopTimer()

	for _, t := range tokens {
		srv.DeleteToken(t)
	}
	db.UserModel.Delete(userObj.Id)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置文件
	if err := config.InitConfig(); err != nil {
		fmt.Println("Failed to init config, error:", err.Error())
		panic(err)
	}

	// 加载日志设置
	if err := log.InitLog(config.Config.LogPath, config.SERVICE_NAME); err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}

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
}

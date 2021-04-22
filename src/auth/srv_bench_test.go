package main

import (
	"eago-auth/conf"
	db "eago-auth/database"
	"eago-auth/srv"
	"eago-common/log"
	"eago-common/redis"
	"fmt"
	"github.com/sirupsen/logrus"
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
		srv.RemoveToken(t)
	}
	db.UserModel.Remove(userObj.Id)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置文件
	if err := conf.InitConfig(); err != nil {
		fmt.Println("Failed to init config, error:", err.Error())
		panic(err)
	}

	// 加载日志设置
	err := log.InitLog(
		conf.Config.LogPath,
		conf.APP_NAME,
		conf.TIMESTAMP_FORMAT,
		logrus.DebugLevel,
	)
	if err != nil {
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
		conf.Config.RedisAddress,
		conf.Config.RedisPassword,
		conf.APP_NAME,
		conf.Config.RedisDb,
	)
}

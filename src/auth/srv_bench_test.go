package main

import (
	"eago/auth/conf"
	"eago/auth/model"
	"eago/auth/srv/local"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/redis"
	"fmt"
	"runtime"
	"testing"
)

// 测试创建Token性能
func Benchmark_NewToken(b *testing.B) {
	model.NewUser("bench_test", "bench_test", true)
	userObj, _ := model.GetUser(model.Query{"username=?": "bench_test"})
	tokens := make([]string, 0)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tokens = append(tokens, local.NewToken(userObj))
	}
	b.StopTimer()

	for _, t := range tokens {
		local.RemoveToken(t)
	}
	model.RemoveUser(userObj.Id)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载日志设置
	err := log.InitLog(
		conf.Config.LogPath,
		conf.MODULAR_NAME,
		conf.TIMESTAMP_FORMAT,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}

	// 初始化DAO
	model.SetDb(orm.InitMysql(
		conf.Config.MysqlAddress,
		conf.Config.MysqlUser,
		conf.Config.MysqlPassword,
		conf.Config.MysqlDbName,
	))

	// 初始化Redis
	redis.InitRedis(
		conf.Config.RedisAddress,
		conf.Config.RedisPassword,
		conf.MODULAR_NAME,
		conf.Config.RedisDb,
	)
}

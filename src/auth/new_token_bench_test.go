package main

import (
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/auth/srv/builtin"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/redis"
	"fmt"
	"runtime"
	"testing"
)

// 测试创建Token性能
func Benchmark_NewToken(b *testing.B) {
	dao.NewUser("bench_test", "bench_test", true)
	userObj, _ := dao.GetUser(dao.Query{"username=?": "bench_test"})
	tokens := make([]string, 0)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tokens = append(tokens, builtin.NewToken(userObj))
	}
	b.StopTimer()

	for _, t := range tokens {
		builtin.RemoveToken(t)
	}
	dao.RemoveUser(userObj.Id)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载日志设置
	err := log.InitLog(
		conf.Conf.LogPath,
		conf.SERVICE_NAME,
		conf.TIMESTAMP_FORMAT,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}

	// 初始化DAO
	dao.Init(orm.InitMysql(
		conf.Conf.MysqlAddress,
		conf.Conf.MysqlUser,
		conf.Conf.MysqlPassword,
		conf.Conf.MysqlDbName,
	))

	// 初始化Redis
	redis.InitRedis(
		conf.Conf.RedisAddress,
		conf.Conf.RedisPassword,
		conf.SERVICE_NAME,
		conf.Conf.RedisDb,
	)
}

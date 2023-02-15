package main

import (
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/redis"
	"eago/common/service"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	auth service.EagoSrv

	authDao   *dao.Dao
	authRedis *redis.RedisTool

	authConf *conf.Conf
	authLg   *logger.Logger
)

func main() {
	auth = NewAuthApi(authDao, authRedis, authConf, authLg)

	e := make(chan error)
	go func() {
		e <- auth.Start()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			if err != nil {
				authLg.ErrorWithFields(logger.Fields{
					"error": err,
				}, "An error occurred while Start.")
			}
			closeAll()
			return
		case sig := <-quit:
			authLg.InfoWithFields(logger.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll
func closeAll() {
	if auth != nil {
		auth.Stop()
	}
	if authLg != nil {
		authLg.Close()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	authConf = conf.NewConfig()

	// 生成Logger
	lg, err := logger.NewLogger(
		logger.LogLevel(authConf.LogLevel),
		logger.LogPath(authConf.LogPath),
		logger.Filename(authConf.Const.ServiceName, "api"),
	)
	if err != nil {
		fmt.Println("An error occurred while logger.NewLogger, error:", err.Error())
		panic(err)
	}
	authLg = lg

	authDao = dao.NewDao(orm.NewMysqlGorm(
		authConf.MysqlAddress,
		authConf.MysqlUser,
		authConf.MysqlPassword,
		authConf.MysqlDbName,
		orm.MysqlMaxIdleConns(authConf.MysqlMaxIdleConns),
		orm.MysqlMaxOpenConns(authConf.MysqlMaxOpenConns),
		orm.UsingOpentracingPlugin(),
	), authLg)

	authRedis = redis.NewRedisTool(
		authConf.RedisAddress,
		authConf.RedisPassword,
		authConf.Const.ServiceName,
		authConf.RedisDb,
		redis.UsingOpentracingHook(),
	)
}

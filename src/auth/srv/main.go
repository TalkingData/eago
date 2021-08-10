package main

import (
	"eago/auth/conf"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/redis"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type AuthService struct{}

func main() {
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

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	srv := micro.NewService(
		micro.Name(conf.RPC_SERVICE_NAME),
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	_ = auth.RegisterAuthServiceHandler(srv.Server(), &AuthService{})

	if err := srv.Run(); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	e := make(chan error)
	go func() {
		e <- srv.Run()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			log.ErrorWithFields(log.Fields{
				"error": err.Error(),
			}, "Error when srv.Run.")
			closeAll()
			return
		case sig := <-quit:
			log.InfoWithFields(log.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			closeAll()
			return
		}
	}
}

// closeAll 关闭全部
func closeAll() {
	orm.Close()
	redis.Close()
	log.Close()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载日志设置
	err := log.InitLog(
		conf.Config.LogPath,
		conf.MODULAR_NAME,
		conf.Config.LogLevel,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}

package main

import (
	"eago/common/log"
	"eago/common/orm"
	"eago/common/redis"
	"eago/task/cli"
	"eago/task/conf"
	"eago/task/model"
	task "eago/task/srv/proto"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type TaskService struct{}

func main() {
	// 初始化DAO
	model.SetDb(orm.InitMysql(
		conf.Config.MysqlAddress,
		conf.Config.MysqlUser,
		conf.Config.MysqlPassword,
		conf.Config.MysqlDbName,
	))

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	srv := micro.NewService(
		micro.Name(conf.RPC_SERVICE_NAME),
		micro.Registry(etcdReg),
		micro.Version("v1"),
	)

	_ = task.RegisterTaskServiceHandler(srv.Server(), &TaskService{})

	cli.InitWorkerCli()

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

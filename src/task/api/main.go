package main

import (
	"eago/common/log"
	"eago/common/orm"
	"eago/task/cli"
	"eago/task/conf"
	"eago/task/model"
	"fmt"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	// 初始化DAO
	model.SetDb(orm.InitMysql(
		conf.Config.MysqlAddress,
		conf.Config.MysqlUser,
		conf.Config.MysqlPassword,
		conf.Config.MysqlDbName,
	))

	cli.InitAuthCli()
	cli.InitTaskCli()

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Config.EtcdAddresses...),
		etcdv3.Auth(conf.Config.EtcdUsername, conf.Config.EtcdPassword),
	)
	apiV1 := web.NewService(
		web.Name(conf.API_SERVICE_NAME),
		web.Registry(etcdReg),
		web.Handler(Engine),
		web.Version("v1"),
	)

	// 初始化WorkerCli
	cli.InitWorkerCli()

	e := make(chan error)
	go func() {
		e <- apiV1.Run()
	}()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case err := <-e:
			log.ErrorWithFields(log.Fields{
				"error": err.Error(),
			}, "Error when apiV1.Run.")
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

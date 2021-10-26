package main

import (
	"context"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/task/cli"
	"eago/task/conf"
	"eago/task/dao"
	task "eago/task/srv/proto"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type TaskService struct{}

func main() {
	// 初始化DAO
	dao.Init(orm.InitMysql(
		conf.Conf.MysqlAddress,
		conf.Conf.MysqlUser,
		conf.Conf.MysqlPassword,
		conf.Conf.MysqlDbName,
	))

	t, c := tracer.NewTracer(conf.RPC_REGISTER_KEY, conf.Conf.JaegerAddress)
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Conf.EtcdAddresses...),
		etcdv3.Auth(conf.Conf.EtcdUsername, conf.Conf.EtcdPassword),
	)
	srv := micro.NewService(
		micro.Name(conf.RPC_REGISTER_KEY),
		micro.Registry(etcdReg),
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		micro.Context(ctx),
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
			if err != nil {
				log.ErrorWithFields(log.Fields{
					"error": err,
				}, "An error occurred while srv.Run.")
			}
			closeAll()
			return
		case sig := <-quit:
			log.InfoWithFields(log.Fields{
				"signal": sig.String(),
			}, "Got quit signal.")
			cancel()
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
		conf.Conf.LogPath,
		conf.SERVICE_NAME,
		conf.Conf.LogLevel,
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}

package main

import (
	"context"
	perm "eago/common/api-suite/permission"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/flow/cli"
	"eago/flow/conf"
	"eago/flow/dao"
	"fmt"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	// 初始化DAO
	dao.Init(orm.InitMysql(
		conf.Conf.MysqlAddress,
		conf.Conf.MysqlUser,
		conf.Conf.MysqlPassword,
		conf.Conf.MysqlDbName,
	))

	// 初始化Tracer
	t, c := tracer.NewTracer(conf.API_REGISTER_KEY, conf.Conf.JaegerAddress)
	defer c.Close()

	opentracing.SetGlobalTracer(t)

	cli.InitAuthCli()
	cli.InitTaskCli()

	perm.SetAuthClient(cli.AuthClient)

	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Conf.EtcdAddresses...),
		etcdv3.Auth(conf.Conf.EtcdUsername, conf.Conf.EtcdPassword),
	)
	apiV1 := web.NewService(
		web.Name(conf.API_REGISTER_KEY),
		web.Registry(etcdReg),
		web.Handler(NewGinEngine()),
		web.Context(ctx),
		web.Version("v1"),
	)

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
package main

import (
	"context"
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/auth/srv/proto"
	"eago/common/broker"
	"eago/common/log"
	"eago/common/orm"
	"eago/common/redis"
	"eago/common/tracer"
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

type AuthService struct{}

var Publisher broker.Publisher

func main() {
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

	t, c := tracer.NewTracer(conf.RPC_REGISTER_KEY, conf.Conf.JaegerAddress)
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.Conf.EtcdAddresses...),
		etcdv3.Auth(conf.Conf.EtcdUsername, conf.Conf.EtcdPassword),
	)
	srv := micro.NewService(
		micro.Name(conf.RPC_REGISTER_KEY),
		micro.Address(conf.Conf.SrvListen),
		micro.Version("v1"),
		micro.Registry(etcdReg),
		micro.RegisterTTL(conf.Conf.RegisterTtl),
		micro.RegisterInterval(conf.Conf.RegisterInterval),
		micro.Context(ctx),
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		micro.Broker(broker.NewBroker(conf.Conf.KafkaAddresses)),
	)

	// 初始化broker
	if Publisher == nil {
		Publisher = broker.NewPublisher(conf.SERVICE_NAME)
	}

	_ = auth.RegisterAuthServiceHandler(srv.Server(), &AuthService{})

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
	redis.Close()
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

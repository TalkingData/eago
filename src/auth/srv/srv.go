package main

import (
	"context"
	"eago/auth/biz"
	"eago/auth/conf"
	"eago/auth/dao"
	authpb "eago/auth/proto"
	"eago/auth/srv/service"
	"eago/common/broker"
	"eago/common/logger"
	"eago/common/redis"
	"eago/common/tracer"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	opentracingGo "github.com/opentracing/opentracing-go"
)

type authSrv struct {
	srv micro.Service

	dao   *dao.Dao
	redis *redis.RedisTool
	pub   broker.Publisher

	biz *biz.Biz

	conf   *conf.Conf
	logger *logger.Logger
	tracer tracer.Tracer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewAuthSrv(dao *dao.Dao, redis *redis.RedisTool, conf *conf.Conf, logger *logger.Logger) *authSrv {
	ctx, cancel := context.WithCancel(context.Background())

	// 生成etcdRegistry
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.EtcdAddresses...),
		etcdv3.Auth(conf.EtcdUsername, conf.EtcdPassword),
	)

	// 生成Tracer
	_tracer, err := tracer.NewJaegerTracer(
		tracer.RegisterKey(conf.Const.RpcRegisterKey),
		tracer.JaegerHostPort(conf.JaegerAddress),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	opentracingGo.SetGlobalTracer(_tracer.GetTracer())

	// 生成Broker
	_broker, err := broker.NewKafkaBroker(conf.KafkaAddresses)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init KafkaBroker: %v\n", err))
	}

	// 生成Srv
	_srv := micro.NewService(
		micro.Name(conf.Const.RpcRegisterKey),
		micro.Address(conf.SrvListen),
		micro.Registry(etcdReg),
		micro.RegisterTTL(conf.MicroRegisterTtl),
		micro.RegisterInterval(conf.MicroRegisterInterval),
		micro.Context(ctx),
		micro.WrapHandler(opentracing.NewHandlerWrapper(_tracer.GetTracer())),
		micro.Broker(_broker),
	)

	// 生成Publisher
	_pub := broker.NewPublisher(
		_broker,
		broker.ServiceName(conf.Const.ServiceName),
		broker.Logger(logger),
	)

	// 生成Biz
	_biz := biz.NewBiz(dao, redis, _pub, conf, logger)

	_ = authpb.RegisterAuthServiceHandler(_srv.Server(), service.NewAuthService(dao, redis, _biz, conf, logger))

	return &authSrv{
		srv: _srv,

		dao:   dao,
		redis: redis,
		pub:   _pub,

		biz: _biz,

		conf:   conf,
		logger: logger,
		tracer: _tracer,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (as *authSrv) Start() error {
	as.logger.Info("Starting auth srv ...")

	return as.srv.Run()
}

func (as *authSrv) Stop() {
	if as.cancelFunc != nil {
		as.cancelFunc()
	}
}

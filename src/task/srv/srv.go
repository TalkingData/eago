package main

import (
	"context"
	"eago/common/logger"
	"eago/common/redis"
	"eago/common/tracer"
	"eago/task/biz"
	"eago/task/conf"
	"eago/task/dao"
	taskpb "eago/task/proto"
	"eago/task/srv/service"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	opentracingGo "github.com/opentracing/opentracing-go"
)

type taskSrv struct {
	srv micro.Service

	dao   *dao.Dao
	redis *redis.RedisTool

	biz *biz.Biz

	conf   *conf.Conf
	logger *logger.Logger
	tracer tracer.Tracer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTaskSrv(dao *dao.Dao, redis *redis.RedisTool, conf *conf.Conf, logger *logger.Logger) *taskSrv {
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

	// 生成Srv
	_srv := micro.NewService(
		micro.Name(conf.Const.RpcRegisterKey),
		micro.Address(conf.SrvListen),
		micro.Registry(etcdReg),
		micro.RegisterTTL(conf.MicroRegisterTtl),
		micro.RegisterInterval(conf.MicroRegisterInterval),
		micro.Context(ctx),
		micro.WrapHandler(opentracing.NewHandlerWrapper(_tracer.GetTracer())),
	)

	// 生成Biz
	_biz := biz.NewBiz(dao, redis, conf, logger)

	_ = taskpb.RegisterTaskServiceHandler(_srv.Server(), service.NewTaskService(dao, redis, _biz, conf, logger))

	return &taskSrv{
		srv: _srv,

		dao:   dao,
		redis: redis,

		biz: _biz,

		conf:   conf,
		logger: logger,
		tracer: _tracer,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (ts *taskSrv) Start() error {
	ts.logger.Info("Starting task srv ...")

	return ts.srv.Run()
}

func (ts *taskSrv) Stop() {
	if ts.cancelFunc != nil {
		ts.cancelFunc()
	}
}

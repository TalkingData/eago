package main

import (
	"context"
	"eago/common/api"
	perm "eago/common/api/permission"
	"eago/common/logger"
	"eago/common/service"
	"eago/common/tracer"
	"eago/task/api/handler"
	"eago/task/conf"
	"eago/task/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
)

type taskApi struct {
	api web.Service

	handler *handler.TaskHandler

	conf   *conf.Conf
	logger *logger.Logger
	tracer tracer.Tracer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTaskApi(dao *dao.Dao, conf *conf.Conf, logger *logger.Logger) service.EagoSrv {
	ctx, cancel := context.WithCancel(context.Background())

	// 生成etcdRegistry
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(conf.EtcdAddresses...),
		etcdv3.Auth(conf.EtcdUsername, conf.EtcdPassword),
	)

	// 生成Tracer
	_tracer, err := tracer.NewJaegerTracer(
		tracer.RegisterKey(conf.Const.ApiRegisterKey),
		tracer.JaegerHostPort(conf.JaegerAddress),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	opentracing.SetGlobalTracer(_tracer.GetTracer())

	_handler := handler.NewTaskHandler(dao, conf, logger)

	_api := web.NewService(
		web.Name(conf.Const.ApiRegisterKey),
		web.Address(conf.ApiListen),
		web.Handler(newGinEngine(conf.GinMode, conf, logger, _handler)),
		web.Registry(etcdReg),
		web.RegisterTTL(conf.MicroRegisterTtl),
		web.RegisterInterval(conf.MicroRegisterInterval),
		web.Context(ctx),
	)

	return &taskApi{
		api: _api,

		handler: _handler,

		conf:   conf,
		logger: logger,
		tracer: _tracer,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (ta *taskApi) Start() error {
	ta.logger.Info("Starting task api ...")

	return ta.api.Run()
}

func (ta *taskApi) Stop() {
	if ta.cancelFunc != nil {
		ta.cancelFunc()
	}
}

func newGinEngine(ginMode string, _conf *conf.Conf, logger *logger.Logger, h *handler.TaskHandler) *gin.Engine {
	gin.SetMode(ginMode)

	engine := gin.New()
	engine.Use(api.GinCustomLogger(logger), gin.Recovery(), api.OpentracingMiddleware)

	g := engine.Group("/task", perm.MustLogin(h.GetAuthCli()))
	{
		// 根据当前登录用户权限列出菜单
		g.GET("/menus", h.ListMenus)

		// 列出所有Worker
		g.GET("/workers", perm.MustRole(_conf.Const.AdminRole), h.ListWorkers)

		// Task模块
		tr := g.Group("/tasks")
		{
			tr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewTask)
			tr.DELETE("/:task_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveTask)
			tr.PUT("/:task_id", perm.MustRole(_conf.Const.AdminRole), h.SetTask)
			tr.GET("", perm.MustRole(_conf.Const.AdminRole), api.PagingQueryMiddleware, h.PagedListTasks)

			// 调用任务
			tr.POST("/:task_id/call", perm.MustRole(_conf.Const.AdminRole), h.CallTask)
		}

		// Schedule模块
		sr := g.Group("/schedules")
		{
			sr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewSchedule)
			sr.DELETE("/:schedule_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveSchedule)
			sr.PUT("/:schedule_id", perm.MustRole(_conf.Const.AdminRole), h.SetSchedule)
			sr.GET("", perm.MustRole(_conf.Const.AdminRole), api.PagingQueryMiddleware, h.PagedListSchedules)
		}

		// ResultTables模块
		rpr := g.Group("/result_partitions")
		{
			// 仅测试用，不对外开放的方法
			// 新建结果分区并建立结果表和日志表
			rpr.POST(
				"/with_create_tables",
				perm.MustRole(_conf.Const.AdminRole),
				h.NewResultPartitionsWithCreateTables,
			)

			// 列出所有结果分区
			rpr.GET("", h.ListResultPartitions)
		}

		// Result模块
		rr := g.Group("/results")
		{
			// 按分区ID列出所有结果
			rr.GET("/:result_partition_id", api.PagingQueryMiddleware, h.PagedListResults)
			// 手动结束任务
			rr.DELETE("/:result_partition_id/:result_id", perm.MustRole(_conf.Const.AdminRole), h.KillTask)

			// 按任务唯一ID查询结果
			rr.GET("/task_unique_id/:task_unique_id", h.GetResultByTaskUniqueId)
			// 按任务唯一ID手动结束任务
			rr.DELETE(
				"/task_unique_id/:task_unique_id",
				perm.MustRole(_conf.Const.AdminRole),
				h.KillTaskByTaskUniqueId,
			)
		}

		// Log模块
		// 按分区ID列出所有结果日志
		g.GET("/logs/:result_partition_id/:result_id", h.ListLogs)
	}

	// Log模块
	// 以WebSocket方式按分区ID列出所有结果日志
	engine.GET(
		"/task/logs/:result_partition_id/:result_id/ws",
		perm.MustLoginWs(h.GetAuthCli()),
		h.WsListLogs,
	)

	engine.NoRoute(api.PageNotFound)

	return engine
}

package main

import (
	"context"
	"eago/common/api"
	perm "eago/common/api/permission"
	"eago/common/broker"
	"eago/common/logger"
	"eago/common/service"
	"eago/common/tracer"
	"eago/flow/api/handler"
	"eago/flow/biz"
	"eago/flow/conf"
	"eago/flow/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
)

type flowApi struct {
	api web.Service

	handler *handler.FlowHandler

	conf   *conf.Conf
	logger *logger.Logger
	tracer tracer.Tracer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewFlowApi(dao *dao.Dao, conf *conf.Conf, logger *logger.Logger) service.EagoSrv {
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

	// 生成Broker
	_broker, err := broker.NewKafkaBroker(conf.KafkaAddresses)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init KafkaBroker: %v\n", err))
	}

	// 生成Publisher
	_pub := broker.NewPublisher(
		_broker,
		broker.ServiceName(conf.Const.ServiceName),
		broker.Logger(logger),
	)

	_handler := handler.NewFlowHandler(dao, biz.NewBiz(dao, _pub, conf, logger), conf, logger)

	_api := web.NewService(
		web.Name(conf.Const.ApiRegisterKey),
		web.Address(conf.ApiListen),
		web.Handler(newGinEngine(conf.GinMode, conf, logger, _handler)),
		web.Registry(etcdReg),
		web.RegisterTTL(conf.MicroRegisterTtl),
		web.RegisterInterval(conf.MicroRegisterInterval),
		web.Context(ctx),
	)

	return &flowApi{
		api: _api,

		handler: _handler,

		conf:   conf,
		logger: logger,
		tracer: _tracer,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (fa *flowApi) Start() error {
	fa.logger.Info("Starting flow api ...")

	return fa.api.Run()
}

func (fa *flowApi) Stop() {
	if fa.cancelFunc != nil {
		fa.cancelFunc()
	}
}

func newGinEngine(ginMode string, _conf *conf.Conf, logger *logger.Logger, h *handler.FlowHandler) *gin.Engine {
	gin.SetMode(ginMode)

	engine := gin.New()
	engine.Use(api.GinCustomLogger(logger), gin.Recovery(), api.OpentracingMiddleware)

	fGroup := engine.Group("/flow", perm.MustLogin(h.GetAuthCli()))
	{
		// 根据当前登录用户权限列出菜单
		fGroup.GET("/menus", h.ListMenus)

		// Instance流程实例模块
		iR := fGroup.Group("/instances")
		{
			// 处理指定流程实例
			iR.PUT("/:instance_id/handle", h.HandleInstance)

			// 列出我发起的流程实例
			iR.GET("/my", api.PagingQueryMiddleware, h.PagedListMyInstances)
			// 列出我代办的流程实例
			iR.GET("/todo", api.PagingQueryMiddleware, h.PagedListTodoInstances)
			// 列出我已办的流程实例
			iR.GET("/done", api.PagingQueryMiddleware, h.PagedListDoneInstances)

			// 列出所有流程实例，要求管理员权限
			iR.GET("", perm.MustRole(_conf.Const.AdminRole), api.PagingQueryMiddleware, h.PagedListInstances)
		}

		// Log审批日志模块
		lR := fGroup.Group("/logs")
		{
			// 新增指定流程实例审批日志
			lR.POST("/:instance_id", h.NewLog)
			// 列出指定流程实例审批日志
			lR.GET("/:instance_id", h.ListLogs)
		}

		// Category类别模块
		cR := fGroup.Group("/categories")
		{
			// 列出所有类别
			cR.GET("", h.ListCategories)

			// 列出指定类别中所关联流程
			cR.GET("/:category_id/flows", h.ListCategoryFlows)

			// 新增类别，要求管理员权限
			cR.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewCategory)
			// 删除类别，要求管理员权限
			cR.DELETE("/:category_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveCategory)
			// 更新类别，要求管理员权限
			cR.PUT("/:category_id", perm.MustRole(_conf.Const.AdminRole), h.SetCategory)
		}

		// Flow流程模块
		flR := fGroup.Group("/flows")
		{
			// 发起流程
			flR.POST("/:flow_id/instantiate", h.InstantiateFlow)

			// 分页列出所有流程
			flR.GET("", api.PagingQueryMiddleware, h.PagedListFlows)

			// 新建流程，要求管理员权限
			flR.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewFlow)
			// 删除流程，要求管理员权限
			flR.DELETE("/:flow_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveFlow)
			// 更新流程，要求管理员权限
			flR.PUT("/:flow_id", perm.MustRole(_conf.Const.AdminRole), h.SetFlow)
		}

		// Form表单模块
		fR := fGroup.Group("/forms")
		{
			// 获取指定表单
			fR.GET("/:form_id", h.GetForm)

			// 新建表单，要求管理员权限
			fR.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewForm)
			// 变更表单，要求管理员权限
			fR.PUT("/:form_id", perm.MustRole(_conf.Const.AdminRole), h.SetForm)
			// 列出所有表单，要求管理员权限
			fR.GET("", perm.MustRole(_conf.Const.AdminRole), api.PagingQueryMiddleware, h.PagedListForms)

			// 列出指定表单所关联流程，要求管理员权限
			fR.GET("/:form_id/flows", perm.MustRole(_conf.Const.AdminRole), h.ListFormFlows)
		}

		// Node节点模块
		nR := fGroup.Group("/nodes")
		{
			// 列出指定节点链
			nR.GET("/:node_id/chain", h.GetNodeChain)

			// 新增节点，要求管理员权限
			nR.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewNode)
			// 删除节点，要求管理员权限
			nR.DELETE("/:node_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveNode)
			// 更新节点，要求管理员权限
			nR.PUT("/:node_id", perm.MustRole(_conf.Const.AdminRole), h.SetNode)
			// 列出所有节点，要求管理员权限
			nR.GET("", perm.MustRole(_conf.Const.AdminRole), api.PagingQueryMiddleware, h.PagedListNodes)

			// 添加触发器至节点，要求管理员权限
			nR.POST("/:node_id/triggers", perm.MustRole(_conf.Const.AdminRole), h.AddTrigger2Node)
			// 移除指定节点中触发器，要求管理员权限
			nR.DELETE("/:node_id/triggers/:trigger_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveNodesTrigger)
			// 列出指定节点中所有触发器，要求管理员权限
			nR.GET("/:node_id/triggers", perm.MustRole(_conf.Const.AdminRole), h.ListNodesTriggers)

			// 列出指定节点所关联流程，要求管理员权限
			nR.GET("/:node_id/flows", perm.MustRole(_conf.Const.AdminRole), h.ListNodesFlows)
		}

		// Trigger触发器模块
		tR := fGroup.Group("/triggers", perm.MustRole(_conf.Const.AdminRole))
		{
			// 新增触发器，要求管理员权限
			tR.POST("", h.NewTrigger)
			// 删除触发器，要求管理员权限
			tR.DELETE("/:trigger_id", h.RemoveTrigger)
			// 变更触发器，要求管理员权限
			tR.PUT("/:trigger_id", h.SetTrigger)
			// 列出所有触发器，要求管理员权限
			tR.GET("", api.PagingQueryMiddleware, h.PagedListTriggers)

			// 列出指定触发器所关联节点，要求管理员权限
			tR.GET("/:trigger_id/nodes", perm.MustRole(_conf.Const.AdminRole), h.ListTriggersNodes)
		}
	}

	engine.NoRoute(api.PageNotFound)

	return engine
}

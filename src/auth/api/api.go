package main

import (
	"context"
	"eago/auth/api/handler"
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/common/api"
	perm "eago/common/api/permission"
	"eago/common/logger"
	"eago/common/redis"
	"eago/common/service"
	"eago/common/tracer"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
)

type authApi struct {
	api web.Service

	handler *handler.AuthHandler

	conf   *conf.Conf
	logger *logger.Logger
	tracer tracer.Tracer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewAuthApi(dao *dao.Dao, redis *redis.RedisTool, conf *conf.Conf, logger *logger.Logger) service.EagoSrv {
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

	_handler := handler.NewAuthHandler(dao, redis, conf, logger)

	_api := web.NewService(
		web.Name(conf.Const.ApiRegisterKey),
		web.Address(conf.ApiListen),
		web.Handler(newGinEngine(conf.GinMode, conf, logger, _handler)),
		web.Registry(etcdReg),
		web.RegisterTTL(conf.MicroRegisterTtl),
		web.RegisterInterval(conf.MicroRegisterInterval),
		web.Context(ctx),
	)

	return &authApi{
		api: _api,

		handler: _handler,

		conf:   conf,
		logger: logger,
		tracer: _tracer,

		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (aa *authApi) Start() error {
	aa.logger.Info("Starting auth api ...")

	return aa.api.Run()
}

func (aa *authApi) Stop() {
	if aa.cancelFunc != nil {
		aa.cancelFunc()
	}
}

func newGinEngine(ginMode string, _conf *conf.Conf, logger *logger.Logger, h *handler.AuthHandler) *gin.Engine {
	gin.SetMode(ginMode)

	//engine := gin.Default()
	//engine.Use(api.OpentracingMiddleware)
	engine := gin.New()
	engine.Use(api.GinCustomLogger(logger), gin.Recovery(), api.OpentracingMiddleware)

	// 登录
	engine.POST("/auth/login", h.ReadLoginForm, h.CrowdLogin, h.DatabaseLogin, h.LoginFailed)
	engine.POST("/auth/login/eagle", h.LoginFromEagle)
	engine.GET("/auth/token/content", h.GetTokenContent)

	g := engine.Group("/auth", perm.MustLogin(h.GetAuthCli()))
	{
		g.POST("heartbeat", h.Heartbeat)
		g.DELETE("logout", h.Logout)

		// 根据当前登录用户权限列出菜单
		g.GET("/menus", h.ListMenus)

		// Product模块
		pr := g.Group("/products")
		{
			// 新建产品线
			pr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewProduct)
			// 删除产品线
			pr.DELETE("/:product_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveProduct)
			// 更新产品线
			pr.PUT("/:product_id", perm.MustRole(_conf.Const.AdminRole), h.SetProduct)
			// 列出所有产品线
			pr.GET("", api.PagingQueryMiddleware, h.PagedListProducts)

			// 添加用户至产品线
			pr.POST(
				"/:product_id/users",
				perm.MustCurrUserInProductOrRole("product_id", _conf.Const.AdminRole, true),
				h.AddUser2Product,
			)
			// 移除产品线中用户
			pr.DELETE(
				"/:product_id/users/:user_id",
				perm.MustCurrUserInProductOrRole("product_id", _conf.Const.AdminRole, true),
				h.RemoveProductsUser,
			)
			// 设置用户是否是产品线Owner
			pr.PUT("/:product_id/users/:user_id", perm.MustRole(_conf.Const.AdminRole), h.SetProductsOwner)
			// 列出产品线中所有用户
			pr.GET("/:product_id/users", h.ListProductsUsers)
		}

		// Department模块
		dr := g.Group("/departments")
		{
			// 新建部门
			dr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewDepartment)
			// 删除部门
			dr.DELETE("/:department_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveDepartment)
			// 更新部门
			dr.PUT("/:department_id", perm.MustRole(_conf.Const.AdminRole), h.SetDepartment)
			// 列出所有部门
			dr.GET("", api.PagingQueryMiddleware, h.PagedListDepartments)
			// 列出指定部门子树
			dr.GET("/:department_id/tree", h.ListDepartmentTree)
			// 以树结构列出所有部门
			dr.GET("/tree", h.ListDepartmentFullTree)

			// 添加部门至角色
			dr.POST("/:department_id/users", perm.MustRole(_conf.Const.AdminRole), h.AddUser2Department)
			// 移除部门中用户
			dr.DELETE("/:department_id/users/:user_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveDepartmentsUser)
			// 设置用户是否是部门Owner
			dr.PUT("/:department_id/users/:user_id", perm.MustRole(_conf.Const.AdminRole), h.SetDepartmentsOwner)
			// 列出部门中所有用户
			dr.GET("/:department_id/users", h.ListDepartmentsUsers)
		}

		// Group模块
		gr := g.Group("/groups")
		{
			// 新建组
			gr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewGroup)
			// 删除组
			gr.DELETE("/:group_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveGroup)
			// 更新组
			gr.PUT("/:group_id", perm.MustRole(_conf.Const.AdminRole), h.SetGroup)
			// 列出所有组
			gr.GET("", api.PagingQueryMiddleware, h.PagedListGroups)

			// 添加用户至组
			gr.POST(
				"/:group_id/users",
				perm.MustCurrUserInGroupOrRole("group_id", _conf.Const.AdminRole, true),
				h.AddUser2Group,
			)
			// 移除组中用户
			gr.DELETE(
				"/:group_id/users/:user_id",
				perm.MustCurrUserInGroupOrRole("group_id", _conf.Const.AdminRole, true),
				h.RemoveGroupsUser,
			)
			// 设置用户是否是组Owner
			gr.PUT("/:group_id/users/:user_id", perm.MustRole(_conf.Const.AdminRole), h.SetGroupsOwner)
			// 列出组中所有用户
			gr.GET("/:group_id/users", h.ListGroupsUsers)
		}

		// Role模块
		rr := g.Group("/roles")
		{
			// 新建角色
			rr.POST("", perm.MustRole(_conf.Const.AdminRole), h.NewRole)
			// 删除角色
			rr.DELETE("/:role_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveRole)
			// 更新角色
			rr.PUT("/:role_id", perm.MustRole(_conf.Const.AdminRole), h.SetRole)
			// 列出所有角色
			rr.GET("", api.PagingQueryMiddleware, h.PagedListRoles)

			// 添加用户至角色
			rr.POST("/:role_id/users", perm.MustRole(_conf.Const.AdminRole), h.AddUser2Role)
			// 移除角色中用户
			rr.DELETE("/:role_id/users/:user_id", perm.MustRole(_conf.Const.AdminRole), h.RemoveRolesUser)
			// 列出角色所有用户
			rr.GET("/:role_id/users", h.ListRolesUsers)
		}

		// User模块
		ur := g.Group("/users")
		{
			// 更新用户
			ur.PUT("/:user_id",
				perm.MustCurrUserOrRole("user_id", _conf.Const.AdminRole),
				h.SetUser)
			// 列出所有用户-分页
			ur.GET("", api.PagingQueryMiddleware, h.PagedListUsers)

			// 列出用户所有角色
			ur.GET("/:user_id/roles",
				perm.MustCurrUserOrRole("user_id", _conf.Const.AdminRole),
				h.ListUsersRoles)
			// 列出用户所有产品线
			ur.GET("/:user_id/products",
				perm.MustCurrUserOrRole("user_id", _conf.Const.AdminRole),
				h.ListUsersProducts)
			// 列出用户所有组
			ur.GET("/:user_id/groups",
				perm.MustCurrUserOrRole("user_id", _conf.Const.AdminRole),
				h.ListUsersGroups)
			// 获得指定用户所在部门
			ur.GET("/:user_id/department",
				perm.MustCurrUserOrRole("user_id",
					_conf.Const.AdminRole), h.GetUsersDepartment)
			// 获得指定用户所在部门链，包含所有层级情况
			ur.GET("/:user_id/department/chain",
				perm.MustCurrUserOrRole("user_id", _conf.Const.AdminRole),
				h.GetUsersDepartmentChain)

			// 用户交接
			ur.GET("/handover/:user_id/:target_user_id",
				perm.MustRole(_conf.Const.AdminRole),
				h.MakeUserHandover)
		}
	}

	engine.NoRoute(api.PageNotFound)

	return engine
}

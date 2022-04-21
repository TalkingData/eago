package main

import (
	h "eago/auth/api/handler"
	"eago/auth/api/middleware"
	"eago/auth/conf"
	"eago/common/api-suite/handler"
	pg "eago/common/api-suite/pagination"
	perm "eago/common/api-suite/permission"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(conf.Conf.GinMode)

	engine := gin.Default()
	engine.Use(tracer.TracerMiddleware)

	// 登录
	engine.POST("/auth/login", middleware.ReadLoginForm(), h.IamLogin, h.DatabaseLogin, h.LoginFailed)
	engine.GET("/auth/token/content", h.GetTokenContent)

	aGroup := engine.Group("/auth", perm.MustLogin())
	{
		aGroup.POST("heartbeat", h.Heartbeat)
		aGroup.DELETE("logout", h.Logout)

		// Product模块
		pr := aGroup.Group("/products")
		{
			// 新建产品线
			pr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewProduct)
			// 删除产品线
			pr.DELETE("/:product_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveProduct)
			// 更新产品线
			pr.PUT("/:product_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetProduct)
			// 列出所有产品线
			pr.GET("", pg.ListPageHelper(), h.PagedListProducts)

			// 添加用户至产品线
			pr.POST(
				"/:product_id/users",
				perm.MustCurrUserInProductOrRole("product_id", conf.ADMIN_ROLE_NAME, true),
				h.AddUser2Product,
			)
			// 移除产品线中用户
			pr.DELETE(
				"/:product_id/users/:user_id",
				perm.MustCurrUserInProductOrRole("product_id", conf.ADMIN_ROLE_NAME, true),
				h.RemoveProductUser,
			)
			// 设置用户是否是产品线Owner
			pr.PUT("/:product_id/users/:user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetUserIsProductOwner)
			// 列出产品线中所有用户
			pr.GET("/:product_id/users", h.ListProductUsers)
		}

		// Department模块
		dr := aGroup.Group("/departments")
		{
			// 新建部门
			dr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewDepartment)
			// 删除部门
			dr.DELETE("/:department_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveDepartment)
			// 更新部门
			dr.PUT("/:department_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetDepartment)
			// 列出所有部门
			dr.GET("", pg.ListPageHelper(), h.PagedListDepartments)
			// 列出指定部门子树
			dr.GET("/:department_id/tree", h.ListDepartmentTree)
			// 以树结构列出所有部门
			dr.GET("/tree", h.ListDepartmentsTree)

			// 添加部门至角色
			dr.POST("/:department_id/users", perm.MustRole(conf.ADMIN_ROLE_NAME), h.AddUser2Department)
			// 移除部门中用户
			dr.DELETE("/:department_id/users/:user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveDepartmentUser)
			// 设置用户是否是部门Owner
			dr.PUT("/:department_id/users/:user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetUserIsDepartmentOwner)
			// 列出部门中所有用户
			dr.GET("/:department_id/users", h.ListDepartmentUsers)
		}

		// Group模块
		gr := aGroup.Group("/groups")
		{
			// 新建组
			gr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewGroup)
			// 删除组
			gr.DELETE("/:group_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveGroup)
			// 更新组
			gr.PUT("/:group_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetGroup)
			// 列出所有组
			gr.GET("", pg.ListPageHelper(), h.PagedListGroups)

			// 添加用户至组
			gr.POST(
				"/:group_id/users",
				perm.MustCurrUserInGroupOrRole("group_id", conf.ADMIN_ROLE_NAME, true),
				h.AddUser2Group,
			)
			// 移除组中用户
			gr.DELETE(
				"/:group_id/users/:user_id",
				perm.MustCurrUserInGroupOrRole("group_id", conf.ADMIN_ROLE_NAME, true),
				h.RemoveGroupUser,
			)
			// 设置用户是否是组Owner
			gr.PUT("/:group_id/users/:user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetUserIsGroupOwner)
			// 列出组中所有用户
			gr.GET("/:group_id/users", h.ListGroupUsers)
		}

		// Role模块
		rr := aGroup.Group("/roles")
		{
			// 新建角色
			rr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewRole)
			// 删除角色
			rr.DELETE("/:role_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveRole)
			// 更新角色
			rr.PUT("/:role_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetRole)
			// 列出所有角色
			rr.GET("", pg.ListPageHelper(), h.PagedListRoles)

			// 添加用户至角色
			rr.POST("/:role_id/users", perm.MustRole(conf.ADMIN_ROLE_NAME), h.AddUser2Role)
			// 移除角色中用户
			rr.DELETE("/:role_id/users/:user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveRoleUser)
			// 列出角色所有用户
			rr.GET("/:role_id/users", h.ListRoleUsers)
		}

		// User模块
		ur := aGroup.Group("/users")
		{
			// 更新用户
			ur.PUT("/:user_id", perm.MustCurrUserOrRole("user_id", conf.ADMIN_ROLE_NAME), h.SetUser)
			// 列出所有用户
			ur.GET("", pg.ListPageHelper(), h.PagedListUsers)

			// 列出用户所有角色
			ur.GET("/:user_id/roles", perm.MustCurrUserOrRole("user_id", conf.ADMIN_ROLE_NAME), h.ListUserRoles)

			// 列出用户所有产品线
			ur.GET("/:user_id/products", perm.MustCurrUserOrRole("user_id", conf.ADMIN_ROLE_NAME), h.ListUserProducts)

			// 列出用户所有组
			ur.GET("/:user_id/groups", perm.MustCurrUserOrRole("user_id", conf.ADMIN_ROLE_NAME), h.ListUserGroups)

			// 列出用户所有部门
			ur.GET("/:user_id/department", perm.MustCurrUserOrRole("user_id", conf.ADMIN_ROLE_NAME), h.GetUserDepartment)

			// 用户交接
			ur.GET("/handover/:user_id/:target_user_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.MakeUserHandover)

		}

	}

	engine.NoRoute(handler.PageNoFound)

	return engine
}

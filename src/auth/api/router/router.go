package router

import (
	h "eago-auth/api/handler"
	m "eago-auth/api/middleware"
	"eago-auth/conf"
	"eago-auth/docs"
	"eago-common/api-suite/handler"
	pg "eago-common/api-suite/pagination"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"time"
)

var Engine *gin.Engine

// InitEngine
func InitEngine() {
	Engine = gin.Default()

	Engine.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Token",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	// Swagger文档
	docs.SwaggerInfo.Title = "Eago auth API document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/auth"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger Api docs
	Engine.GET("/auth/apidoc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 登录
	Engine.POST("/auth/login", m.ReadLoginForm(), h.IamLogin, h.DatabaseLogin, h.LoginFailed)

	aGroup := Engine.Group("/auth", m.MustLogin())
	{
		aGroup.POST("heartbeat", h.Heartbeat)
		aGroup.DELETE("logout", h.Logout)
		aGroup.GET("/token/content", h.GetTokenContent)

		// User模块
		ur := aGroup.Group("/users")
		{
			// 更新用户
			ur.PUT("/:user_id", m.MustCurrUserOrRole("user_id", conf.Config.AdminRoleName), h.SetUser)
			// 列出所有用户
			ur.GET("", pg.ListPageHelper(), h.ListUsers)

			// 列出用户所有角色
			ur.GET("/:user_id/roles", m.MustCurrUserOrRole("user_id", conf.Config.AdminRoleName), h.ListUserRoles)

			// 列出用户所有产品线
			ur.GET("/:user_id/products", m.MustCurrUserOrRole("user_id", conf.Config.AdminRoleName), h.ListUserProducts)

			// 列出用户所有组
			ur.GET("/:user_id/groups", m.MustCurrUserOrRole("user_id", conf.Config.AdminRoleName), h.ListUserGroups)

			// 列出用户所有部门
			ur.GET("/:user_id/department", m.MustCurrUserOrRole("user_id", conf.Config.AdminRoleName), h.GetUserDepartment)
		}

		// Product模块
		pr := aGroup.Group("/products")
		{
			// 新建产品线
			pr.POST("", m.MustRole(conf.Config.AdminRoleName), h.NewProduct)
			// 删除产品线
			pr.DELETE("/:product_id", m.MustRole(conf.Config.AdminRoleName), h.RemoveProduct)
			// 更新产品线
			pr.PUT("/:product_id", m.MustRole(conf.Config.AdminRoleName), h.SetProduct)
			// 列出所有产品线
			pr.GET("", pg.ListPageHelper(), h.ListProducts)

			// 添加用户至产品线
			pr.POST(
				"/:product_id/users",
				m.MustCurrUserInProductOrRole("product_id", conf.Config.AdminRoleName, true),
				h.AddUser2Product,
			)
			// 移除产品线中用户
			pr.DELETE(
				"/:product_id/users/:user_id",
				m.MustCurrUserInProductOrRole("product_id", conf.Config.AdminRoleName, true),
				h.RemoveProductUser,
			)
			// 设置用户是否是产品线Owner
			pr.PUT("/:product_id/users/:user_id", m.MustRole(conf.Config.AdminRoleName), h.SetUserIsProductOwner)
			// 列出产品线中所有用户
			pr.GET("/:product_id/users", h.ListProductUsers)
		}

		// Department模块
		dr := aGroup.Group("/departments")
		{
			// 新建部门
			dr.POST("", m.MustRole(conf.Config.AdminRoleName), h.NewDepartment)
			// 删除部门
			dr.DELETE("/:department_id", m.MustRole(conf.Config.AdminRoleName), h.RemoveDepartment)
			// 更新部门
			dr.PUT("/:department_id", m.MustRole(conf.Config.AdminRoleName), h.SetDepartment)
			// 列出所有部门
			dr.GET("", pg.ListPageHelper(), h.ListDepartments)
			// 列出指定部门子树
			dr.GET("/:department_id/tree", h.ListDepartmentTree)
			// 以树结构列出所有部门
			// 真实URI = "/v1/auth/departments/tree"
			dr.GET("/:department_id", m.Dispatch("/v1/auth/departments/tree"), h.ListDepartmentsTree)

			// 添加部门至角色
			dr.POST("/:department_id/users", m.MustRole(conf.Config.AdminRoleName), h.AddUser2Department)
			// 移除部门中用户
			dr.DELETE("/:department_id/users/:user_id", m.MustRole(conf.Config.AdminRoleName), h.RemoveDepartmentUser)
			// 设置用户是否是部门Owner
			dr.PUT("/:department_id/users/:user_id", m.MustRole(conf.Config.AdminRoleName), h.SetUserIsDepartmentOwner)
			// 列出部门中所有用户
			dr.GET("/:department_id/users", h.ListDepartmentUsers)
		}

		// Group模块
		gr := aGroup.Group("/groups")
		{
			// 新建组
			gr.POST("", m.MustRole(conf.Config.AdminRoleName), h.NewGroup)
			// 删除组
			gr.DELETE("/:group_id", m.MustRole(conf.Config.AdminRoleName), h.RemoveGroup)
			// 更新组
			gr.PUT("/:group_id", m.MustRole(conf.Config.AdminRoleName), h.SetGroup)
			// 列出所有组
			gr.GET("", pg.ListPageHelper(), h.ListGroups)

			// 添加用户至组
			gr.POST(
				"/:group_id/users",
				m.MustCurrUserInGroupOrRole("group_id", conf.Config.AdminRoleName, true),
				h.AddUser2Group,
			)
			// 移除组中用户
			gr.DELETE(
				"/:group_id/users/:user_id",
				m.MustCurrUserInGroupOrRole("group_id", conf.Config.AdminRoleName, true),
				h.RemoveGroupUser,
			)
			// 设置用户是否是组Owner
			gr.PUT("/:group_id/users/:user_id", m.MustRole(conf.Config.AdminRoleName), h.SetUserIsGroupOwner)
			// 列出组中所有用户
			gr.GET("/:group_id/users", h.ListGroupUsers)
		}

		// Role模块
		rr := aGroup.Group("/roles", m.MustRole(conf.Config.AdminRoleName))
		{
			// 新建角色
			rr.POST("", h.NewRole)
			// 删除角色
			rr.DELETE("/:role_id", h.RemoveRole)
			// 更新角色
			rr.PUT("/:role_id", h.SetRole)
			// 列出所有角色
			rr.GET("", pg.ListPageHelper(), h.ListRoles)

			// 添加用户至角色
			rr.POST("/:role_id/users", h.AddUser2Role)
			// 移除角色中用户
			rr.DELETE("/:role_id/users/:user_id", h.RemoveRoleUser)
			// 列出角色所有用户
			rr.GET("/:role_id/users", h.ListRoleUsers)
		}

	}

	Engine.NoRoute(handler.PageNoFound)
}

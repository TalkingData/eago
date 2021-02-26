package router

import (
	h "eago-auth/api/handler"
	m "eago-auth/api/middleware"
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

func InitEngine() error {
	Engine = gin.Default()

	Engine.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	// Swagger文档
	docs.SwaggerInfo.Title = "Eago auth API document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1/auth"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger Api docs
	Engine.GET("/v1/auth/apidocs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 登录
	Engine.POST("/v1/auth/login", m.ReadLoginForm(), h.CrowdLogin, h.DatabaseLogin, h.LoginFailed)

	v1 := Engine.Group("/v1/auth", m.MustLogin())
	{
		v1.POST("heartbeat", h.Heartbeat)
		v1.DELETE("logout", h.Logout)
		v1.GET("/token/content", h.GetTokenContent)

		// User模块
		ur := v1.Group("/users")
		{
			// 更新用户
			ur.PUT("/:user_id", m.MustCurrUserOrRole("user_id", "auth_admin"), h.SetUser)
			// 列出所有用户
			ur.GET("", m.MustRole("auth_admin"), pg.ListPageHelper(), h.ListUsers)

			// 列出用户所有角色
			ur.GET("/:user_id/roles", m.MustCurrUserOrRole("user_id", "auth_admin"), h.ListUserRoles)

			// 列出用户所有产品线
			ur.GET("/:user_id/products", m.MustCurrUserOrRole("user_id", "auth_admin"), h.ListUserProducts)

			// 列出用户所有组
			ur.GET("/:user_id/groups", m.MustCurrUserOrRole("user_id", "auth_admin"), h.ListUserGroups)

			// 列出用户所有部门
			ur.GET("/:user_id/department", m.MustCurrUserOrRole("user_id", "auth_admin"), h.GetUserDepartment)
		}

		// Product模块
		pr := v1.Group("/products", m.MustRole("auth_admin"))
		{
			// 新建产品线
			pr.POST("", h.NewProduct)
			// 删除产品线
			pr.DELETE("/:product_id", h.DeleteProduct)
			// 更新产品线
			pr.PUT("/:product_id", h.SetProduct)
			// 列出所有产品线
			pr.GET("", pg.ListPageHelper(), h.ListProducts)

			// 添加用户至产品线
			pr.POST("/:product_id/users", h.AddUser2Product)
			// 移除产品线中用户
			pr.DELETE("/:product_id/users/:user_id", h.RemoveProductUser)
			// 设置用户是否是产品线Owner
			pr.PUT("/:product_id/users/:user_id", h.SetUserIsProductOwner)
			// 列出产品线中所有用户
			pr.GET("/:product_id/users", h.ListProductUsers)
		}

		// Department模块
		dr := v1.Group("/departments", m.MustRole("auth_admin"))
		{
			// 新建部门
			dr.POST("", h.NewDepartment)
			// 删除部门
			dr.DELETE("/:department_id", h.DeleteDepartment)
			// 更新部门
			dr.PUT("/:department_id", h.SetDepartment)
			// 列出所有部门
			dr.GET("", pg.ListPageHelper(), h.ListDepartments)
			// 列出指定部门子树
			dr.GET("/:department_id/tree", h.ListDepartmentTree)
			// 以树结构列出所有部门
			// 真实URI = "/v1/auth/departments/tree"
			dr.GET("/:department_id", m.Dispatch("/v1/auth/departments/tree"), h.ListDepartmentsTree)

			// 添加部门至角色
			dr.POST("/:department_id/users", h.AddUser2Department)
			// 移除部门中用户
			dr.DELETE("/:department_id/users/:user_id", h.RemoveDepartmentUser)
			// 设置用户是否是部门Owner
			dr.PUT("/:department_id/users/:user_id", h.SetUserIsDepartmentOwner)
			// 列出部门中所有用户
			dr.GET("/:department_id/users", h.ListDepartmentUsers)
		}

		// Role模块
		rr := v1.Group("/roles", m.MustRole("auth_admin"))
		{
			// 新建角色
			rr.POST("", h.NewRole)
			// 删除角色
			rr.DELETE("/:role_id", h.DeleteRole)
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

		// Group模块
		gr := v1.Group("/groups", m.MustRole("auth_admin"))
		{
			// 新建组
			gr.POST("", h.NewGroup)
			// 删除组
			gr.DELETE("/:group_id", h.DeleteGroup)
			// 更新组
			gr.PUT("/:group_id", h.SetGroup)
			// 列出所有组
			gr.GET("", pg.ListPageHelper(), h.ListGroups)

			// 添加用户至组
			gr.POST("/:group_id/users", h.AddUser2Group)
			// 移除组中用户
			gr.DELETE("/:group_id/users/:user_id", h.RemoveGroupUser)
			// 设置用户是否是组Owner
			gr.PUT("/:group_id/users/:user_id", h.SetUserIsGroupOwner)
			// 列出组中所有用户
			gr.GET("/:group_id/users", h.ListGroupUsers)
		}
	}

	Engine.NoRoute(handler.PageNoFound)

	return nil
}

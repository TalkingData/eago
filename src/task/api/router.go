package main

import (
	"eago/common/api-suite/handler"
	pg "eago/common/api-suite/pagination"
	"eago/task/api/docs"
	h "eago/task/api/handler"
	m "eago/task/api/middleware"
	"eago/task/conf"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var Engine *gin.Engine

// init
func init() {
	Engine = gin.Default()

	// Swagger文档
	docs.SwaggerInfo.Title = "Eago task API document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/task"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger Api docs
	Engine.GET("/task/apidoc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tGroup := Engine.Group("/task", m.MustLogin())
	{
		// 列出所有Worker
		tGroup.GET("/workers", m.MustRole(conf.ADMIN_ROLE_NAME), h.ListWorkers)

		// Task模块
		tr := tGroup.Group("/tasks")
		{
			tr.POST("", m.MustRole(conf.ADMIN_ROLE_NAME), h.NewTask)
			tr.DELETE("/:task_id", m.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveTask)
			tr.PUT("/:task_id", m.MustRole(conf.ADMIN_ROLE_NAME), h.SetTask)
			tr.GET("", m.MustRole(conf.ADMIN_ROLE_NAME), pg.ListPageHelper(), h.ListTasks)

			// 调用任务
			tr.POST("/:task_id/call", m.MustRole(conf.ADMIN_ROLE_NAME), h.CallTask)
		}

		// ResultTables模块
		rpr := tGroup.Group("/result_partitions")
		{
			// 仅测试用，不对外开放的方法
			// 新建结果分区并建立结果表和日志表
			rpr.POST("/with_create_tables", m.MustRole(conf.ADMIN_ROLE_NAME), h.NewResultPartitionsWithCreateTables)

			// 列出所有结果分区
			rpr.GET("", h.ListResultPartitions)
		}

		// Result模块
		rr := tGroup.Group("/results")
		{
			// 按分区ID列出所有结果
			rr.GET("/:result_partition_id", pg.ListPageHelper(), h.ListResults)
			// 手动结束任务
			rr.DELETE("/:result_partition_id/:result_id", m.MustRole(conf.ADMIN_ROLE_NAME), h.KillTask)
		}

		// Log模块
		// 按分区ID列出所有结果日志
		tGroup.GET("/logs/:result_partition_id/:result_id", pg.ListPageHelper(), h.ListLogs)

	}
	// 以WebSocket方式按分区ID列出所有结果日志
	Engine.GET("/task/logs/:result_partition_id/:result_id/ws", h.WsListLogs)

	Engine.NoRoute(handler.PageNoFound)
}

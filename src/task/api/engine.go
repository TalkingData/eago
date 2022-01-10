package main

import (
	"eago/common/api-suite/handler"
	pg "eago/common/api-suite/pagination"
	perm "eago/common/api-suite/permission"
	"eago/common/tracer"
	h "eago/task/api/handler"
	"eago/task/conf"
	"github.com/gin-gonic/gin"
)

// NewGinEngine
func NewGinEngine() *gin.Engine {
	gin.SetMode(conf.Conf.GinMode)

	engine := gin.Default()
	engine.Use(tracer.TracerMiddleware)

	tGroup := engine.Group("/task", perm.MustLogin())
	{
		// 列出所有Worker
		tGroup.GET("/workers", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListWorkers)

		// Task模块
		tr := tGroup.Group("/tasks")
		{
			tr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewTask)
			tr.DELETE("/:task_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveTask)
			tr.PUT("/:task_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetTask)
			tr.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), pg.ListPageHelper(), h.ListTasks)

			// 调用任务
			tr.POST("/:task_id/call", perm.MustRole(conf.ADMIN_ROLE_NAME), h.CallTask)
		}

		// Schedule模块
		sr := tGroup.Group("/schedules")
		{
			sr.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewSchedule)
			sr.DELETE("/:schedule_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveSchedule)
			sr.PUT("/:schedule_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetSchedule)
			sr.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), pg.ListPageHelper(), h.ListSchedules)
		}

		// ResultTables模块
		rpr := tGroup.Group("/result_partitions")
		{
			// 仅测试用，不对外开放的方法
			// 新建结果分区并建立结果表和日志表
			rpr.POST("/with_create_tables", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewResultPartitionsWithCreateTables)

			// 列出所有结果分区
			rpr.GET("", h.ListResultPartitions)
		}

		// Result模块
		rr := tGroup.Group("/results")
		{
			// 按分区ID列出所有结果
			rr.GET("/:result_partition_id", pg.ListPageHelper(), h.ListResults)
			// 手动结束任务
			rr.DELETE("/:result_partition_id/:result_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.KillTask)
		}

		// Log模块
		// 按分区ID列出所有结果日志
		tGroup.GET("/logs/:result_partition_id/:result_id", h.ListLogs)
	}

	// Log模块
	// 以WebSocket方式按分区ID列出所有结果日志
	engine.GET("/task/logs/:result_partition_id/:result_id/ws", perm.MustLoginWs(), h.WsListLogs)

	engine.NoRoute(handler.PageNoFound)

	return engine
}

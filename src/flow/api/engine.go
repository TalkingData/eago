package main

import (
	"eago/common/api-suite/handler"
	pg "eago/common/api-suite/pagination"
	perm "eago/common/api-suite/permission"
	"eago/common/tracer"
	h "eago/flow/api/handler"
	"eago/flow/conf"
	"github.com/gin-gonic/gin"
)

// NewGinEngine
func NewGinEngine() *gin.Engine {
	gin.SetMode(conf.Conf.GinMode)

	engine := gin.Default()
	engine.Use(tracer.TracerMiddleware)

	fGroup := engine.Group("/flow", perm.MustLogin())
	{
		// Instance流程实例模块
		iR := fGroup.Group("/instances")
		{
			// 处理指定流程实例
			iR.PUT("/:instance_id/handle", h.HandleInstance)

			// 列出我发起的流程实例
			iR.GET("/my", pg.ListPageHelper(), h.ListMyInstances)
			// 列出我代办的流程实例
			iR.GET("/todo", pg.ListPageHelper(), h.ListTodoInstances)
			// 列出我已办的流程实例
			iR.GET("/done", pg.ListPageHelper(), h.ListDoneInstances)

			// 列出所有流程实例，要求管理员权限
			iR.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), pg.ListPageHelper(), h.ListInstances)
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
			cR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewCategory)
			// 删除类别，要求管理员权限
			cR.DELETE("/:category_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveCategory)
			// 更新类别，要求管理员权限
			cR.PUT("/:category_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetCategory)
		}

		// Flow流程模块
		flR := fGroup.Group("/flows")
		{
			// 发起流程
			flR.POST("/:flow_id/instantiate", h.InstantiateFlow)

			// 分页列出所有流程
			flR.GET("", pg.ListPageHelper(), h.ListFlows)

			// 新建流程，要求管理员权限
			flR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewFlow)
			// 删除流程，要求管理员权限
			flR.DELETE("/:flow_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveFlow)
			// 更新流程，要求管理员权限
			flR.PUT("/:flow_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetFlow)
		}

		// Form表单模块
		fR := fGroup.Group("/forms")
		{
			// 获取指定表单
			fR.GET("/:form_id", h.GetForm)

			// 新建表单，要求管理员权限
			fR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewForm)
			// 变更表单，要求管理员权限
			fR.PUT("/:form_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetForm)
			// 列出所有表单，要求管理员权限
			fR.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), pg.ListPageHelper(), h.ListForms)

			// 列出指定表单所关联流程，要求管理员权限
			fR.GET("/:form_id/flows", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListFormFlows)
		}

		// Node节点模块
		nR := fGroup.Group("/nodes")
		{
			// 列出指定节点链
			nR.GET("/:node_id/chain", h.GetNodeChain)

			// 新增节点，要求管理员权限
			nR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewNode)
			// 删除节点，要求管理员权限
			nR.DELETE("/:node_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveNode)
			// 更新节点，要求管理员权限
			nR.PUT("/:node_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetNode)
			// 列出所有节点，要求管理员权限
			nR.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListNodes)

			// 添加触发器至节点，要求管理员权限
			nR.POST("/:node_id/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME), h.AddTrigger2Node)
			// 移除指定节点中触发器，要求管理员权限
			nR.DELETE("/:node_id/triggers/:trigger_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveNodeTrigger)
			// 列出指定节点中所有触发器，要求管理员权限
			nR.GET("/:node_id/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListNodeTriggers)

			// 列出指定节点所关联流程，要求管理员权限
			nR.GET("/:node_id/flows", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListNodeFlows)
		}

		// Trigger触发器模块
		tR := fGroup.Group("/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME))
		{
			// 新增触发器，要求管理员权限
			tR.POST("", h.NewTrigger)
			// 删除触发器，要求管理员权限
			tR.DELETE("/:trigger_id", h.RemoveTrigger)
			// 变更触发器，要求管理员权限
			tR.PUT("/:trigger_id", h.SetTrigger)
			// 列出所有触发器，要求管理员权限
			tR.GET("", pg.ListPageHelper(), h.ListTriggers)

			// 列出指定触发器所关联节点，要求管理员权限
			tR.GET("/:trigger_id/nodes", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListTriggerNodes)
		}
	}

	engine.NoRoute(handler.PageNoFound)

	return engine
}

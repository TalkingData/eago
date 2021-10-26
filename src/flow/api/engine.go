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
			iR.GET("", pg.ListPageHelper(), h.ListInstances)
			iR.PUT("/:instance_id/handle", h.HandleInstance)
		}

		// Flow流程模块
		flR := fGroup.Group("/flows")
		{
			// 新建流程
			flR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewFlow)
			flR.DELETE("/:flow_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveFlow)
			flR.PUT("/:flow_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetFlow)
			flR.GET("", pg.ListPageHelper(), h.ListFlows)

			// 发起流程
			flR.POST("/:flow_id/instantiate", h.InstantiateFlow)
		}

		// Form表单模块
		fR := fGroup.Group("/forms", perm.MustRole(conf.ADMIN_ROLE_NAME))
		{
			// 新建表单
			fR.POST("", h.NewForm)
			fR.PUT("/:form_id", h.SetForm)
			fR.GET("", pg.ListPageHelper(), h.ListForms)
		}

		// Trigger触发器模块
		tR := fGroup.Group("/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME))
		{
			tR.POST("", h.NewTrigger)
			tR.DELETE("/:trigger_id", h.RemoveTrigger)
			tR.PUT("/:trigger_id", h.SetTrigger)
			tR.GET("", pg.ListPageHelper(), h.ListTriggers)
		}

		nR := fGroup.Group("/nodes")
		{
			nR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewNode)
			nR.DELETE("/:node_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveNode)
			nR.PUT("/:node_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetNode)
			// 列出指定节点链
			nR.GET("/:node_id/chain", h.GetNodeChain)
			nR.GET("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListNodes)

			// 添加触发器至节点
			nR.POST("/:node_id/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME), h.AddTrigger2Node)
			// 移除节点中触发器
			nR.DELETE("/:node_id/triggers/:trigger_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveNodeTrigger)
			// 列出节点中所有触发器
			nR.GET("/:node_id/triggers", perm.MustRole(conf.ADMIN_ROLE_NAME), h.ListNodeTriggers)

		}

		lR := fGroup.Group("/logs")
		{
			lR.POST("/:instance_id", h.NewLog)
			lR.GET("/:instance_id", h.ListLogs)
		}

		// Category 类别模块
		cR := fGroup.Group("/categories")
		{
			cR.POST("", perm.MustRole(conf.ADMIN_ROLE_NAME), h.NewCategory)
			cR.DELETE("/:category_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.RemoveCategory)
			cR.PUT("/:category_id", perm.MustRole(conf.ADMIN_ROLE_NAME), h.SetCategory)
			cR.GET("", h.ListCategories)
		}
	}

	engine.NoRoute(handler.PageNoFound)

	return engine
}

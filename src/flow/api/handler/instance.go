package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"eago/flow/srv/builtin"
	"github.com/gin-gonic/gin"
	"strconv"
)

// HandleInstance 处理流程实例
func HandleInstance(c *gin.Context) {
	insId, err := strconv.Atoi(c.Param("instance_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "instance_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var hdIns dto.HandleInstance
	// 序列化request body
	if err := c.ShouldBindJSON(&hdIns); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")
	// 验证数据
	if m := hdIns.Validate(insId, tc["Username"].(string)); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 执行实际流程实例处理
	if err = builtin.HandleInstance(hdIns.Instance, &hdIns); err != nil {
		m := msg.UndefinedError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "instance_id", insId)
}

// PagedListInstances 列出所有流程实例-分页
func PagedListInstances(c *gin.Context) {
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")

	query := dao.Query{}
	// 设置查询filter
	liq := dto.PagedListInstancesQuery{}
	if c.ShouldBindQuery(&liq) == nil {
		_ = liq.UpdateDefaultQuery(query, tc["Username"].(string))
	}
	pagedListInstances(c, query)
}

// PagedListMyInstances 列出我发起的流程实例-分页
func PagedListMyInstances(c *gin.Context) {
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")

	query := dao.Query{}
	// 设置查询filter
	liq := dto.PagedListInstancesQuery{}
	if c.ShouldBindQuery(&liq) == nil {
		_ = liq.UpdateMyInstancesQuery(query, tc["Username"].(string))
	}
	pagedListInstances(c, query)
}

// PagedListTodoInstances 列出我代办的流程实例-分页
func PagedListTodoInstances(c *gin.Context) {
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")

	query := dao.Query{}
	// 设置查询filter
	liq := dto.PagedListInstancesQuery{}
	if c.ShouldBindQuery(&liq) == nil {
		_ = liq.UpdateTodoInstancesQuery(query, tc["Username"].(string))
	}
	pagedListInstances(c, query)
}

// PagedListDoneInstances 列出我已办的流程实例-分页
func PagedListDoneInstances(c *gin.Context) {
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")

	query := dao.Query{}
	// 设置查询filter
	liq := dto.PagedListInstancesQuery{}
	if c.ShouldBindQuery(&liq) == nil {
		_ = liq.UpdateDoneInstancesQuery(query, tc["Username"].(string))
	}
	pagedListInstances(c, query)
}

// pagedListInstances 列出所有流程实例-分页
func pagedListInstances(c *gin.Context, query dao.Query) {
	paged, ok := dao.PagedListInstances(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UndefinedError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "instances", paged)
}

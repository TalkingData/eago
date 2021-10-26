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
		m := msg.UnknownError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "instance_id", insId)
}

// ListInstances 列出所有流程实例
func ListInstances(c *gin.Context) {
	// 取出TokenContent
	tc := c.GetStringMap("TokenContent")

	query := dao.Query{}
	// 设置查询filter
	liq := dto.ListInstancesQuery{}
	if c.ShouldBindQuery(&liq) == nil {
		_ = liq.UpdateQuery(query, tc["Username"].(string))
	}

	paged, ok := dao.PagedListInstances(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UnknownError
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "instances", paged)
}

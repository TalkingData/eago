package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewTrigger 新建触发器
func NewTrigger(c *gin.Context) {
	var tFrm dto.NewTrigger

	// 序列化request body
	if err := c.ShouldBindJSON(&tFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := tFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	t := dao.NewTrigger(tFrm.Name, *tFrm.Description, tFrm.TaskCodename, tFrm.Arguments, tc["Username"].(string))
	if t == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "trigger", t)
}

// RemoveTrigger 删除触发器
func RemoveTrigger(c *gin.Context) {
	tId, err := strconv.Atoi(c.Param("trigger_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "trigger_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rtFrm dto.RemoveTrigger
	// 验证数据
	if m := rtFrm.Validate(tId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveTrigger(tId); !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetTrigger 更新触发器
func SetTrigger(c *gin.Context) {
	tId, err := strconv.Atoi(c.Param("trigger_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "trigger_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var stFrm dto.SetTrigger
	// 序列化request body
	if err = c.ShouldBindJSON(&stFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := stFrm.Validate(tId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	t, ok := dao.SetTrigger(tId, stFrm.Name, *stFrm.Description, stFrm.TaskCodename, stFrm.Arguments, tc["Username"].(string))
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "trigger", t)
}

// ListTriggers 列出所有触发器
func ListTriggers(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	ltq := dto.ListTriggersQuery{}
	if c.ShouldBindQuery(&ltq) == nil {
		_ = ltq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListTriggers(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "triggers", paged)
}

// ListTriggerNodes 列出触发器所关联节点
func ListTriggerNodes(c *gin.Context) {
	tId, err := strconv.Atoi(c.Param("trigger_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "trigger_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var ltnFrm dto.ListTriggerNodes
	// 验证数据
	if m := ltnFrm.Validate(tId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	u, ok := dao.ListTriggerNodes(tId)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "nodes", u)
}

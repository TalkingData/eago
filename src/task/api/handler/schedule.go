package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewSchedule 新建计划任务
func NewSchedule(c *gin.Context) {
	var schFrm dto.NewSchedule

	// 序列化request body
	if err := c.ShouldBindJSON(&schFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := schFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	s := dao.NewSchedule(
		schFrm.TaskCodename,
		schFrm.Expression,
		schFrm.Arguments,
		*schFrm.Description,
		*schFrm.Timeout,
		*schFrm.Disabled,
		tc["Username"].(string),
	)
	// 新建失败
	if s == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "schedule", s)
}

// RemoveSchedule 删除计划任务
func RemoveSchedule(c *gin.Context) {
	schId, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "schedule_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var tFrm dto.RemoveSchedule
	if m := tFrm.Validate(schId); m != nil {
		// 数据验证未通过
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveSchedule(schId); !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetSchedule 更新计划任务
func SetSchedule(c *gin.Context) {
	schId, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "schedule_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var schFrm dto.SetSchedule
	// 序列化request body
	if err := c.ShouldBindJSON(&schFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := schFrm.Validate(schId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	sch, ok := dao.SetSchedule(
		schId,
		schFrm.TaskCodename,
		schFrm.Expression,
		schFrm.Arguments,
		*schFrm.Description,
		*schFrm.Timeout,
		*schFrm.Disabled,
		tc["Username"].(string),
	)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "schedule", sch)
}

// ListSchedules 列出所有计划任务
func ListSchedules(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	ltq := dto.ListTasksQuery{}
	if c.ShouldBindQuery(&ltq) == nil {
		_ = ltq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListSchedules(
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

	w.WriteSuccessPayload(c, "schedules", paged)
}

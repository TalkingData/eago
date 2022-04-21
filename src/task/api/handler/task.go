package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/dto"
	"eago/task/srv/builtin"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CallTask 调用任务
func CallTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "task_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var ctFrm dto.CallTask
	// 序列化request body
	if err = c.ShouldBindJSON(&ctFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := ctFrm.Validate(taskId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 调用任务
	tc := c.GetStringMap("TokenContent")
	tUid, err := builtin.CallTask(ctFrm.TaskCodeName, ctFrm.Arguments, tc["Username"].(string), *ctFrm.Timeout)
	// 调用失败
	if err != nil {
		m := msg.CallTaskFailed
		log.ErrorWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "task_unique_id", tUid)
}

// NewTask 新建任务
func NewTask(c *gin.Context) {
	var ntFrm dto.NewTask

	// 序列化request body
	if err := c.ShouldBindJSON(&ntFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := ntFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	t := dao.NewTask(
		*ntFrm.Disabled,
		*ntFrm.Category,
		ntFrm.Codename,
		*ntFrm.Description,
		ntFrm.FormalParams,
		tc["Username"].(string),
	)
	// 新建失败
	if t == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "task", t)
}

// RemoveTask 删除任务
func RemoveTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "task_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var tFrm dto.RemoveTask
	if m := tFrm.Validate(taskId); m != nil {
		// 数据验证未通过
		m.WriteRest(c)
		return
	}

	// 删除
	if ok := dao.RemoveTask(taskId); !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetTask 更新任务
func SetTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "task_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var taskFrm dto.SetTask
	// 序列化request body
	if err := c.ShouldBindJSON(&taskFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := taskFrm.Validate(taskId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 更新
	tc := c.GetStringMap("TokenContent")
	task, ok := dao.SetTask(
		taskId,
		*taskFrm.Disabled,
		*taskFrm.Category,
		taskFrm.Codename,
		*taskFrm.Description,
		taskFrm.FormalParams,
		tc["Username"].(string),
	)
	// 更新失败
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "task", task)
}

// PagedListTasks 列出所有任务-分页
func PagedListTasks(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	ltq := dto.ListTasksQuery{}
	if c.ShouldBindQuery(&ltq) == nil {
		_ = ltq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListTasks(
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

	w.WriteSuccessPayload(c, "tasks", paged)
}

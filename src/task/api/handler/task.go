package handler

import (
	"database/sql"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/model"
	"eago/task/srv/local"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CallTask 调用任务
// @Summary 调用任务
// @Tags 任务
// @Param token header string true "Token"
// @Param data body model.CallTask true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","task":{"id":4,"category":3,"codename":"auth.test_task14","description":"desc","arguments":"{}","disabled":false,"created_at":"2021-03-03 16:51:42","created_by":"test","updated_at":"2021-03-03 16:51:42","updated_by":""}}"
// @Router /tasks/{task_id}/call [POST]
func CallTask(c *gin.Context) {
	var callTask model.CallTask

	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'task_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&callTask); err != nil {
		resp := msg.WarnInvalidBody.GenResponse(err.Error())
		log.Warn(resp.String())
		resp.Write(c)
		return
	}

	tObj, ok := model.GetTask(model.Query{"id=?": taskId, "disabled=?": false})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.GetTask.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}
	if tObj == nil {
		resp := msg.WarnNotFound.GenResponse("Keep task exist and not disabled.")
		log.Warn(resp.String())
		resp.Write(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	tUid, err := local.CallTask(tObj.Codename, callTask.Arguments, tc["Username"].(string), *callTask.Timeout)
	if err != nil {
		resp := msg.ErrCallTask.GenResponse(err.Error())
		log.ErrorWithFields(log.Fields{
			"id":    tObj.Id,
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("task_unique_id", tUid)
	resp.Write(c)
}

// NewTask 新建任务
// @Summary 新建任务
// @Tags 任务
// @Param token header string true "Token"
// @Param data body model.Task true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","task":{"id":4,"category":3,"codename":"auth.test_task14","description":"desc","arguments":"{}","disabled":false,"created_at":"2021-03-03 16:51:42","created_by":"test","updated_at":"2021-03-03 16:51:42","updated_by":""}}"
// @Router /tasks [POST]
func NewTask(c *gin.Context) {
	var task model.Task

	// 序列化request body
	if err := c.ShouldBindJSON(&task); err != nil {
		resp := msg.WarnInvalidBody.GenResponse(err.Error())
		log.Warn(resp.String())
		resp.Write(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	t := model.NewTask(task.Disabled, task.Category, task.Codename, *task.Description, task.Arguments, tc["Username"].(string))
	if t == nil {
		resp := msg.ErrDatabase.GenResponse("Error in model.NewTask.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("task", t)
	resp.Write(c)
}

// RemoveTask 删除任务
// @Summary 删除任务
// @Tags 任务
// @Param token header string true "Token"
// @Param task_id path string true "任务ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /tasks/{task_id} [DELETE]
func RemoveTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'task_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if ok := model.RemoveTask(taskId); !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.model.RemoveTask.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	msg.Success.GenResponse().Write(c)
}

// SetTask 更新任务
// @Summary 更新任务
// @Tags 任务
// @Param token header string true "Token"
// @Param task_id path string true "任务ID"
// @Param data body model.Task true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","task":{"id":4,"category":3,"codename":"auth.test_task14","description":"desc","arguments":"{}","disabled":false,"created_at":"2021-03-03 16:51:42","created_by":"test","updated_at":"2021-03-03 16:51:42","updated_by":""}}"
// @Router /tasks/{task_id} [PUT]
func SetTask(c *gin.Context) {
	var taskFrm model.Task

	taskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'task_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&taskFrm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'category', 'codename', 'description', 'arguments' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	task, ok := model.SetTask(taskId, taskFrm.Disabled, taskFrm.Category, taskFrm.Codename, *taskFrm.Description, taskFrm.Arguments, tc["Username"].(string))
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.SetTask.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("task", task)
	resp.Write(c)
}

// ListTasks 列出所有任务
// @Summary 列出所有任务
// @Tags 任务
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"tasks":[{"id":1,"category":0,"codename":"auth.sync_department","description":"同步部门","arguments":"{}","disabled":false,"created_at":"2021-02-23 07:24:14","created_by":"","updated_at":null,"updated_by":""}],"total":1}"
// @Router /tasks [GET]
func ListTasks(c *gin.Context) {
	var query model.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"codename LIKE @query OR description LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListTasks(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.PageListTasks.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "tasks")
	resp.Write(c)
}

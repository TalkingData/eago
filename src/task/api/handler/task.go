package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/tracer"
	"eago/task/api/form"
	"eago/task/conf/msg"
	taskpb "eago/task/proto"
	"github.com/gin-gonic/gin"
)

// CallTask 调用任务
func (th *TaskHandler) CallTask(c *gin.Context) {
	taskId, err := ext.ParamUint32(c, "task_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "task_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.CallTask{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, th.dao, taskId); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 调用任务
	req := taskpb.CallTaskReq{
		TaskCodename: frm.TaskCodeName,
		Timeout:      *frm.Timeout,
		Arguments:    []byte(frm.Arguments),
		Caller:       perm.MustGetTokenContent(c).Username,
	}
	rsp, err := th.taskCli.CallTask(ctx, &req)
	// 调用失败
	if err != nil {
		m := msg.MsgCallTaskFailed.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "task_unique_id", rsp.TaskUniqueId)
}

// NewTask 新建任务
func (th *TaskHandler) NewTask(c *gin.Context) {
	frm := form.NewTaskForm{}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 序列化request body
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, th.dao); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 新建
	task, err := th.dao.NewTask(
		ctx,
		*frm.Disabled,
		*frm.Category,
		frm.Codename,
		*frm.Description,
		frm.FormalParams,
		perm.MustGetTokenContent(c).Username,
	)
	// 新建失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "task", task)
}

// RemoveTask 删除任务
func (th *TaskHandler) RemoveTask(c *gin.Context) {
	taskId, err := ext.ParamUint32(c, "task_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "task_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveTaskForm{}
	if m := frm.Validate(ctx, th.dao, taskId); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 删除
	if err = th.dao.RemoveTask(ctx, taskId); err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetTask 更新任务
func (th *TaskHandler) SetTask(c *gin.Context) {
	taskId, err := ext.ParamUint32(c, "task_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "task_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.SetTaskForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, th.dao, taskId); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 更新
	task, err := th.dao.SetTask(
		ctx,
		taskId,
		*frm.Disabled,
		*frm.Category,
		frm.Codename,
		*frm.Description,
		frm.FormalParams,
		perm.MustGetTokenContent(c).Username,
	)
	// 更新失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "task", task)
}

// PagedListTasks 列出所有任务-分页
func (th *TaskHandler) PagedListTasks(c *gin.Context) {
	pFrm := form.PagedListTasksParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	paged, err := th.dao.PagedListTasks(
		tracer.ExtractTraceCtxFromGin(c),
		pFrm.GenQuery(),
		c.GetInt(global.GinCtxPageKey),
		c.GetInt(global.GinCtxPageSizeKey),
		c.GetStringSlice(global.GinCtxOrderByKey)...,
	)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "tasks", paged)
}

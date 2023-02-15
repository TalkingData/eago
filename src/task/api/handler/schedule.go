package handler

import (
	"eago/common/api/ext"
	perm "eago/common/api/permission"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/tracer"
	"eago/task/api/form"
	"eago/task/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewSchedule 新建计划任务
func (th *TaskHandler) NewSchedule(c *gin.Context) {
	frm := form.NewScheduleForm{}

	// 序列化request body
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 新建
	schObj, err := th.dao.NewSchedule(
		tracer.ExtractTraceCtxFromGin(c),
		frm.TaskCodename,
		frm.Expression,
		frm.Arguments,
		*frm.Description,
		*frm.Timeout,
		*frm.Disabled,
		perm.MustGetTokenContent(c).Username,
	)
	// 新建失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "schedule", schObj)
}

// RemoveSchedule 删除计划任务
func (th *TaskHandler) RemoveSchedule(c *gin.Context) {
	schId, err := ext.ParamUint32(c, "schedule_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "schedule_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.RemoveScheduleForm{}
	if m := frm.Validate(ctx, th.dao, schId); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if err = th.dao.RemoveSchedule(ctx, schId); err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// SetSchedule 更新计划任务
func (th *TaskHandler) SetSchedule(c *gin.Context) {
	schId, err := ext.ParamUint32(c, "schedule_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "schedule_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	frm := form.SetScheduleForm{}
	// 序列化request body
	if err = c.ShouldBindJSON(&frm); err != nil {
		m := cMsg.MsgSerializeFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 验证数据
	if m := frm.Validate(ctx, th.dao, schId); m != nil {
		// 数据验证未通过
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	sch, err := th.dao.SetSchedule(
		ctx,
		schId,
		frm.TaskCodename,
		frm.Expression,
		frm.Arguments,
		*frm.Description,
		*frm.Timeout,
		*frm.Disabled,
		perm.MustGetTokenContent(c).Username,
	)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "schedule", sch)
}

// PagedListSchedules 列出所有计划任务-分页
func (th *TaskHandler) PagedListSchedules(c *gin.Context) {
	// 设置查询filter
	pFrm := form.PagedListTasksParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		th.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in TaskHandler.PagedListSchedules, skipped it.")
	}

	paged, err := th.dao.PagedListSchedules(
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

	ext.WriteSuccessPayload(c, "schedules", paged)
}

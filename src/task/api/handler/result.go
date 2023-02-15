package handler

import (
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/tracer"
	"eago/task/api/form"
	"eago/task/conf/msg"
	"github.com/gin-gonic/gin"
)

// KillTask 手动结束任务
func (th *TaskHandler) KillTask(c *gin.Context) {
	// 获得分区ID
	resPartId, err := ext.ParamUint32(c, "result_partition_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_partition_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 查找分区
	part, err := th.dao.GetResultPartitionsPartition(ctx, resPartId)
	if err != nil {
		m := msg.MsgKillTaskPartitionNotFoundFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields().Append("result_partition_id", resPartId), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 找不到分区
	if len(part) < 1 {
		m := msg.MsgKillTaskPartitionNotFoundFailed
		th.logger.WarnWithFields(logger.Fields{
			"result_partition_id": resPartId,
		}, m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 获得结果ID
	rId, err := ext.ParamUint32(c, "result_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 执行手动结束任务
	if err = th.biz.KillTask(ctx, th.biz.TaskUniqueIdEncode(part, rId)); err != nil {
		m := msg.MsgKillTaskFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// KillTaskByTaskUniqueId 按任务唯一ID手动结束任务
func (th *TaskHandler) KillTaskByTaskUniqueId(c *gin.Context) {
	// 执行手动结束任务
	if err := th.biz.KillTask(tracer.ExtractTraceCtxFromGin(c), c.Param("task_unique_id")); err != nil {
		m := msg.MsgKillTaskFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccess(c)
}

// PagedListResults 按分区列出所有结果-分页
func (th *TaskHandler) PagedListResults(c *gin.Context) {
	// 获得分区ID
	resPartId, err := ext.ParamUint32(c, "result_partition_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_partition_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ctx := tracer.ExtractTraceCtxFromGin(c)

	// 根据分区获取结果表前缀
	part, err := th.dao.GetResultPartitionsPartition(ctx, resPartId)
	if err != nil {
		m := msg.MsgKillTaskPartitionNotFoundFailed.SetError(err)
		th.logger.WarnWithFields(m.ToLoggerFields().Append("result_partition_id", resPartId), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}
	// 找不到分区
	if len(part) < 1 {
		m := msg.MsgKillTaskPartitionNotFoundFailed
		th.logger.WarnWithFields(logger.Fields{
			"result_partition_id": resPartId,
		}, m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 设置查询filter
	pFrm := form.ListResultsParamsForm{}
	if err = c.ShouldBindQuery(&pFrm); err != nil {
		th.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in TaskHandler.PagedListResults, skipped it.")
	}

	// 查询结果
	paged, err := th.dao.PagedListResultsByPartition(
		ctx,
		pFrm.GenQuery(),
		part,
		c.GetInt(global.GinCtxPageKey),
		c.GetInt(global.GinCtxPageSizeKey),
		c.GetStringSlice(global.GinCtxOrderByKey)...,
	)
	// 查询失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "results", paged)
}

// GetResultByTaskUniqueId 按任务唯一ID查询结果
func (th *TaskHandler) GetResultByTaskUniqueId(c *gin.Context) {
	// 获得任务唯一ID
	taskUniqueId := c.Param("task_unique_id")
	// 解码任务唯一ID
	part, resId, err := th.biz.TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "task_unique_id")
		th.logger.WarnWithFields(m.ToLoggerFields().Append("task_unique_id", taskUniqueId), m.GetMsg())
		m.Write2GinCtx(c)
	}

	// 查询结果
	result, err := th.dao.GetResult(tracer.ExtractTraceCtxFromGin(c), part, resId)
	// 查询失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "result", result)
}

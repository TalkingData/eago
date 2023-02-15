package handler

import (
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/logger"
	"eago/common/tracer"
	"eago/task/api/form"
	"eago/task/conf/msg"
	"github.com/gin-gonic/gin"
)

// NewResultPartitionsWithCreateTables 新增结果分区并建立结果表和日志表
// 仅测试用，不对外开放的方法
func (th *TaskHandler) NewResultPartitionsWithCreateTables(c *gin.Context) {
	frm := form.NewResultPartitionForm{}

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
	resPart, err := th.dao.NewResultPartitionWithCreateTables(tracer.ExtractTraceCtxFromGin(c), frm.Partition)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "result_partition", resPart)
}

// ListResultPartitions 列出所有结果分区
func (th *TaskHandler) ListResultPartitions(c *gin.Context) {
	pFrm := form.ListResultPartitionsParamsForm{}
	if err := c.ShouldBindQuery(&pFrm); err != nil {
		th.logger.WarnWithFields(logger.Fields{
			"params": c.Params,
			"error":  err,
		}, "An error occurred while c.ShouldBindQuery in TaskHandler.ListResultPartitions, skipped it.")
	}

	// 列出所有分区
	rp, err := th.dao.ListResultPartitions(tracer.ExtractTraceCtxFromGin(c), pFrm.GenQuery())
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "result_partitions", rp)
}

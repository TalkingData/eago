package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/dto"
	"github.com/gin-gonic/gin"
)

// NewResultPartitionsWithCreateTables 新增结果分区并建立结果表和日志表
// 仅测试用，不对外开放的方法
func NewResultPartitionsWithCreateTables(c *gin.Context) {
	var rpForm dto.NewResultPartition

	// 序列化request body
	if err := c.ShouldBindJSON(&rpForm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := rpForm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	rp := dao.NewResultPartitionWithCreateTables(rpForm.Partition)
	// 新建失败
	if rp == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "result_partition", rp)
}

// ListResultPartitions 列出所有结果分区
func ListResultPartitions(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lrp := dto.ListResultPartitionsQuery{}
	if c.ShouldBindQuery(&lrp) == nil {
		_ = lrp.UpdateQuery(query)
	}

	// 列出所有分区
	rp, ok := dao.ListResultPartitions(query)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "result_partitions", rp)
}

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

// KillTask 手动结束任务
func KillTask(c *gin.Context) {
	// 获得分区ID
	rpId, err := strconv.Atoi(c.Param("result_partition_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "result_partition_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 查找分区
	p, ok := dao.GetResultPartitionsPartition(rpId)
	if !ok {
		m := msg.KillTaskPartitionNotFoundFailed
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
		}, m.String())
		m.WriteRest(c)
		return
	}
	// 找不到分区
	if p == "" {
		m := msg.KillTaskPartitionNotFoundFailed
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
		}, m.String())
		m.WriteRest(c)
		return
	}

	// 获得结果ID
	rId, err := strconv.Atoi(c.Param("result_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "result_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 执行手动结束任务
	err = builtin.KillTask(builtin.TaskUniqueIdEncode(p, rId))
	if err != nil {
		m := msg.KillTaskFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListResults 按分区列出所有结果
func ListResults(c *gin.Context) {
	// 获得分区ID
	rpId, err := strconv.Atoi(c.Param("result_partition_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "result_partition_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 根据分区获取结果表前缀
	tableSuffix, ok := dao.GetResultPartitionsPartition(rpId)
	if !ok {
		m := msg.ListResultsPartNotFoundFailed
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	query := dao.Query{}
	// 设置查询filter
	lrq := dto.ListResultsQuery{}
	if c.ShouldBindQuery(&lrq) == nil {
		_ = lrq.UpdateQuery(query)
	}

	// 查询结果
	paged, ok := dao.PagedListResultsByPartition(
		query,
		// 根据分区ID查找Result表名
		tableSuffix,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	// 查询失败
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "results", paged)
}

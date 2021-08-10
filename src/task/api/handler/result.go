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

// KillTask 手动结束任务
// @Summary 手动结束任务
// @Tags 结果
// @Param token header string true "Token"
// @Param result_partition_id path string true "结果分区ID"
// @Param result_id path string true "结果ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /results/{result_partition_id}/{result_id} [DELETE]
func KillTask(c *gin.Context) {
	// 获得分区ID
	rpId, err := strconv.Atoi(c.Param("result_partition_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'result_partition_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	p, ok := model.GetResultPartitionsPartition(rpId)
	if !ok {
		resp := msg.WarnNotFound.GenResponse("Got not success returned status.")
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
			"error":               resp.Error().Error(),
		}, "Error when find record from model.GetResultPartitionTableSuffix.")
		resp.Write(c)
		return
	}
	if p == "" {
		resp := msg.WarnNotFound.GenResponse("Got empty partition.")
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
			"error":               resp.Error().Error(),
		}, "Error when find record from model.GetResultPartitionTableSuffix.")
		resp.Write(c)
		return
	}

	// 获得结果ID
	rId, err := strconv.Atoi(c.Param("result_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'result_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 执行手动结束任务
	err = local.KillTask(local.TaskUniqueIdEncode(p, rId))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field run KillTask.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	msg.Success.GenResponse().Write(c)
}

// ListResults 按分区列出所有结果
// @Summary 按分区列出所有结果
// @Tags 结果
// @Param token header string true "Token"
// @Param status query string false "状态过滤条件"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Param result_partition_id path string true "结果分区ID"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"results":[{"id":1,"task_id":15,"task_name":"task.test_task_2021_03_02_15_32_55","status":0,"worker":"W_119","arguments":"{}","start_at":"2021-03-23 10:42:22","end_at":"2021-03-23T10:55:29+08:00"}],"total":1}"
// @Router /results/{result_partition_id} [GET]
func ListResults(c *gin.Context) {
	query := make(model.Query)

	// 获得分区ID
	rpId, err := strconv.Atoi(c.Param("result_partition_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'result_partition_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	tableSuffix, ok := model.GetResultPartitionsPartition(rpId)
	if !ok {
		resp := msg.WarnNotFound.GenResponse("Find record from model.GetResultPartitionsPartition.")
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
		}, resp.String())
		resp.Write(c)
		return
	}

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query["task_codename LIKE @query"] = sql.Named("query", likeQuery)
	}

	// 状态筛选
	status, err := strconv.Atoi(c.DefaultQuery("status", ""))
	if err == nil {
		query["status = ?"] = status
	}

	paged, ok := model.PagedListResultsByPartition(
		query,
		// 根据分区ID查找Result表名
		tableSuffix,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.PagedListWithPartition.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "results")
	resp.Write(c)
}

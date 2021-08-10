package handler

import (
	"database/sql"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

// NewResultPartitionsWithCreateTables 新建结果分区并建立结果表和日志表
// 仅测试用，不对外开放的方法
func NewResultPartitionsWithCreateTables(c *gin.Context) {
	var rpForm model.ResultPartition

	// 序列化request body
	if err := c.ShouldBindJSON(&rpForm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'partition' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	rp := model.NewResultPartitionWithCreateTables(rpForm.Partition)
	if rp == nil {
		resp := msg.ErrDatabase.GenResponse("Error in model.NewResultPartitionWithCreateTables.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("result_partition", rp)
	resp.Write(c)
	return
}

// ListResultPartitions 列出所有结果分区
// @Summary 列出所有结果分区
// @Tags 结果分区
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success","result_partitions":[{"id":32,"partition":"202103"}]}"
// @Router /result_partitions [GET]
func ListResultPartitions(c *gin.Context) {
	var query = model.Query{}

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query["partition LIKE @query"] = sql.Named("query", likeQuery)

	}

	rp, ok := model.ListResultPartitions(query)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListResultPartitions.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("result_partitions", rp)
	resp.Write(c)
}

package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ListLogs 按分区列出所有结果日志
func ListLogs(c *gin.Context) {

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
		m := msg.ListLogsPartNotFoundFailed
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 获得结果ID
	resultId, err := strconv.Atoi(c.Param("result_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "result_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	query := dao.Query{}
	// 查询结果日志
	query["result_id"] = resultId
	logs, ok := dao.ListLogs(query, tableSuffix)
	// 查询失败
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "logs", logs)
}

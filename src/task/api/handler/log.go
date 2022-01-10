package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/task/conf"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/worker"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
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

	// 查询结果日志
	logs, ok := dao.ListLogs(dao.Query{"result_id=?": resultId}, tableSuffix)
	// 查询失败
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "logs", logs)
}

// WsListLogs 以WebSocket方式按分区ID列出所有结果日志
func WsListLogs(c *gin.Context) {
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

	// 创建默认query
	query := dao.Query{
		"result_id=?": resultId,
		"id>?":        0,
	}

	// 将请求升级为WebSocket
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		_ = ws.Close()
	}()

	for {
		logs, ok := dao.ListLogs(query, tableSuffix)
		if !ok {
			break
		}

		// 循环格式化日志
		for _, l := range *logs {
			content := fmt.Sprintf("[%s] %s", l.CreatedAt.Format(conf.TIMESTAMP_FORMAT), l.Content)
			err = ws.WriteMessage(1, []byte(content))
			if err != nil {
				break
			}
			query["id>?"] = l.Id
		}

		// 获取结果
		r, _ := dao.GetResult(tableSuffix, resultId)
		// 如果获取到结果，并且结果状态不是运行态，则结束循环
		if r != nil && r.Status <= worker.TASK_SUCCESS_END_STATUS {
			break
		}

		time.Sleep(conf.TASK_LOG_REFRESH_INTERVAL_MS * time.Millisecond)
	}
}

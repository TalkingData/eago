package handler

import (
	"eago/common/log"
	"eago/task/conf"
	"eago/task/conf/msg"
	"eago/task/model"
	"eago/task/worker"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

// ListLogs 按分区列出所有结果日志
// @Summary 按分区列出所有结果日志
// @Tags 结果日志
// @Param token header string true "Token"
// @Param result_partition_id path string true "结果分区ID"
// @Param result_id path string true "结果ID"
// @Success 200 {string} string "{"code":0,"message":"Success","logs":[{"id":4,"result_id":1,"content":"Task 1, done.","CreatedAt":"2021-03-23 10:48:22"}]}"
// @Router /logs/{result_partition_id}/{result_id} [GET]
func ListLogs(c *gin.Context) {
	var query = model.Query{}

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

	tableSuffix, ok := model.GetResultPartitionTableSuffix(rpId)
	if !ok {
		resp := msg.WarnNotFound.GenResponse("Find record from model.GetResultPartitionTableSuffix.")
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
		}, resp.String())
		resp.Write(c)
		return
	}

	// 获得结果ID
	resultId, err := strconv.Atoi(c.Param("result_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'result_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}
	query["result_id"] = resultId

	logs, ok := model.ListLogs(query, tableSuffix)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in model.ListLogs.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("logs", logs)
	resp.Write(c)
}

// WsListLogs 以WebSocket方式按分区ID列出所有结果日志
// @Summary 以WebSocket方式按分区ID列出所有结果日志
// @Tags 结果日志
// @Param token header string true "Token"
// @Router /logs/{result_partition_id}/{result_id}/ws [GET]
func WsListLogs(c *gin.Context) {
	var query = model.Query{}

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

	tableSuffix, ok := model.GetResultPartitionTableSuffix(rpId)
	if !ok {
		resp := msg.WarnNotFound.GenResponse("Find record from model.GetResultPartitionTableSuffix.")
		log.WarnWithFields(log.Fields{
			"result_partition_id": rpId,
		}, resp.String())
		resp.Write(c)
		return
	}

	// 获得结果ID
	resultId, err := strconv.Atoi(c.Param("result_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'result_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())

		resp.Write(c)
		return
	}
	query["result_id"] = resultId
	query["id>?"] = 0

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	_, _, _ = ws.ReadMessage()

	for {
		logs, ok := model.ListLogs(query, tableSuffix)
		if !ok {
			break
		}

		for _, l := range *logs {
			msg := fmt.Sprintf("[%s] %s", l.CreatedAt.Format(conf.TIMESTAMP_FORMAT), l.Content)
			err = ws.WriteMessage(1, []byte(msg))
			if err != nil {
				break
			}
			query["id>?"] = l.Id
		}

		r, ok := model.GetResult(tableSuffix, resultId)
		if !ok {
			break
		}
		if r.Status <= worker.TASK_SUCCESS_END_STATUS {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

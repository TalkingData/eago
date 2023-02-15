package handler

import (
	"eago/common/api/ext"
	cMsg "eago/common/code_msg"
	"eago/common/global"
	"eago/common/orm"
	"eago/common/tracer"
	"eago/task/conf/msg"
	"eago/task/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// ListLogs 按分区列出所有结果日志
func (th *TaskHandler) ListLogs(c *gin.Context) {
	// 获得分区ID
	rpId, err := ext.ParamUint32(c, "result_partition_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_partition_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 根据分区获取结果表前缀
	part, err := th.dao.GetResultPartitionsPartition(tracer.ExtractTraceCtxFromGin(c), rpId)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 获得结果ID
	resultId, err := ext.ParamUint32(c, "result_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 查询结果日志
	logs, err := th.dao.ListLogsByPartition(tracer.ExtractTraceCtxFromGin(c), part, resultId)
	// 查询失败
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	ext.WriteSuccessPayload(c, "logs", logs)
}

// WsListLogs 以WebSocket方式按分区ID列出所有结果日志
func (th *TaskHandler) WsListLogs(c *gin.Context) {
	// 获得分区ID
	rpId, err := ext.ParamUint32(c, "result_partition_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_partition_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 根据分区获取结果表前缀
	part, err := th.dao.GetResultPartitionsPartition(tracer.ExtractTraceCtxFromGin(c), rpId)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	if len(part) < 1 {
		m := msg.MsgListLogsPartNotFoundFailed.SetError(err)
		th.logger.ErrorWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 获得结果ID
	resultId, err := ext.ParamUint32(c, "result_id")
	if err != nil {
		m := cMsg.MsgInvalidUriFailed.SetError(err, "result_id")
		th.logger.WarnWithFields(m.ToLoggerFields(), m.GetMsg())
		m.Write2GinCtx(c)
		return
	}

	// 创建默认query
	query := orm.Query{
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

	ctx := tracer.ExtractTraceCtxFromGin(c)
	for {
		logs, err := th.dao.ListLogsByPartition(ctx, part, resultId)
		if err != nil {
			break
		}

		// 循环格式化日志
		for _, l := range logs {
			content := fmt.Sprintf("[%s] %s", l.CreatedAt.Format(global.TimestampFormat), l.Content)
			err = ws.WriteMessage(1, []byte(content))
			if err != nil {
				break
			}
			query["id>?"] = l.Id
		}

		// 获取结果
		resObj, _ := th.dao.GetResult(ctx, part, resultId)
		// 如果获取到结果，并且结果状态不是运行态，则结束循环
		if resObj != nil && resObj.Id > 1 && resObj.Status <= dto.TaskResultStatusSuccessEnd {
			break
		}

		time.Sleep(th.conf.Const.TaskLogRefreshIntervalMs * time.Millisecond)
	}
}

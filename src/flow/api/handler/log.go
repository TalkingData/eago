package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewLog 新建指定流程实例的审批日志
func NewLog(c *gin.Context) {
	// 获得流程实例ID
	insId, err := strconv.Atoi(c.Param("instance_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "instance_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lgFrm dto.NewLog
	// 序列化request body
	if err = c.ShouldBindJSON(&lgFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := lgFrm.Validate(insId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	l := dao.NewLog(insId, lgFrm.Result, *lgFrm.Content, tc["Username"].(string))
	if l == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "log", l)
}

// ListLogs 列出指定流程实例的所有审批日志
func ListLogs(c *gin.Context) {
	// 获得流程实例ID
	insId, err := strconv.Atoi(c.Param("instance_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "instance_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lgFrm dto.ListLog
	// 验证数据
	if m := lgFrm.Validate(insId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	logs, ok := dao.ListLogs(dao.Query{"instance_id=?": insId})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "logs", logs)
}

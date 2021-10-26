package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"eago/flow/srv/builtin"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

// InstantiateFlow 发起流程
func InstantiateFlow(c *gin.Context) {
	fId, err := strconv.Atoi(c.Param("flow_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "flow_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var insFrm dto.InstantiateFlow
	// 序列化request body
	if err = c.ShouldBindJSON(&insFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := insFrm.Validate(fId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 反序列化FormData
	fData := make(map[string]interface{})
	if err = json.Unmarshal([]byte(*insFrm.FormData), &fData); err != nil {
		m := msg.SerializeFailed.SetError(err, "无法反序列化表单数据内容")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	// 将创建人信息存入FormData
	fData["user_id"] = tc["UserId"].(int32)
	fData["username"] = tc["Username"].(string)
	fData["phone"] = tc["Phone"].(string)

	// 序列化FormData
	fDataStr, err := json.Marshal(fData)
	if err != nil {
		m := msg.SerializeFailed.SetError(err, "无法序列化表单数据内容")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 调用流程服务，发起流程，返回流程实例ID
	insId, err := builtin.InstantiateFlow(insFrm.FormId, string(fDataStr), tc["Username"].(string))
	if err != nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "instance_id", insId)
}

// NewFlow 新建流程
func NewFlow(c *gin.Context) {
	var flFrm dto.NewFlow

	// 序列化request body
	if err := c.ShouldBindJSON(&flFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := flFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	fl := dao.NewFlow(
		flFrm.Name,
		flFrm.CategoriesId,
		*flFrm.Description,
		*flFrm.Disabled,
		flFrm.FormId,
		flFrm.FirstNodeId,
		tc["Username"].(string),
	)
	// 新建失败
	if fl == nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "flow", fl)
}

// RemoveFlow 删除流程
func RemoveFlow(c *gin.Context) {
	fId, err := strconv.Atoi(c.Param("flow_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "flow_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rfFrm dto.RemoveFlow
	// 验证数据
	if m := rfFrm.Validate(fId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveFlow(fId); !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetFlow 更新流程
func SetFlow(c *gin.Context) {

	fId, err := strconv.Atoi(c.Param("flow_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "flow_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var flFrm dto.SetFlow
	// 序列化request body
	if err = c.ShouldBindJSON(&flFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := flFrm.Validate(fId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	f, ok := dao.SetFlow(
		fId,
		flFrm.Name,
		flFrm.CategoriesId,
		*flFrm.Description,
		*flFrm.Disabled,
		flFrm.FormId,
		flFrm.FirstNodeId,
		tc["Username"].(string),
	)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "flow", f)
}

// ListFlows 列出所有流程
func ListFlows(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lfq := dto.ListFlowsQuery{}
	if c.ShouldBindQuery(&lfq) == nil {
		_ = lfq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListFlows(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while model.PagedListFlows.")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "flows", paged)
}

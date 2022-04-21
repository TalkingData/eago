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

// NewForm 新建表单
func NewForm(c *gin.Context) {
	var frm dto.NewForm

	// 序列化request body
	if err := c.ShouldBindJSON(&frm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := frm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	f := dao.NewForm(frm.Name, *frm.Disabled, *frm.Description, *frm.Body, tc["Username"].(string))
	if f == nil {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "form", f)
}

// SetForm 更新表单
func SetForm(c *gin.Context) {
	frmId, err := strconv.Atoi(c.Param("form_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "form_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var setFrm dto.SetForm
	// 序列化request body
	if err := c.ShouldBindJSON(&setFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := setFrm.Validate(frmId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	f, ok := dao.SetForm(frmId, setFrm.Name, *setFrm.Disabled, *setFrm.Description, tc["Username"].(string))
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "form", f)
}

// GetForm 获取指定表单
func GetForm(c *gin.Context) {
	frmId, err := strconv.Atoi(c.Param("form_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "form_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var gFrm dto.GetForm
	// 验证数据
	if m := gFrm.Validate(frmId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	f, ok := dao.GetForm(dao.Query{"id=?": frmId})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "form", f)
}

// PagedListForms 列出所有表单-分页
func PagedListForms(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lfq := dto.PagedListFormsQuery{}
	if c.ShouldBindQuery(&lfq) == nil {
		_ = lfq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListForms(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "forms", paged)
}

// ListFormFlows 列出表单所关联流程
func ListFormFlows(c *gin.Context) {
	frmId, err := strconv.Atoi(c.Param("form_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "form_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lfrFrm dto.ListFormRelations
	// 验证数据
	if m := lfrFrm.Validate(frmId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	fl, ok := dao.ListFlows(dao.Query{"form_id=?": frmId})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "flows", fl)
}

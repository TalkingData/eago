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

// NewCategory 新建类别
func NewCategory(c *gin.Context) {
	var cFrm dto.NewCategory

	// 序列化request body
	if err := c.ShouldBindJSON(&cFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := cFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	tc := c.GetStringMap("TokenContent")
	cat := dao.NewCategory(cFrm.Name, tc["Username"].(string))
	// 新建失败
	if cat == nil {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "category", cat)
}

// RemoveCategory 删除类别
func RemoveCategory(c *gin.Context) {
	cId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "category_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rcFrm dto.RemoveCategory
	// 验证数据
	if m := rcFrm.Validate(cId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if ok := dao.RemoveCategory(cId); !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetCategory 更新类别
func SetCategory(c *gin.Context) {
	cId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "category_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var cFrm dto.SetCategory
	// 序列化request body
	if err = c.ShouldBindJSON(&cFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := cFrm.Validate(cId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	tc := c.GetStringMap("TokenContent")
	cat, ok := dao.SetCategory(cId, cFrm.Name, tc["Username"].(string))
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "category", cat)
}

// ListCategories 列出所有类别
func ListCategories(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lcq := dto.ListCategoriesQuery{}
	if c.ShouldBindQuery(&lcq) == nil {
		_ = lcq.UpdateQuery(query)
	}

	cs, ok := dao.ListCategories(query)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "categories", cs)
}

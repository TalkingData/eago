package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewForm struct
type NewForm struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	Body        *string `json:"body" valid:"MinSize(2)"`
}

func (n *NewForm) Valid(v *validation.Validation) {
	if ct, _ := dao.GetFormCount(dao.Query{"name=?": n.Name}); ct > 0 {
		_ = v.SetError("Name", "表单名称已存在")
	}
}

func (n *NewForm) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(n)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// SetForm struct
type SetForm struct {
	formId int

	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	Disabled    *bool   `json:"disabled" valid:"Required"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (s *SetForm) Valid(v *validation.Validation) {
	if ct, _ := dao.GetFormCount(dao.Query{"name=?": s.Name, "id<>?": s.formId}); ct > 0 {
		_ = v.SetError("Name", "表单名称已存在")
	}
}

func (s *SetForm) Validate(frmId int) *message.Message {
	s.formId = frmId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(s)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// GetForm struct
type GetForm struct{}

func (*GetForm) Validate(frmId int) *message.Message {
	// 验证表单是否存在
	if ct, _ := dao.GetFormCount(dao.Query{"id=?": frmId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("表单不存在")
	}

	return nil
}

// PagedListFormsQuery struct
type PagedListFormsQuery struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (q *PagedListFormsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(id LIKE @query OR "+
			"name LIKE @query OR "+
			"description LIKE @query OR "+
			"created_by LIKE @query OR "+
			"updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Disabled != nil {
		query["disabled=?"] = *q.Disabled
	}

	return nil
}

// ListFormRelations struct
type ListFormRelations struct{}

func (*ListFormRelations) Validate(frmId int) *message.Message {
	// 验证表单是否存在
	if ct, _ := dao.GetFormCount(dao.Query{"id=?": frmId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("表单不存在")
	}

	return nil
}

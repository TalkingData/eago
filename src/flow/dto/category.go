package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewCategory struct
type NewCategory struct {
	Name string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
}

// Valid
func (n *NewCategory) Valid(v *validation.Validation) {
	if ct, _ := dao.GetCategoriesCount(dao.Query{"name=?": n.Name}); ct > 0 {
		_ = v.SetError("Name", "类别名称已存在")
	}
}

// Validate
func (n *NewCategory) Validate() *message.Message {
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

// RemoveCategory struct
type RemoveCategory struct{}

// Validate
func (*RemoveCategory) Validate(cId int) *message.Message {
	// 验证类别是否存在
	if ct, _ := dao.GetCategoriesCount(dao.Query{"id=?": cId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("类别不存在")
	}

	// 验证类别与流程关联
	if ct, _ := dao.GetFlowCount(dao.Query{"id=?": cId}); ct > 0 {
		return msg.AssociatedCategoryFlowFailed
	}

	return nil
}

// SetCategory struct
type SetCategory struct {
	categoryId int

	Name string `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
}

// Valid
func (s *SetCategory) Valid(v *validation.Validation) {
	if ct, _ := dao.GetCategoriesCount(dao.Query{"name=?": s.Name, "id<>?": s.categoryId}); ct > 0 {
		_ = v.SetError("Name", "类别名称已存在")
	}
}

// Validate
func (s *SetCategory) Validate(cId int) *message.Message {
	s.categoryId = cId
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

// ListCategoriesQuery struct
type ListCategoriesQuery struct {
	Query *string `form:"query"`
}

// UpdateQuery
func (q *ListCategoriesQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["name LIKE @query OR id LIKE @query OR created_by LIKE @query OR updated_by LIKE @query"] = sql.Named("query", likeQuery)
	}

	return nil
}

// ListCategoriesRelations struct
type ListCategoriesRelations struct{}

// Validate
func (*ListCategoriesRelations) Validate(cId int) *message.Message {
	// 验证类别是否存在
	if ct, _ := dao.GetCategoriesCount(dao.Query{"id=?": cId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("类别不存在")
	}

	return nil
}

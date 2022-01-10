package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// InstantiateFlow struct
type InstantiateFlow struct {
	flowId int

	FormId   int
	Name     string
	FormData *string `json:"form_data" valid:"MinSize(2)"`
}

// Validate
func (i *InstantiateFlow) Validate(fId int) *message.Message {
	flow, ok := dao.GetFlow(dao.Query{"id=?": fId})
	if !ok {
		return msg.UnknownError
	}

	if flow == nil || flow.Id == 0 {
		return msg.NotFoundFailed.SetDetail("流程不存在")
	}

	if *flow.Disabled {
		return msg.NotFoundFailed.SetDetail("无法发起一个禁用的流程")
	}

	i.Name = flow.Name
	i.FormId = flow.Id

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(i)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type NewFlow struct {
	Name          string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	InstanceTitle string  `json:"instance_title" valid:"Required;MinSize(3);MaxSize(200)"`
	CategoriesId  *int    `json:"categories_id"`
	Disabled      *bool   `json:"disabled" valid:"Required"`
	Description   *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	FormId        int     `json:"form_id" valid:"Required"`
	FirstNodeId   int     `json:"first_node_id" valid:"Required"`
}

// Valid
func (n *NewFlow) Valid(v *validation.Validation) {
	if ct, _ := dao.GetFlowCount(dao.Query{"name=?": n.Name}); ct > 0 {
		_ = v.SetError("Name", "流程名称已存在")
	}

	if ct, _ := dao.GetFormCount(dao.Query{"id=?": n.FormId, "disabled=?": false}); ct < 1 {
		_ = v.SetError("FormId", "找不到所选表单，请确定该表单存在并且不是禁用状态")
	}
}

// Validate
func (n *NewFlow) Validate() *message.Message {
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

// RemoveFlow struct
type RemoveFlow struct{}

// Validate
func (*RemoveFlow) Validate(tId int) *message.Message {
	// 验证流程是否存在
	if ct, _ := dao.GetFlowCount(dao.Query{"id=?": tId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("流程不存在")
	}

	return nil
}

// SetFlow struct
type SetFlow struct {
	flowId int

	Name          string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	InstanceTitle string  `json:"instance_title" valid:"Required;MinSize(3);MaxSize(200)"`
	CategoriesId  *int    `json:"categories_id"`
	Disabled      *bool   `json:"disabled" valid:"Required"`
	Description   *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	FormId        int     `json:"form_id" valid:"Required"`
	FirstNodeId   int     `json:"first_node_id" valid:"Required"`
}

// Valid
func (s *SetFlow) Valid(v *validation.Validation) {
	if ct, _ := dao.GetFlowCount(dao.Query{"name=?": s.Name, "id<>?": s.flowId}); ct > 0 {
		_ = v.SetError("Name", "流程名称已存在")
	}

	if ct, _ := dao.GetFormCount(dao.Query{"id=?": s.FormId, "disabled=?": false}); ct < 1 {
		_ = v.SetError("FormId", "找不到所选表单，请确定该表单存在并且不是禁用状态")
	}
}

// Validate
func (s *SetFlow) Validate(flId int) *message.Message {
	s.flowId = flId
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

// ListFlowsQuery struct
type ListFlowsQuery struct {
	Query        *string `form:"query"`
	Disabled     *bool   `form:"disabled"`
	CategoriesId *int    `form:"categories_id"`
}

// UpdateQuery
func (q *ListFlowsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(flows.id LIKE @query OR "+
			"flows.name LIKE @query OR "+
			"flows.description LIKE @query OR "+
			"flows.created_by LIKE @query OR "+
			"flows.updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Disabled != nil {
		query["flows.disabled=?"] = *q.Disabled
	}

	if q.CategoriesId != nil {
		query["flows.categories_id=?"] = *q.CategoriesId
	}

	return nil
}

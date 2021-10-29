package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewTrigger struct
type NewTrigger struct {
	Name         string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	TaskCodename string  `json:"task_codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2);MaxSize(4000)"`
}

// Valid
func (nt *NewTrigger) Valid(v *validation.Validation) {
	if ct, _ := dao.GetTriggerCount(dao.Query{"name=?": nt.Name}); ct > 0 {
		_ = v.SetError("Name", "触发器名称已存在")
	}
}

// Validate
func (nt *NewTrigger) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(nt)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveTrigger struct
type RemoveTrigger struct{}

// Validate
func (*RemoveTrigger) Validate(tId int) *message.Message {
	// 验证触发器是否存在
	if ct, _ := dao.GetTriggerCount(dao.Query{"id=?": tId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("触发器不存在")
	}

	// 验证触发器是否有关联存在
	if ct, _ := dao.GetNodeTriggerCount(dao.Query{"trigger_id=?": tId}); ct > 0 {
		return msg.AssociatedTriggerNodeFailed
	}

	return nil
}

// SetTrigger struct
type SetTrigger struct {
	triggerId int

	Name         string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z0-9-\u4e00-\u9fa5]+$/)"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
	TaskCodename string  `json:"task_codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2);MaxSize(4000)"`
}

// Valid
func (sg *SetTrigger) Valid(v *validation.Validation) {
	if ct, _ := dao.GetTriggerCount(dao.Query{"name=?": sg.Name, "id<>?": sg.triggerId}); ct > 0 {
		_ = v.SetError("Name", "触发器名称已存在")
	}
}

// Validate
func (sg *SetTrigger) Validate(tId int) *message.Message {
	sg.triggerId = tId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(sg)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListTriggersQuery struct
type ListTriggersQuery struct {
	Query *string `form:"query"`
}

// UpdateQuery
func (ltq *ListTriggersQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if ltq.Query != nil && *ltq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *ltq.Query)
		query["name LIKE @query OR description LIKE @query OR task_codename LIKE @query OR id LIKE @query OR created_by LIKE @query OR updated_by LIKE @query"] = sql.Named("query", likeQuery)
	}

	return nil
}

// ListTriggerNodes struct
type ListTriggerNodes struct{}

// Validate
func (*ListTriggerNodes) Validate(tId int) *message.Message {
	// 验证触发器是否存在
	if ct, _ := dao.GetTriggerCount(dao.Query{"id=?": tId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("触发器不存在")
	}

	return nil
}

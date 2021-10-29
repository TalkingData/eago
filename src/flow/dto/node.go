package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewNode struct {
	Name              string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId          *int    `json:"parent_id"`
	Category          int     `json:"category" valid:"Required"`
	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`
	VisibleFields     string  `json:"visible_fields" valid:"Required;MinSize(2)"`
	EditableFields    string  `json:"editable_fields" valid:"Required;MinSize(2)"`
}

// Validate
func (n *NewNode) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据¬
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

// RemoveNode struct
type RemoveNode struct{}

// Validate
func (*RemoveNode) Validate(nId int) *message.Message {
	// 验证节点是否存在
	if ct, _ := dao.GetNodeCount(dao.Query{"id=?": nId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("节点不存在")
	}

	// 验证是否存在父节点关联
	if ct, _ := dao.GetNodeCount(dao.Query{"parent_id=?": nId}); ct > 0 {
		return msg.AssociatedNodeFailed
	}

	// 验证节点是否存在触发器关联
	if ct, _ := dao.GetNodeTriggerCount(dao.Query{"node_id=?": nId}); ct > 0 {
		return msg.AssociatedNodeTriggerFailed
	}

	// 验证节是否有关联存在
	if ct, _ := dao.GetFlowCount(dao.Query{"first_node_id=?": nId}); ct < 1 {
		return msg.AssociatedNodeFlowFailed
	}

	return nil
}

type SetNode struct {
	Name              string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId          *int    `json:"parent_id"`
	Category          int     `json:"category" valid:"Required"`
	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`
	VisibleFields     string  `json:"visible_fields" valid:"Required;MinSize(2)"`
	EditableFields    string  `json:"editable_fields" valid:"Required;MinSize(2)"`
}

// Validate
func (s *SetNode) Validate(nId int) *message.Message {
	// 验证节是否存在
	if ct, _ := dao.GetNodeCount(dao.Query{"id=?": nId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("节点不存在")
	}

	valid := validation.Validation{}
	// 验证数据¬
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

// ListNodesQuery struct
type ListNodesQuery struct {
	Query    *string `form:"query"`
	Category *bool   `form:"category"`
}

// UpdateQuery
func (q *ListNodesQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(name LIKE @query OR description LIKE @query OR id LIKE @query OR created_by LIKE @query OR updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Category != nil {
		query["category=?"] = *q.Category
	}

	return nil
}

// AddTrigger2Node struct
type AddTrigger2Node struct {
	nodeId int

	TriggerId int `json:"trigger_id" valid:"Required"`
}

// Valid
func (apr *AddTrigger2Node) Valid(v *validation.Validation) {
	// 验证触发器是否存在
	if ct, _ := dao.GetTriggerCount(dao.Query{"id=?": apr.TriggerId}); ct < 1 {
		_ = v.SetError("TriggerId", "触发器不存在")
	}

	// 验证触发器是否已经属于该节点
	if ct, _ := dao.GetNodeTriggerCount(dao.Query{"node_id": apr.TriggerId, "trigger_id=?": apr.nodeId}); ct > 0 {
		_ = v.SetError("TriggerId", "触发器已经属于该节点")
	}
}

// Validate
func (apr *AddTrigger2Node) Validate(nId int) *message.Message {
	// 验证节点是否存在
	if ct, _ := dao.GetNodeCount(dao.Query{"id=?": nId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("节点不存在")
	}

	apr.nodeId = nId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(apr)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveNodeTrigger struct
type RemoveNodeTrigger struct{}

// Validate
func (*RemoveNodeTrigger) Validate(nId, tId int) *message.Message {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetNodeTriggerCount(dao.Query{"node_id": nId, "trigger_id=?": tId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

// ListNodeRelations struct
type ListNodeRelations struct{}

// Validate
func (*ListNodeRelations) Validate(nId int) *message.Message {
	// 验证节点是否存在
	if ct, _ := dao.GetNodeCount(dao.Query{"id=?": nId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("节点不存在")
	}

	return nil
}

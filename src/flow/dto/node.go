package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/model"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"regexp"
	"strings"
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

func (n *NewNode) Valid(v *validation.Validation) {
	// 验证AssigneeCondition
	if m := validateAssigneeCondition(n.Category, n.AssigneeCondition); m != "" {
		_ = v.SetError("AssigneeCondition", m)
	}
}

func (n *NewNode) Validate() *message.Message {
	if n.ParentId != nil {
		// 验证父节点是否存在
		if ct, _ := dao.GetNodeCount(dao.Query{"id=?": n.ParentId}); ct < 1 {
			return msg.AssociatedParentNodeNotFoundFailed
		}

		// 验证是否重复关联父节点
		if ct, _ := dao.GetNodeCount(dao.Query{"parent_id=?": n.ParentId}); ct > 0 {
			return msg.AssociatedParentNodeDuplicateFailed
		}
	}

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

// RemoveNode struct
type RemoveNode struct{}

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
	if ct, _ := dao.GetFlowCount(dao.Query{"first_node_id=?": nId}); ct > 0 {
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

func (s *SetNode) Valid(v *validation.Validation) {
	// 验证AssigneeCondition
	if m := validateAssigneeCondition(s.Category, s.AssigneeCondition); m != "" {
		_ = v.SetError("AssigneeCondition", m)
	}
}

func (s *SetNode) Validate(nId int) *message.Message {
	if s.ParentId != nil {
		// 验证父节点是否存在
		if ct, _ := dao.GetNodeCount(dao.Query{"id=?": s.ParentId}); ct < 1 {
			return msg.AssociatedParentNodeNotFoundFailed
		}

		// 验证是否重复关联父节点
		if ct, _ := dao.GetNodeCount(dao.Query{"parent_id=?": s.ParentId, "id<>?": nId}); ct > 0 {
			return msg.AssociatedParentNodeDuplicateFailed
		}

		// 不能将自身节点关联为父节点
		if *s.ParentId == nId {
			return msg.AssociatedParentNodeSelfFailed
		}
	}

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

// PagedListNodesQuery struct
type PagedListNodesQuery struct {
	Query    *string `form:"query"`
	Name     *string `form:"name"`
	ParentId *string `form:"parent_id"`
	Category *int    `form:"category"`
}

func (q *PagedListNodesQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(nodes.id LIKE @query OR "+
			"nodes.name LIKE @query OR "+
			"nodes.created_by LIKE @query OR "+
			"nodes.updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Name != nil {
		query["nodes.name LIKE ?"] = fmt.Sprintf("%%%s%%", *q.Name)
	}

	if q.Category != nil {
		query["nodes.category=?"] = *q.Category
	}

	if q.ParentId != nil {
		if strings.ToLower(*q.ParentId) == "null" {
			query["nodes.parent_id"] = nil
		} else {
			query["nodes.parent_id=?"] = *q.ParentId
		}
	}

	return nil
}

// AddTrigger2Node struct
type AddTrigger2Node struct {
	nodeId int

	TriggerId int `json:"trigger_id" valid:"Required"`
}

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

func (*RemoveNodeTrigger) Validate(nId, tId int) *message.Message {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetNodeTriggerCount(dao.Query{"node_id": nId, "trigger_id=?": tId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

// ListNodeRelations struct
type ListNodeRelations struct{}

func (*ListNodeRelations) Validate(nId int) *message.Message {
	// 验证节点是否存在
	if ct, _ := dao.GetNodeCount(dao.Query{"id=?": nId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("节点不存在")
	}

	return nil
}

// validateAssigneeCondition 验证审批人条件是否合法
func validateAssigneeCondition(category int, acStr *string) string {
	// 首节点不需要有AssigneeCondition
	if category == conf.NODE_CATEGORY_FIRST {
		if acStr != nil && *acStr != "{}" {
			return "不需要设置审批人条件，请设置其为空对象"
		}
		*acStr = "{}"
		return ""
	}

	// 反序列化AssigneeCondition
	acObj := model.AssigneeCondition{}
	err := json.Unmarshal([]byte(*acStr), &acObj)
	// 反序列化错误
	if err != nil {
		return "格式不合法，必须包含：condition、getter、data属性"
	}

	// Condition表达式不合法
	if acObj.Condition != conf.AC_INITIATOR &&
		acObj.Condition != conf.AC_INITIATORS_DEPARTMENTS_OWNER &&
		acObj.Condition != conf.AC_INITIATORS_PARENT_DEPARTMENTS_OWNER &&
		acObj.Condition != conf.AC_SPECIFIED_USERS &&
		acObj.Condition != conf.AC_SPECIFIED_PRODUCT_OWNER &&
		acObj.Condition != conf.AC_SPECIFIED_GROUP_OWNER &&
		acObj.Condition != conf.AC_SPECIFIED_DEPARTMENT_OWNER &&
		acObj.Condition != conf.AC_SPECIFIED_ROLE {
		return "不支持所输入的condition表达式"
	}

	// Getter不合法
	if acObj.Getter != conf.GETTER_DIRECT &&
		acObj.Getter != conf.GETTER_FIELD {
		return "不支持所输入的getter方法"
	}

	// Data长度不合法
	if len(acObj.Data) < 3 {
		return "所输入的data长度不能小于3个字符"
	}

	// Data不合法
	r, _ := regexp.Compile("^[a-zA-Z-_][a-zA-Z0-9-.,_@]{1,}$")
	if !r.Match([]byte(acObj.Data)) {
		return "所输入的data只能包含：英文大小写字母和\"@,.-_\"，并且以英文大小写字母开头"
	}

	return ""
}

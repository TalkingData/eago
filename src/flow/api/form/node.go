package form

import (
	"context"
	"database/sql"
	cMsg "eago/common/code_msg"
	"eago/common/orm"
	"eago/flow/conf/msg"
	"eago/flow/dao"
	"eago/flow/dto"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"regexp"
	"strings"
)

type NewNodeForm struct {
	Name              string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId          *uint32 `json:"parent_id"`
	Category          int32   `json:"category" valid:"Required"`
	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`
	VisibleFields     string  `json:"visible_fields" valid:"Required;MinSize(2)"`
	EditableFields    string  `json:"editable_fields" valid:"Required;MinSize(2)"`
}

func (f *NewNodeForm) Valid(v *validation.Validation) {
	// 验证AssigneeCondition
	if m := validateAssigneeCondition(f.Category, f.AssigneeCondition); m != "" {
		_ = v.SetError("AssigneeCondition", m)
	}
}

func (f *NewNodeForm) Validate(ctx context.Context, dao *dao.Dao) *cMsg.CodeMsg {
	if f.ParentId != nil {
		// 验证父节点是否存在
		if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": f.ParentId}); !exist {
			return msg.MsgAssociatedParentNodeNotFoundFailed
		}

		// 验证是否重复关联父节点
		if exist, _ := dao.IsNodeExist(ctx, orm.Query{"parent_id=?": f.ParentId}); exist {
			return msg.MsgAssociatedParentNodeDuplicateFailed
		}
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type RemoveNodeForm struct{}

func (*RemoveNodeForm) Validate(ctx context.Context, dao *dao.Dao, nodeId uint32) *cMsg.CodeMsg {
	// 验证节点是否存在
	if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": nodeId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("节点不存在")
	}

	// 验证是否存在父节点关联
	if exist, _ := dao.IsNodeExist(ctx, orm.Query{"parent_id=?": nodeId}); exist {
		return msg.MsgAssociatedNodeFailed
	}

	// 验证节点是否存在触发器关联
	if ct, _ := dao.GetNodesTriggerCount(ctx, orm.Query{"node_id=?": nodeId}); ct > 0 {
		return msg.MsgAssociatedNodeTriggerFailed
	}

	// 验证节是否有关联存在
	if exist, _ := dao.IsFlowExist(ctx, orm.Query{"first_node_id=?": nodeId}); exist {
		return msg.MsgAssociatedNodeFlowFailed
	}

	return nil
}

type SetNodeForm struct {
	Name              string  `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId          *uint32 `json:"parent_id"`
	Category          int32   `json:"category" valid:"Required"`
	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`
	VisibleFields     string  `json:"visible_fields" valid:"Required;MinSize(2)"`
	EditableFields    string  `json:"editable_fields" valid:"Required;MinSize(2)"`
}

func (s *SetNodeForm) Valid(v *validation.Validation) {
	// 验证AssigneeCondition
	if m := validateAssigneeCondition(s.Category, s.AssigneeCondition); m != "" {
		_ = v.SetError("AssigneeCondition", m)
	}
}

func (s *SetNodeForm) Validate(ctx context.Context, dao *dao.Dao, nodeId uint32) *cMsg.CodeMsg {
	if s.ParentId != nil {
		// 验证父节点是否存在
		if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": s.ParentId}); !exist {
			return msg.MsgAssociatedParentNodeNotFoundFailed
		}

		// 验证是否重复关联父节点
		if exist, _ := dao.IsNodeExist(ctx, orm.Query{"parent_id=?": s.ParentId, "id<>?": nodeId}); exist {
			return msg.MsgAssociatedParentNodeDuplicateFailed
		}

		// 不能将自身节点关联为父节点
		if *s.ParentId == nodeId {
			return msg.MsgAssociatedParentNodeSelfFailed
		}
	}

	// 验证节是否存在
	if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": nodeId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("节点不存在")
	}

	valid := validation.Validation{}
	// 验证数据¬
	ok, err := valid.Valid(s)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type PagedListNodesParamsForm struct {
	Query    *string `form:"query"`
	Name     *string `form:"name"`
	ParentId *string `form:"parent_id"`
	Category *int    `form:"category"`
}

func (pf *PagedListNodesParamsForm) GenQuery() orm.Query {
	query := orm.Query{}

	// 通用Query
	if pf.Query != nil && *pf.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *pf.Query)
		query["(nodes.id LIKE @query OR "+
			"nodes.name LIKE @query OR "+
			"nodes.created_by LIKE @query OR "+
			"nodes.updated_by LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if pf.Name != nil {
		query["nodes.name LIKE ?"] = fmt.Sprintf("%%%s%%", *pf.Name)
	}

	if pf.Category != nil {
		query["nodes.category=?"] = *pf.Category
	}

	if pf.ParentId != nil {
		if strings.ToLower(*pf.ParentId) == "null" {
			query["nodes.parent_id"] = nil
		} else {
			query["nodes.parent_id=?"] = *pf.ParentId
		}
	}

	return query
}

type AddTrigger2NodeForm struct {
	TriggerId uint32 `json:"trigger_id" valid:"Required"`

	nodeId uint32

	dao *dao.Dao
	ctx context.Context
}

func (f *AddTrigger2NodeForm) Valid(v *validation.Validation) {
	// 验证触发器是否存在
	if exist, _ := f.dao.IsTriggerExist(f.ctx, orm.Query{"id=?": f.TriggerId}); !exist {
		_ = v.SetError("TriggerId", "触发器不存在")
	}

	// 验证触发器是否已经属于该节点
	if ct, _ := f.dao.GetNodesTriggerCount(f.ctx, orm.Query{"node_id": f.TriggerId, "trigger_id=?": f.nodeId}); ct > 0 {
		_ = v.SetError("TriggerId", "触发器已经属于该节点")
	}
}

func (f *AddTrigger2NodeForm) Validate(ctx context.Context, dao *dao.Dao, nodeId uint32) *cMsg.CodeMsg {
	// 验证节点是否存在
	if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": nodeId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("节点不存在")
	}

	f.nodeId = nodeId

	f.ctx = ctx
	f.dao = dao

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(f)
	if err != nil {
		return cMsg.MsgValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return cMsg.MsgValidateFailed.SetError(valid.Errors)
	}

	return nil
}

type RemoveNodesTriggerForm struct{}

func (*RemoveNodesTriggerForm) Validate(ctx context.Context, dao *dao.Dao, nodeId, triId uint32) *cMsg.CodeMsg {
	// 验证用户是否已经属于该产品线
	if ct, _ := dao.GetNodesTriggerCount(ctx, orm.Query{"node_id": nodeId, "trigger_id=?": triId}); ct < 1 {
		return cMsg.MsgNotFoundFailed
	}

	return nil
}

type ListNodesRelationsForm struct{}

func (*ListNodesRelationsForm) Validate(ctx context.Context, dao *dao.Dao, nodeId uint32) *cMsg.CodeMsg {
	// 验证节点是否存在
	if exist, _ := dao.IsNodeExist(ctx, orm.Query{"id=?": nodeId}); !exist {
		return cMsg.MsgNotFoundFailed.SetDetail("节点不存在")
	}

	return nil
}

// validateAssigneeCondition 验证审批人条件是否合法
func validateAssigneeCondition(category int32, acStr *string) string {
	// 首节点不需要有AssigneeCondition
	if category == dto.NodeCategoryFirst {
		if acStr != nil && *acStr != "{}" {
			return "不需要设置审批人条件，请设置其为空对象"
		}
		*acStr = "{}"
		return ""
	}

	// 反序列化AssigneeCondition
	acObj := dto.AssigneeCondition{}
	err := json.Unmarshal([]byte(*acStr), &acObj)
	// 反序列化错误
	if err != nil {
		return "格式不合法，必须包含：condition、getter、data属性"
	}

	// Condition表达式不合法
	if _, ok := dto.AssigneeConditionsAllowed[acObj.Condition]; !ok {
		return "不支持所输入的condition表达式"

	}

	// Getter不合法
	if _, ok := dto.GetterAllowed[acObj.Getter]; !ok {
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

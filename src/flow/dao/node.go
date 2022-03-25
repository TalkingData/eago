package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewNode 新增节点
func NewNode(name string, parentId *int, category int, entryCondition, assigneeCondition *string, vFields, eFields, createdBy string) *model.Node {
	log.Info("dao.NewNode called.")
	defer log.Info("dao.NewNode end.")

	var n = model.Node{
		Name:              name,
		ParentId:          parentId,
		Category:          category,
		EntryCondition:    entryCondition,
		AssigneeCondition: assigneeCondition,
		VisibleFields:     vFields,
		EditableFields:    eFields,
		CreatedBy:         createdBy,
	}

	if res := db.Create(&n); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":               name,
			"parent_id":          parentId,
			"category":           category,
			"entry_condition":    entryCondition,
			"assignee_condition": assigneeCondition,
			"visible_fields":     vFields,
			"editable_fields":    eFields,
			"created_by":         createdBy,
			"error":              res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &n
}

// RemoveNode 删除节点
func RemoveNode(nId int) bool {
	res := db.Delete(model.Node{}, "id=?", nId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    nId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetNode 更新节点
func SetNode(id int, name string, parentId *int, category int, entryCondition, assigneeCondition *string, vFields, eFields, updatedBy string) (*model.Node, bool) {
	log.Info("dao.SetNode called.")
	defer log.Info("dao.SetNode end.")

	n := model.Node{}

	res := db.Model(&model.Node{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":               name,
			"parent_id":          parentId,
			"category":           category,
			"entry_condition":    entryCondition,
			"assignee_condition": assigneeCondition,
			"visible_fields":     vFields,
			"editable_fields":    eFields,
			"updated_by":         updatedBy,
		}).
		First(&n)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":               name,
			"parent_id":          parentId,
			"category":           category,
			"entry_condition":    entryCondition,
			"assignee_condition": assigneeCondition,
			"visible_fields":     vFields,
			"editable_fields":    eFields,
			"updated_by":         updatedBy,
			"error":              res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &n, true
}

// GetNode 查询单个节点
func GetNode(query Query) (*model.Node, bool) {
	log.Info("dao.GetNode called.")
	defer log.Info("dao.GetNode end.")

	var (
		n = model.Node{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&n); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &n, true
}

// GetNodeCount 查询节点数量
func GetNodeCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Node{})

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Count.")
		return count, false
	}
	return count, true
}

// ListNodes 查询节点
func ListNodes(query Query) ([]model.Node, bool) {
	var d = db
	ns := make([]model.Node, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ns); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found")
			return ns, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return ns, true
}

// PagedListNodes 查询节点-分页
func PagedListNodes(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	log.Info("dao.PagedListNodes called.")
	defer log.Info("dao.PagedListNodes end.")

	var d = db.Model(&model.Node{}).Select(
		"nodes.id, " +
			"nodes.name, " +
			"nodes.parent_id, " +
			"p.name AS parent_name, " +
			"nodes.category, " +
			"nodes.entry_condition, " +
			"nodes.assignee_condition, " +
			"nodes.visible_fields, " +
			"nodes.editable_fields, " +
			"nodes.created_at, " +
			"nodes.created_by, " +
			"nodes.updated_at, " +
			"nodes.updated_by",
	).Joins("LEFT JOIN nodes AS p ON p.id = nodes.parent_id")
	ns := make([]model.ListNodes, pageSize)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ns)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// Node2Chain 将节点转化为链结构
func Node2Chain(n *model.Node) *model.NodeChain {
	ts, ok := ListNodeTriggers(n.Id)
	triggers := make([]model.NodesTrigger, 0)
	if ok && len(ts) > 0 {
		triggers = append(triggers, ts...)
	}

	return &model.NodeChain{
		Id:                n.Id,
		Name:              n.Name,
		Category:          n.Category,
		EntryCondition:    *n.EntryCondition,
		AssigneeCondition: *n.AssigneeCondition,
		Assignees:         nil,
		Triggers:          triggers,
		VisibleFields:     n.VisibleFields,
		EditableFields:    n.EditableFields,
		SubNode:           nil,
	}
}

// GetNodeChain 递推列出节点链
func GetNodeChain(parentNode *model.NodeChain) bool {
	node, ok := GetNode(Query{"parent_id": parentNode.Id})
	if !ok {
		return false
	}
	if node == nil {
		return true
	}
	subNode := Node2Chain(node)
	parentNode.SubNode = subNode
	return GetNodeChain(subNode)
}

// AddNodeTrigger 关联表操作::添加触发器至节点
func AddNodeTrigger(nodeId, triggerId int, createdBy string) bool {
	log.Info("dao.AddNodeTrigger called.")
	defer log.Info("dao.AddNodeTrigger end.")

	var nt = model.NodeTrigger{
		NodeId:    nodeId,
		TriggerId: triggerId,
		CreatedBy: createdBy,
	}

	if res := db.Create(&nt); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"node_id":    nodeId,
			"trigger_id": triggerId,
			"created_at": nt.CreatedAt,
			"created_by": createdBy,
			"error":      res.Error,
		}, "An error occurred while db.Create.")
		return false
	}

	return true
}

// RemoveNodeTrigger 关联表操作::移除节点中触发器
func RemoveNodeTrigger(nodeId, triggerId int) bool {
	log.Info("dao.RemoveNodeTrigger called.")
	defer log.Info("dao.RemoveNodeTrigger end.")

	res := db.Delete(model.NodeTrigger{}, "node_id=? AND trigger_id=?", nodeId, triggerId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"node_id":    nodeId,
			"trigger_id": triggerId,
			"error":      res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// GetNodeTriggerCount 关联表操作::获得节点中所有触发器数量
func GetNodeTriggerCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Trigger{}).
		Select("triggers.id AS id, " +
			"triggers.name AS name, " +
			"triggers.description AS description, " +
			"triggers.arguments AS arguments").
		Joins("LEFT JOIN node_triggers AS nt ON triggers.id = nt.trigger_id")

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Count.")
		return count, false
	}
	return count, true
}

// ListNodeTriggers 关联表操作::列出节点中所有触发器
func ListNodeTriggers(nodeId int) ([]model.NodesTrigger, bool) {
	log.Info("dao.ListNodeTriggers called.")
	defer log.Info("dao.ListNodeTriggers end.")

	var d = db.Model(&model.Trigger{})
	nts := make([]model.NodesTrigger, 0)

	res := d.Select("triggers.id AS id, "+
		"triggers.name AS name, "+
		"triggers.description AS description, "+
		"triggers.arguments AS arguments").
		Joins("LEFT JOIN node_triggers AS nt ON triggers.id = nt.trigger_id").
		Where("nt.node_id=?", nodeId).
		Find(&nts)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found")
			return nts, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return nts, true
}

package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/model"
)

// NewNode 新增节点
func (d *Dao) NewNode(
	ctx context.Context,
	name string, parentId *uint32, category int32,
	entryCondition, assigneeCondition *string,
	vFields, eFields, createdBy string,
) (*model.Node, error) {
	node := &model.Node{
		Name:              name,
		ParentId:          parentId,
		Category:          category,
		EntryCondition:    entryCondition,
		AssigneeCondition: assigneeCondition,
		VisibleFields:     vFields,
		EditableFields:    eFields,
		CreatedBy:         createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&node)
	return node, res.Error
}

// RemoveNode 删除节点
func (d *Dao) RemoveNode(ctx context.Context, nodeId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Node{}, "id=?", nodeId)
	return res.Error
}

// SetNode 更新节点
func (d *Dao) SetNode(
	ctx context.Context,
	id uint32,
	name string, parentId *uint32, category int32,
	entryCondition, assigneeCondition *string,
	vFields, eFields, updatedBy string,
) (node *model.Node, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Node{}).
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
		First(&node)

	return node, res.Error
}

// GetNode 查询单个节点
func (d *Dao) GetNode(ctx context.Context, q orm.Query) (node *model.Node, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&node)
	return node, res.Error
}

// GetNodeCount 查询节点数量
func (d *Dao) GetNodeCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Node{})).Count(&count)
	return count, res.Error
}

// IsNodeExist 查询节点是否存在
func (d *Dao) IsNodeExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetNodeCount(ctx, q)
	return count > 0, err
}

// ListNodes 查询节点
func (d *Dao) ListNodes(ctx context.Context, q orm.Query) (nodes []*model.Node, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&nodes)
	return nodes, res.Error
}

// PagedListNodes 查询节点-分页
func (d *Dao) PagedListNodes(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	db := q.Where(d.getDbWithCtx(ctx)).Model(&model.Node{}).
		Select("nodes.id, " +
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
			"nodes.updated_by").
		Joins("LEFT JOIN nodes AS p ON p.id = nodes.parent_id")

	nodes := make([]*model.ListNodes, pageSize)

	return orm.PagingQuery(db, page, pageSize, &nodes, orderBy...)
}

// GetNodeChain 递推列出节点链
func (d *Dao) GetNodeChain(ctx context.Context, parentNode *model.NodeChain) error {
	node, err := d.GetNode(ctx, orm.Query{"parent_id": parentNode.Id})
	if err != nil {
		return err
	}

	if node == nil || node.Id < 1 {
		return nil
	}

	subNode := d.Node2Chain(ctx, node)
	parentNode.SubNode = subNode
	return d.GetNodeChain(ctx, subNode)
}

// Node2Chain 将节点转化为链结构
func (d *Dao) Node2Chain(ctx context.Context, n *model.Node) *model.NodeChain {
	ts, err := d.ListNodesTriggers(ctx, n.Id)

	nodeTris := make([]*model.NodeTriggers, 0)
	if err == nil && len(ts) > 0 {
		nodeTris = append(nodeTris, ts...)
	}

	return &model.NodeChain{
		Id:                n.Id,
		Name:              n.Name,
		Category:          n.Category,
		EntryCondition:    *n.EntryCondition,
		AssigneeCondition: *n.AssigneeCondition,
		Assignees:         nil,
		Triggers:          nodeTris,
		VisibleFields:     n.VisibleFields,
		EditableFields:    n.EditableFields,
		SubNode:           nil,
	}
}

// AddNodesTrigger 关联表操作::添加触发器至节点
func (d *Dao) AddNodesTrigger(ctx context.Context, nodeId, triggerId uint32, createdBy string) error {
	res := d.getDbWithCtx(ctx).Create(&model.NodeTrigger{
		NodeId:    nodeId,
		TriggerId: triggerId,
		CreatedBy: createdBy,
	})

	return res.Error
}

// RemoveNodesTrigger 关联表操作::移除指定节点中指定触发器
func (d *Dao) RemoveNodesTrigger(ctx context.Context, nodeId, triggerId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.NodeTriggers{}, "node_id=? AND trigger_id=?", nodeId, triggerId)
	return res.Error
}

// GetNodesTriggerCount 关联表操作::获得节点中所有触发器数量
func (d *Dao) GetNodesTriggerCount(ctx context.Context, q orm.Query) (count int64, err error) {
	_db := d.getDbWithCtx(ctx).Model(&model.Trigger{}).
		Select("triggers.id AS id, " +
			"triggers.name AS name, " +
			"triggers.description AS description, " +
			"triggers.arguments AS arguments").
		Joins("LEFT JOIN node_triggers AS nt ON triggers.id = nt.trigger_id")

	res := q.Where(_db).Count(&count)
	return count, res.Error
}

// ListNodesTriggers 关联表操作::列出节点中所有触发器
func (d *Dao) ListNodesTriggers(ctx context.Context, nodeId uint32) (nodeTris []*model.NodeTriggers, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Trigger{}).
		Select("triggers.id AS id, "+
			"triggers.name AS name, "+
			"triggers.description AS description, "+
			"triggers.arguments AS arguments").
		Joins("LEFT JOIN node_triggers AS nt ON triggers.id = nt.trigger_id").
		Where("nt.node_id=?", nodeId).
		Find(&nodeTris)

	return nodeTris, res.Error
}

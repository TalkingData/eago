package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/dto"
	"eago/flow/model"
)

// NewTrigger 创建触发器
func (d *Dao) NewTrigger(
	ctx context.Context, name, description, taskCodename, args, createdBy string,
) (*model.Trigger, error) {
	tri := &model.Trigger{
		Name:         name,
		Description:  &description,
		TaskCodename: taskCodename,
		Arguments:    args,
		CreatedBy:    createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&tri)
	return tri, res.Error
}

// RemoveTrigger 删除触发器
func (d *Dao) RemoveTrigger(ctx context.Context, triId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Trigger{}, "id=?", triId)
	return res.Error
}

// SetTrigger 更新触发器
func (d *Dao) SetTrigger(
	ctx context.Context, id uint32, name, description, taskCodename, args, updatedBy string,
) (tri *model.Trigger, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Trigger{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":          name,
			"description":   description,
			"task_codename": taskCodename,
			"arguments":     args,
			"updated_by":    updatedBy,
		}).
		Limit(1).Find(&tri)

	return tri, res.Error
}

// GetTrigger 查询单个触发器
func (d *Dao) GetTrigger(ctx context.Context, q orm.Query) (tri *model.Trigger, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&tri)
	return tri, res.Error
}

// GetTriggerCount 查询触发器数量
func (d *Dao) GetTriggerCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Trigger{})).Count(&count)
	return count, res.Error
}

// IsTriggerExist 查询触发器是否存在
func (d *Dao) IsTriggerExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetTriggerCount(ctx, q)
	return count > 0, err
}

// ListTriggers 查询触发器
func (d *Dao) ListTriggers(ctx context.Context, q orm.Query) (tris []*model.Trigger, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&tris)
	return tris, res.Error
}

// PagedListTriggers 查询触发器-分页
func (d *Dao) PagedListTriggers(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	tris := make([]*model.Trigger, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Trigger{}))
	return orm.PagingQuery(db, page, pageSize, &tris, orderBy...)
}

// ListTriggersNodes 关联表操作::列出触发器所关联节点
func (d *Dao) ListTriggersNodes(ctx context.Context, triID uint32) (tris []*dto.TriggersNode, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Node{}).
		Select("nodes.id AS id, "+
			"nodes.name AS name, "+
			"nodes.parent_id AS parent_id").
		Joins("LEFT JOIN node_triggers AS nt ON nodes.id = nt.node_id").
		Where("nt.trigger_id=?", triID).
		Find(&tris)

	return tris, res.Error
}

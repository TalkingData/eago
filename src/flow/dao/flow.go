package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/dto"
	"eago/flow/model"
)

// NewFlow 创建流程
func (d *Dao) NewFlow(
	ctx context.Context,
	name, instanceTitle string,
	catId *uint32,
	description string,
	disabled bool,
	frmID, firstNodeID uint32,
	createdBy string,
) (*model.Flow, error) {
	f := &model.Flow{
		Name:          name,
		InstanceTitle: instanceTitle,
		CategoriesId:  catId,
		Disabled:      &disabled,
		Description:   &description,
		FormId:        frmID,
		FirstNodeId:   firstNodeID,
		CreatedBy:     createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&f)
	return f, res.Error
}

// RemoveFlow 删除流程
func (d *Dao) RemoveFlow(ctx context.Context, flowId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Flow{}, "id=?", flowId)
	return res.Error
}

// SetFlow 更新流程
func (d *Dao) SetFlow(
	ctx context.Context,
	id uint32,
	name, instanceTitle string,
	catId *uint32,
	description string,
	disabled bool,
	frmID, firstNodeID uint32,
	updatedBy string,
) (f *model.Flow, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Flow{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":           name,
			"instance_title": instanceTitle,
			"categories_id":  catId,
			"disabled":       disabled,
			"description":    description,
			"form_id":        frmID,
			"first_node_id":  firstNodeID,
			"updated_by":     updatedBy,
		}).
		First(&f)

	return f, res.Error
}

// GetFlow 查询单个流程
func (d *Dao) GetFlow(ctx context.Context, q orm.Query) (f *model.Flow, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&f)
	return f, res.Error
}

// GetFlowCount 查询流程数量
func (d *Dao) GetFlowCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Flow{})).Count(&count)
	return count, res.Error
}

// IsFlowExist 查询流程是否存在
func (d *Dao) IsFlowExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetFlowCount(ctx, q)
	return count > 0, err
}

// ListFlows 查询流程
func (d *Dao) ListFlows(ctx context.Context, q orm.Query) (flows []*model.Flow, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&flows)
	return flows, res.Error
}

// PagedListFlows 查询流程-分页
func (d *Dao) PagedListFlows(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	db := q.Where(d.getDbWithCtx(ctx)).Model(&model.Flow{}).
		Select("flows.id AS id, " +
			"flows.name AS name, " +
			"flows.instance_title AS instance_title, " +
			"flows.categories_id AS categories_id, " +
			"c.name AS categories_name, " +
			"flows.disabled AS disabled, " +
			"flows.description AS description, " +
			"flows.form_id AS form_id, " +
			"f.name AS form_name, " +
			"flows.first_node_id AS first_node_id, " +
			"n.name AS first_node_name, " +
			"flows.created_at AS created_at, " +
			"flows.created_by AS created_by, " +
			"flows.updated_at AS updated_at, " +
			"flows.updated_by AS updated_by").
		Joins("LEFT JOIN categories AS c ON c.id = flows.categories_id").
		Joins("LEFT JOIN forms AS f ON f.id = flows.form_id").
		Joins("LEFT JOIN nodes AS n ON n.id = flows.first_node_id")

	flows := make([]*dto.ListFlows, pageSize)

	return orm.PagingQuery(db, page, pageSize, &flows, orderBy...)
}

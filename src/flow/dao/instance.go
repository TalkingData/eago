package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/model"
)

// NewInstance 创建流程实例
func (d *Dao) NewInstance(
	ctx context.Context,
	formId uint32, status int32,
	name, formData, flowChain, createdBy string,
) (*model.Instance, error) {
	// 保证流程实例名称不超过表最大长度
	if len(name) > d.conf.Const.InstanceNameMaxLength {
		name = name[:d.conf.Const.InstanceNameMaxLength]
	}

	i := &model.Instance{
		Name:      name,
		Status:    status,
		FormId:    formId,
		FormData:  &formData,
		FlowChain: &flowChain,
		CreatedBy: createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&i)
	return i, res.Error
}

// SetInstance 设置流程实例
func (d *Dao) SetInstance(
	ctx context.Context,
	id uint32, status, currStep, assigneesReq int32,
	flowChain, currAssignees, passedAssignees, updatedBy string,
) (ins *model.Instance, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Instance{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status":             status,
			"current_step":       currStep,
			"flow_chain":         flowChain,
			"assignees_required": assigneesReq,
			"current_assignees":  currAssignees,
			"passed_assignees":   passedAssignees,
			"updated_by":         updatedBy,
		}).
		Limit(1).Find(&ins)

	return ins, res.Error
}

// SetHandleInstance 设置流程实例
func (d *Dao) SetHandleInstance(
	ctx context.Context,
	id uint32, status, step, assigneesReq int32,
	formData, currAssignees, passedAssignees, updatedBy string,
) error {
	res := d.getDbWithCtx(ctx).Model(&model.Instance{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status":             status,
			"current_step":       step,
			"form_data":          formData,
			"assignees_required": assigneesReq,
			"current_assignees":  currAssignees,
			"passed_assignees":   passedAssignees,
			"updated_by":         updatedBy,
		})

	return res.Error
}

// GetInstance 查询单个流程实例
func (d *Dao) GetInstance(ctx context.Context, q orm.Query) (inst *model.Instance, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&inst)
	return inst, res.Error
}

// GetInstanceCount 查询流程实例数量
func (d *Dao) GetInstanceCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Instance{})).Count(&count)
	return count, res.Error
}

// IsInstanceExist 查询流程实例是否存在
func (d *Dao) IsInstanceExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetInstanceCount(ctx, q)
	return count > 0, err
}

// PagedListInstances 查询流程实例-分页
func (d *Dao) PagedListInstances(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	insts := make([]*model.Instance, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Instance{}))
	return orm.PagingQuery(db, page, pageSize, &insts, orderBy...)
}

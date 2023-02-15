package dao

import (
	"context"
	"eago/common/orm"
	"eago/task/model"
)

// NewSchedule 新建计划任务
func (d *Dao) NewSchedule(
	ctx context.Context, tCodeName, expr, args, description string, timeout int64, disabled bool, createdBy string,
) (*model.Schedule, error) {
	sch := &model.Schedule{
		TaskCodename: tCodeName,
		Expression:   expr,
		Description:  &description,
		Timeout:      &timeout,
		Arguments:    args,
		Disabled:     &disabled,
		CreatedBy:    createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&sch)
	return sch, res.Error
}

// RemoveSchedule 删除计划任务
func (d *Dao) RemoveSchedule(ctx context.Context, schId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Schedule{}, "id=?", schId)
	return res.Error
}

// SetSchedule 更新计划任务
func (d *Dao) SetSchedule(
	ctx context.Context,
	id uint32, tCodeName, expr, args, description string, timeout int64, disabled bool, updatedBy string,
) (sch *model.Schedule, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Schedule{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"task_codename": tCodeName,
			"expression":    expr,
			"timeout":       timeout,
			"arguments":     args,
			"disabled":      disabled,
			"description":   description,
			"updated_by":    updatedBy,
		}).
		Limit(1).Find(&sch)
	return sch, res.Error
}

// GetSchedule 查询单个计划任务
func (d *Dao) GetSchedule(ctx context.Context, q orm.Query) (sch *model.Task, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&sch)
	return sch, res.Error
}

// GetScheduleCount 查询计划任务数量
func (d *Dao) GetScheduleCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Schedule{})).Count(&count)
	return count, res.Error
}

// IsScheduleExist 查询计划任务是否存在
func (d *Dao) IsScheduleExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetScheduleCount(ctx, q)
	return count > 0, err
}

// ListSchedules 查询计划任务
func (d *Dao) ListSchedules(ctx context.Context, q orm.Query) (schs []*model.Schedule, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&schs)
	return schs, res.Error
}

// PagedListSchedules 查询计划任务-分页
func (d *Dao) PagedListSchedules(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	schs := make([]*model.Schedule, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Schedule{}))
	return orm.PagingQuery(db, page, pageSize, &schs, orderBy...)
}

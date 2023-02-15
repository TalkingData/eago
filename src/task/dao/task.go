package dao

import (
	"context"
	"eago/common/orm"
	"eago/task/model"
)

// NewTask 新建任务
func (d *Dao) NewTask(
	ctx context.Context, disabled bool, category int32, codename, description, fParams, createdBy string,
) (*model.Task, error) {
	t := &model.Task{
		Category:     &category,
		Codename:     codename,
		Description:  &description,
		FormalParams: fParams,
		Disabled:     &disabled,
		CreatedBy:    createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&t)
	return t, res.Error
}

// RemoveTask 删除任务
func (d *Dao) RemoveTask(ctx context.Context, taskId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Task{}, "id=?", taskId)
	return res.Error
}

// SetTask 更新任务
func (d *Dao) SetTask(
	ctx context.Context, id uint32, disabled bool, category int32, codename, description, fParams, updatedBy string,
) (t *model.Task, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Task{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"category":      category,
			"codename":      codename,
			"formal_params": fParams,
			"disabled":      disabled,
			"description":   description,
			"updated_by":    updatedBy,
		}).
		Limit(1).Find(&t)
	return t, res.Error
}

// GetTask 查询单个任务
func (d *Dao) GetTask(ctx context.Context, q orm.Query) (t *model.Task, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&t)
	return t, res.Error
}

// GetTaskCount 查询任务数量
func (d *Dao) GetTaskCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Task{})).Count(&count)
	return count, res.Error
}

// IsTaskExist 查询任务是否存在
func (d *Dao) IsTaskExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetTaskCount(ctx, q)
	return count > 0, err
}

// ListTasks 查询任务
func (d *Dao) ListTasks(ctx context.Context, q orm.Query) (ts []*model.Task, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&ts)
	return ts, res.Error
}

// PagedListTasks 查询任务-分页
func (d *Dao) PagedListTasks(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	ts := make([]*model.Task, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Task{}))
	return orm.PagingQuery(db, page, pageSize, &ts, orderBy...)
}

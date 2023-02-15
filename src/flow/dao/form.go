package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/model"
)

// NewForm 创建表单
func (d *Dao) NewForm(
	ctx context.Context, name string, disabled bool, description, body, createdBy string,
) (*model.Form, error) {
	f := &model.Form{
		Name:        name,
		Disabled:    &disabled,
		Description: &description,
		Body:        &body,
		CreatedBy:   createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&f)
	return f, res.Error
}

// RemoveForm 删除表单
func (d *Dao) RemoveForm(ctx context.Context, frmId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Form{}, "id=?", frmId)
	return res.Error
}

// SetForm 更新表单
func (d *Dao) SetForm(
	ctx context.Context, id uint32, name string, disabled bool, description, updatedBy string,
) (f *model.Form, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Form{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"disabled":    disabled,
			"description": description,
			"updated_by":  updatedBy,
		}).
		Limit(1).Find(&f)
	return f, res.Error
}

// GetForm 查询单个表单
func (d *Dao) GetForm(ctx context.Context, q orm.Query) (f *model.Form, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&f)
	return f, res.Error
}

// GetFormCount 查询表单数量
func (d *Dao) GetFormCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Form{})).Count(&count)
	return count, res.Error
}

// IsFormExist 查询表单是否存在
func (d *Dao) IsFormExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetFormCount(ctx, q)
	return count > 0, err
}

// ListForms 查询表单
func (d *Dao) ListForms(ctx context.Context, q orm.Query) (fs []*model.Form, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&fs)
	return fs, res.Error
}

// PagedListForms 查询表单-分页
func (d *Dao) PagedListForms(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	fs := make([]*model.Form, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Form{}))
	return orm.PagingQuery(db, page, pageSize, &fs, orderBy...)
}

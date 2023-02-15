package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/model"
)

// NewCategory 创建类别
func (d *Dao) NewCategory(ctx context.Context, name, createdBy string) (*model.Categories, error) {
	cat := &model.Categories{
		Name:      name,
		CreatedBy: createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&cat)
	return cat, res.Error
}

// RemoveCategory 删除类别
func (d *Dao) RemoveCategory(ctx context.Context, catId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Categories{}, "id=?", catId)
	return res.Error
}

// SetCategory 更新类别
func (d *Dao) SetCategory(ctx context.Context, id uint32, name, updatedBy string) (cat *model.Categories, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Categories{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":       name,
			"updated_by": updatedBy,
		}).
		Limit(1).Find(&cat)

	return cat, res.Error
}

// GetCategory 查询单个类别
func (d *Dao) GetCategory(ctx context.Context, q orm.Query) (cat *model.Categories, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&cat)
	return cat, res.Error
}

// GetCategoriesCount 查询类别数量
func (d *Dao) GetCategoriesCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Categories{})).Count(&count)
	return count, res.Error
}

// IsCategoryExist 查询类别是否存在
func (d *Dao) IsCategoryExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetCategoriesCount(ctx, q)
	return count > 0, err
}

// ListCategories 查询类别
func (d *Dao) ListCategories(ctx context.Context, q orm.Query) (cats []*model.Categories, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&cats)
	return cats, res.Error
}

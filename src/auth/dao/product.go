package dao

import (
	"context"
	"eago/auth/model"
	"eago/common/orm"
)

// NewProduct 新建产品线
func (d *Dao) NewProduct(
	ctx context.Context, name, alias, description string, disabled *bool,
) (*model.Product, error) {
	prod := &model.Product{
		Name:        name,
		Alias:       alias,
		Disabled:    disabled,
		Description: &description,
	}
	res := d.getDbWithCtx(ctx).Create(&prod)
	return prod, res.Error
}

// RemoveProduct 删除产品线
func (d *Dao) RemoveProduct(ctx context.Context, productId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Product{}, "id=?", productId)
	return res.Error
}

// SetProduct 更新产品线
func (d *Dao) SetProduct(
	ctx context.Context, id uint32, name, alias, description string, disabled bool,
) (prod *model.Product, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Product{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"alias":       alias,
			"disabled":    disabled,
			"description": description,
		}).
		First(&prod)
	return prod, res.Error
}

// GetProduct 查询单个产品线
func (d *Dao) GetProduct(ctx context.Context, q orm.Query) (prod *model.Product, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&prod)
	return prod, res.Error
}

// GetProductCount 查询产品线数量
func (d *Dao) GetProductCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Product{})).Count(&count)
	return count, res.Error
}

// IsProductExist 查询产品线是否存在
func (d *Dao) IsProductExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetProductCount(ctx, q)
	return count > 0, err
}

// ListProducts 查询产品线
func (d *Dao) ListProducts(ctx context.Context, q orm.Query) (prods []*model.Product, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&prods)
	return prods, res.Error
}

// PagedListProducts 分页查询产品线
func (d *Dao) PagedListProducts(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	prods := make([]*model.Product, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Product{}))
	return orm.PagingQuery(db, page, pageSize, &prods, orderBy...)
}

// AddUser2Product 关联表操作::添加指定用户至指定产品线
func (d *Dao) AddUser2Product(ctx context.Context, productId, userId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Create(&model.UserProduct{
		ProductId: productId,
		UserId:    userId,
		IsOwner:   &isOwner,
	})

	return res.Error
}

// RemoveProductsUser 关联表操作::移除指定产品线中指定用户
func (d *Dao) RemoveProductsUser(ctx context.Context, productId, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserProduct{}, "user_id=? AND product_id=?", userId, productId)
	return res.Error
}

// RemoveUsersProducts 关联表操作::移除指定用户所有产品线
func (d *Dao) RemoveUsersProducts(ctx context.Context, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserProduct{}, "user_id=?", userId)
	return res.Error
}

// SetProductsOwner 关联表操作::设置指定用户是否是指定产品线的Owner
func (d *Dao) SetProductsOwner(ctx context.Context, productId, userId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Model(&model.UserProduct{}).
		Where("user_id=? AND product_id=?", userId, productId).
		Update("is_owner", isOwner)
	return res.Error
}

// GetProductsUserCount 关联表操作::列出指定产品线中用户数量
func (d *Dao) GetProductsUserCount(ctx context.Context, q orm.Query) (count int64, err error) {
	_db := d.getDbWithCtx(ctx).Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"up.is_owner AS is_owner, " +
			"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON users.id = up.user_id")

	res := q.Where(_db).Count(&count)
	return count, res.Error
}

// IsEmptyProduct 关联表操作::指定产品线是否为空，不包含用户则为空
func (d *Dao) IsEmptyProduct(ctx context.Context, prodId uint32) (bool, error) {
	count, err := d.GetProductsUserCount(ctx, orm.Query{"product_id=?": prodId})
	return count == 0, err
}

// ListProductsUsers 关联表操作::列出指定产品线中所有用户
func (d *Dao) ListProductsUsers(
	ctx context.Context, productId uint32, q orm.Query,
) (mbrUser []*model.MemberUser, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.User{})).
		Select("users.id AS id, "+
			"users.username AS username, "+
			"up.is_owner AS is_owner, "+
			"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON users.id = up.user_id").
		Where("product_id=?", productId).
		Find(&mbrUser)
	return mbrUser, res.Error
}

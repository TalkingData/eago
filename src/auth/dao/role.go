package dao

import (
	"context"
	"eago/auth/model"
	"eago/common/orm"
)

// NewRole 新建角色
func (d *Dao) NewRole(ctx context.Context, name, description string) (*model.Role, error) {
	r := &model.Role{
		Name:        name,
		Description: &description,
	}
	res := d.getDbWithCtx(ctx).Create(&r)
	return r, res.Error
}

// RemoveRole 删除角色
func (d *Dao) RemoveRole(ctx context.Context, roleId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Role{}, "id=?", roleId)
	return res.Error
}

// SetRole 更新角色
func (d *Dao) SetRole(ctx context.Context, id uint32, name, description string) (r *model.Role, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Role{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).
		Limit(1).Find(&r)
	return r, res.Error
}

// GetRole 查询单个角色
func (d *Dao) GetRole(ctx context.Context, q orm.Query) (r *model.Role, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&r)
	return r, res.Error
}

// GetRoleCount 查询角色数量
func (d *Dao) GetRoleCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Role{})).Count(&count)
	return count, res.Error
}

// IsRoleExist 查询角色是否存在
func (d *Dao) IsRoleExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetRoleCount(ctx, q)
	return count > 0, err
}

// ListRoles 查询角色
func (d *Dao) ListRoles(ctx context.Context, q orm.Query) (rs []*model.Role, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&rs)
	return rs, res.Error
}

// PagedListRoles 分页查询角色
func (d *Dao) PagedListRoles(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	rs := make([]*model.Role, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Role{}))
	return orm.PagingQuery(db, page, pageSize, &rs, orderBy...)
}

// AddUser2Role 关联表操作::添加用户至指定角色
func (d *Dao) AddUser2Role(ctx context.Context, roleId, userId uint32) error {
	res := d.getDbWithCtx(ctx).Create(&model.UserRole{
		RoleId: roleId,
		UserId: userId,
	})
	return res.Error
}

// RemoveRolesUser 关联表操作::移除指定角色中指定用户
func (d *Dao) RemoveRolesUser(ctx context.Context, roleId, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserRole{}, "role_id=? AND user_id=?", roleId, userId)
	return res.Error
}

// RemoveUsersRoles 关联表操作::移除指定用户所有角色
func (d *Dao) RemoveUsersRoles(ctx context.Context, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserRole{}, "user_id=?", userId)
	return res.Error
}

// GetRolesUserCount 关联表操作::获得指定角色中用户数量
func (d *Dao) GetRolesUserCount(ctx context.Context, q orm.Query) (count int64, err error) {
	_db := d.getDbWithCtx(ctx).Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"ur.joined_at AS joined_at").
		Joins("LEFT JOIN user_roles AS ur ON users.id = ur.user_id")

	res := q.Where(_db).Count(&count)
	return count, res.Error
}

// IsEmptyRole 关联表操作::指定角色是否为空，不包含用户则为空
func (d *Dao) IsEmptyRole(ctx context.Context, rId uint32) (bool, error) {
	count, err := d.GetRolesUserCount(ctx, orm.Query{"role_id=?": rId})
	return count == 0, err
}

// ListRolesUsers 关联表操作::列出指定角色中用户
func (d *Dao) ListRolesUsers(ctx context.Context, roleId uint32) (rUsers []*model.RolesUser, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.User{}).
		Select("users.id AS id, "+
			"users.username AS username, "+
			"ur.joined_at AS joined_at").
		Joins("LEFT JOIN user_roles AS ur ON users.id = ur.user_id").
		Where("role_id=?", roleId).
		Find(&rUsers)
	return rUsers, res.Error
}

package dao

import (
	"context"
	"eago/auth/model"
	"eago/common/orm"
)

// NewGroup 新建组
func (d *Dao) NewGroup(ctx context.Context, name, description string) (*model.Group, error) {
	g := &model.Group{
		Name:        name,
		Description: &description,
	}
	res := d.getDbWithCtx(ctx).Create(&g)
	return g, res.Error
}

// RemoveGroup 删除组
func (d *Dao) RemoveGroup(ctx context.Context, groupId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Group{}, "id=?", groupId)
	return res.Error
}

// SetGroup 更新组
func (d *Dao) SetGroup(ctx context.Context, id uint32, name, description string) (g *model.Group, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Group{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).
		Limit(1).Find(&g)
	return g, res.Error
}

// GetGroup 查询单个组
func (d *Dao) GetGroup(ctx context.Context, q orm.Query) (g *model.Group, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&g)
	return g, res.Error
}

// GetGroupCount 查询组数量
func (d *Dao) GetGroupCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Group{})).Count(&count)
	return count, res.Error
}

// IsGroupExist 查询组是否存在
func (d *Dao) IsGroupExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetGroupCount(ctx, q)
	return count > 0, err
}

// ListGroups 查询组
func (d *Dao) ListGroups(ctx context.Context, q orm.Query) (gs []*model.Group, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&gs)
	return gs, res.Error
}

// PagedListGroups 分页查询组
func (d *Dao) PagedListGroups(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	gs := make([]*model.Group, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Group{}))
	return orm.PagingQuery(db, page, pageSize, &gs, orderBy...)
}

// AddUser2Group 关联表操作::添加指定用户至指定组
func (d *Dao) AddUser2Group(ctx context.Context, groupId, userId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Create(&model.UserGroup{
		GroupId: groupId,
		UserId:  userId,
		IsOwner: &isOwner,
	})
	return res.Error
}

// RemoveGroupsUser 关联表操作::移除指定组中指定用户
func (d *Dao) RemoveGroupsUser(ctx context.Context, groupId, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserGroup{}, "group_id=? AND user_id=?", groupId, userId)
	return res.Error
}

// RemoveUsersGroups 关联表操作::移除指定用户所有组
func (d *Dao) RemoveUsersGroups(ctx context.Context, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserGroup{}, "user_id=?", userId)
	return res.Error
}

// SetGroupsOwner 关联表操作::设置用户是否是指定组的Owner
func (d *Dao) SetGroupsOwner(ctx context.Context, groupId, userId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Model(&model.UserGroup{}).
		Where("group_id=? AND user_id=?", groupId, userId).
		Update("is_owner", isOwner)
	return res.Error
}

// GetGroupsUserCount 关联表操作::列出指定组中所有用户数量
func (d *Dao) GetGroupsUserCount(ctx context.Context, q orm.Query) (count int64, err error) {
	_db := d.getDbWithCtx(ctx).Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"ug.is_owner AS is_owner, " +
			"ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON users.id = ug.user_id")

	res := q.Where(_db).Count(&count)
	return count, res.Error
}

// IsEmptyGroup 关联表操作::查询指定组是否为空，不包含用户则为空
func (d *Dao) IsEmptyGroup(ctx context.Context, gId uint32) (bool, error) {
	count, err := d.GetGroupsUserCount(ctx, orm.Query{"group_id=?": gId})
	return count == 0, err
}

// ListGroupsUsers 关联表操作::列出指定组中所有用户
func (d *Dao) ListGroupsUsers(
	ctx context.Context, groupId uint32, q orm.Query,
) (mbrUser []*model.MemberUser, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.User{})).
		Select("users.id AS id, "+
			"users.username AS username, "+
			"ug.is_owner AS is_owner, "+
			"ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON users.id = ug.user_id").
		Where("group_id=?", groupId).
		Find(&mbrUser)
	return mbrUser, res.Error
}

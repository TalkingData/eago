package dao

import (
	"context"
	"eago/auth/model"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/utils"
	"strings"
	"time"
)

// NewUser 新建用户
func (d *Dao) NewUser(ctx context.Context, username, email string, login bool) (*model.User, error) {
	u := &model.User{
		Username: strings.ToLower(username),
		Email:    strings.ToLower(email),
	}

	// 判断时候设置最近登录时间
	if login {
		u.LastLogin = &utils.CustomTime{Time: time.Now()}
	}

	res := d.getDbWithCtx(ctx).Create(&u)
	return u, res.Error
}

// RemoveUser 删除用户
func (d *Dao) RemoveUser(ctx context.Context, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.User{}, "id=?", userId)
	return res.Error
}

// SetUserLastLogin 更新用户最后登录时间，自动更新为当前时间
func (d *Dao) SetUserLastLogin(ctx context.Context, id uint32) error {
	db := d.getDbWithCtx(ctx).Model(&model.User{}).Where("id=?", id)

	res := db.Update("last_login", &utils.CustomTime{Time: time.Now()})
	return res.Error
}

// DisableUser 更新用户为禁用状态
func (d *Dao) DisableUser(ctx context.Context, id uint32) error {
	res := d.getDbWithCtx(ctx).Model(&model.User{}).
		Where("id=?", id).
		Update("disabled", true)
	return res.Error
}

// SetUser 更新用户
func (d *Dao) SetUser(ctx context.Context, id uint32, email, phone string) (u *model.User, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.User{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"email": strings.ToLower(email),
			"phone": phone,
		}).
		First(&u)
	return u, res.Error
}

// GetUser 查询单个用户
func (d *Dao) GetUser(ctx context.Context, q orm.Query) (u *model.User, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&u)
	return u, res.Error
}

// GetUserCount 查询用户数量
func (d *Dao) GetUserCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.User{})).Count(&count)
	return count, res.Error
}

// IsUserExist 查询用户是否存在
func (d *Dao) IsUserExist(ctx context.Context, q orm.Query) (exist bool, err error) {
	count, err := d.GetUserCount(ctx, q)
	return count > 0, err
}

// ListUsers 查询用户
func (d *Dao) ListUsers(ctx context.Context, q orm.Query) (users []*model.User, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&users)
	return users, res.Error
}

// PagedListUsers 分页查询用户
func (d *Dao) PagedListUsers(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	us := make([]*model.User, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.User{}))
	return orm.PagingQuery(db, page, pageSize, &us, orderBy...)
}

// IsSuperUser 查询指定用户是否是Admin
func (d *Dao) IsSuperUser(ctx context.Context, userId uint32) bool {
	u, err := d.GetUser(ctx, orm.Query{"id=?": userId})
	if err != nil {
		d.lg.WarnWithFields(logger.Fields{
			"user_id": userId,
			"error":   err,
		}, "An error occurred while dao.GetUser in dao.IsSuperUser, skipped it.")
		return false
	}
	if u == nil {
		d.lg.WarnWithFields(logger.Fields{
			"user_id": userId,
		}, "Got nil user object in dao.IsSuperUser, skipped it.")
		return false
	}

	return u.IsSuperuser
}

// ListUsersRoles 关联表操作::列出指定用户所有角色
func (d *Dao) ListUsersRoles(ctx context.Context, userId uint32) (roles []*model.Role, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Role{}).
		Joins("LEFT JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_id=?", userId).
		Find(&roles)
	return roles, res.Error
}

// ListUsersProducts 关联表操作::列出指定用户所有产品线
func (d *Dao) ListUsersProducts(ctx context.Context, userId uint32) (upms []*model.UserProductMember, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Product{}).
		Select("products.id AS id, "+
			"products.name AS name, "+
			"products.alias AS alias, "+
			"products.disabled AS disabled, "+
			"up.is_owner AS is_owner, "+
			"up.joined_at AS joined_at").
		Joins("LEFT JOIN user_products AS up ON products.id = up.product_id").
		Where("user_id=?", userId).
		Find(&upms)
	return upms, res.Error
}

// ListUsersGroups 关联表操作::列出指定用户所有组
func (d *Dao) ListUsersGroups(ctx context.Context, userId uint32) (ums []*model.UserGroupMember, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Group{}).
		Select("groups.id AS id, "+
			"groups.name AS name, "+
			"ug.is_owner AS is_owner, "+
			"ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON groups.id = ug.group_id").
		Where("user_id=?", userId).
		Find(&ums)
	return ums, res.Error
}

// GetUsersDepartment 关联表操作::获得指定用户所在部门
func (d *Dao) GetUsersDepartment(ctx context.Context, userId uint32) (udm *model.UserDepartmentMember, err error) {
	res := d.getDbWithCtx(ctx).Model(&model.Department{}).
		Select("departments.id AS id, "+
			"departments.name AS name, "+
			"departments.parent_id AS parent_id, "+
			"ud.is_owner AS is_owner, "+
			"ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON departments.id = ud.department_id").
		Where("user_id=?", userId).
		Limit(1).Find(&udm)
	return udm, res.Error
}

// GetUsersDepartmentChain 关联表操作::列出指定用户所在部门链，包含所有层级情况
func (d *Dao) GetUsersDepartmentChain(ctx context.Context, userId uint32) (res []*model.UserDepartmentNode) {
	// 获取用户所在部门
	userDeptObj, err := d.GetUsersDepartment(ctx, userId)
	if err != nil {
		d.lg.WarnWithFields(logger.Fields{
			"user_id": userId,
			"error":   err,
		}, "An error occurred while dao.GetUsersDepartment in dao.GetUsersDepartmentChain, skipped it.")
		return
	}

	// 如果找不到用户所在部门，直接返回
	if userDeptObj == nil {
		return
	}

	// 添加部门
	res = append(res, &model.UserDepartmentNode{
		Id:       userDeptObj.Id,
		Name:     userDeptObj.Name,
		ParentId: userDeptObj.ParentId,
	})
	// 递归查找上级部门
	d.deptSimpleRecursion(ctx, userDeptObj.ParentId, res)

	// 反转数组，使父部门在前
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return
}

// MakeUserHandover 用户交接
func (d *Dao) MakeUserHandover(
	ctx context.Context, userId, tgtUserId uint32,
) (srcUser *model.User, tgtUser *model.User, err error) {
	tx := d.getDbWithCtx(ctx).Begin()
	defer tx.Rollback()

	// 获得交接用户
	if res := tx.Where("id=?", userId).Find(&srcUser); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to find user.")
		return nil, nil, res.Error
	}

	// 获得交接目标用户
	if res := tx.Where("id=?", tgtUserId).Find(&tgtUser); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"target_user_id": tgtUserId,
			"error":          res.Error,
		}, "Failed to find handover target user.")
		return nil, nil, res.Error
	}

	// 交接产品线Owner权限
	res := tx.Model(&model.UserProduct{}).
		Where("user_id=? AND is_owner=?", userId, true).
		Updates(map[string]interface{}{"user_id": tgtUserId, "joined_at": time.Now()})
	if res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id":        userId,
			"target_user_id": tgtUserId,
			"error":          res.Error,
		}, "Failed to handover user's products.")
		return nil, nil, res.Error
	}

	// 交接组Owner权限
	res = tx.Model(&model.UserGroup{}).
		Where("user_id=? AND is_owner=?", userId, true).
		Updates(map[string]interface{}{"user_id": tgtUserId, "joined_at": time.Now()})
	if res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id":        userId,
			"target_user_id": tgtUserId,
			"error":          res.Error,
		}, "Failed to handover user's groups.")
		return nil, nil, res.Error
	}

	// 交接角色权限
	res = tx.Model(&model.UserRole{}).
		Where("user_id=?", userId).
		Updates(map[string]interface{}{"user_id": tgtUserId, "joined_at": time.Now()})
	if res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id":        userId,
			"target_user_id": tgtUserId,
			"error":          res.Error,
		}, "Failed to handover user's roles.")
		return nil, nil, res.Error
	}

	// 删除所在产品线
	if res = tx.Where("user_id=?", userId).Delete(model.UserProduct{}); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to remove user products.")
		return nil, nil, res.Error
	}

	// 删除所在部门
	if res = tx.Where("user_id=?", userId).Delete(model.UserDepartment{}); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to remove user department.")
		return nil, nil, res.Error
	}

	// 删除所在组
	if res = tx.Where("user_id=?", userId).Delete(model.UserGroup{}); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to remove user groups.")
		return nil, nil, res.Error
	}

	// 删除所在角色
	if res = tx.Where("user_id=?", userId).Delete(model.UserRole{}); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to remove user roles.")
		return nil, nil, res.Error
	}

	// 删除用户
	if res = tx.Where("id=?", userId).Delete(model.User{}); res.Error != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "Failed to remove user.")
		return nil, nil, res.Error
	}

	tx.Commit()
	return
}

// deptSimpleRecursion 递归查找部门
func (d *Dao) deptSimpleRecursion(ctx context.Context, pId *uint32, array []*model.UserDepartmentNode) {
	if pId == nil {
		return
	}

	pDept, err := d.GetDepartment(ctx, orm.Query{"id": pId})
	if err != nil {
		d.lg.WarnWithFields(logger.Fields{
			"department_id": pId,
			"error":         err,
		}, "An error occurred while dao.GetDepartment in dao.deptSimpleRecursion, skipped it.")
		return
	}
	// 如果找不到用户所在部门，直接返回
	if pDept == nil {
		return
	}

	// 添加部门
	array = append(array, &model.UserDepartmentNode{
		Id:       pDept.Id,
		Name:     pDept.Name,
		ParentId: pDept.ParentId,
	})
	// 递归查找上级部门
	d.deptSimpleRecursion(ctx, pDept.ParentId, array)
}

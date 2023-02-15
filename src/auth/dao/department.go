package dao

import (
	"context"
	"eago/auth/dto"
	"eago/auth/model"
	"eago/common/logger"
	"eago/common/orm"
	"gorm.io/gorm"
)

// NewDepartment 新建部门
func (d *Dao) NewDepartment(ctx context.Context, name string, parentId *uint32) (*model.Department, error) {
	dept := &model.Department{
		Name:     name,
		ParentId: parentId,
	}
	res := d.getDbWithCtx(ctx).Create(&dept)
	return dept, res.Error
}

// RemoveDepartment 删除部门
func (d *Dao) RemoveDepartment(ctx context.Context, deptId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.Department{}, "id=?", deptId)
	return res.Error
}

// EmptyDepartment 清空部门
func (d *Dao) EmptyDepartment(ctx context.Context) error {
	res := d.getDbWithCtx(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model.UserDepartment{})
	if res.Error != nil {
		d.lg.WarnWithFields(logger.Fields{
			"error": res.Error,
		}, "An error occurred while db.Delete(model.UserDepartment{}) in dao.EmptyDepartment.")
		return res.Error
	}
	res = d.getDbWithCtx(ctx).Model(&model.UserDepartment{}).Exec("ALTER TABLE user_departments AUTO_INCREMENT=1")
	if res.Error != nil {
		d.lg.WarnWithFields(logger.Fields{
			"error": res.Error,
		}, "An error occurred while 'ALTER TABLE user_departments AUTO_INCREMENT=1' in dao.EmptyDepartment.")
		return res.Error
	}

	res = d.getDbWithCtx(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model.Department{})
	if res.Error != nil {
		d.lg.WarnWithFields(logger.Fields{
			"error": res.Error,
		}, "An error occurred while db.Delete(model.Department{}) in dao.EmptyDepartment.")
		return res.Error
	}
	res = d.getDbWithCtx(ctx).Model(&model.Department{}).Exec("ALTER TABLE departments AUTO_INCREMENT=1")
	if res.Error != nil {
		d.lg.WarnWithFields(logger.Fields{
			"error": res.Error,
		}, "An error occurred while 'ALTER TABLE departments AUTO_INCREMENT=1' in dao.EmptyDepartment.")
		return res.Error
	}
	return nil
}

// SetDepartment 更新部门
func (d *Dao) SetDepartment(ctx context.Context, id uint32, name string, parentId *uint32) (dept *model.Department, err error) {
	updatesMap := map[string]interface{}{"name": name}
	if parentId != nil {
		updatesMap["parent_id"] = *parentId
	}
	res := d.getDbWithCtx(ctx).Model(&model.Department{}).
		Where("id=?", id).
		Updates(updatesMap).
		Find(&dept)
	return dept, res.Error
}

// GetDepartment 查询单个部门
func (d *Dao) GetDepartment(ctx context.Context, q orm.Query) (dept *model.Department, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&dept)
	return dept, res.Error
}

// GetDepartmentCount 查询部门数量
func (d *Dao) GetDepartmentCount(ctx context.Context, q orm.Query) (count int64, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.Department{})).Count(&count)
	return count, res.Error
}

// IsDepartmentExist 查询部门是否存在
func (d *Dao) IsDepartmentExist(ctx context.Context, q orm.Query) (bool, error) {
	count, err := d.GetDepartmentCount(ctx, q)
	return count > 0, err
}

// ListDepartments 查询部门
func (d *Dao) ListDepartments(ctx context.Context, q orm.Query) (depts []*model.Department, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&depts)
	return depts, res.Error
}

// PagedListDepartments 分页查询部门
func (d *Dao) PagedListDepartments(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	depts := make([]*model.Department, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Department{}))
	return orm.PagingQuery(db, page, pageSize, &depts, orderBy...)
}

// ListDepartmentTree 以树形式列出所有部门
func (d *Dao) ListDepartmentTree(pNode *dto.DepartmentTree, deptList []*model.Department) {
	// 部门列表为空时直接返回
	if deptList == nil {
		return
	}

	for _, dept := range deptList {
		// 跳过root节点
		if dept.ParentId == nil {
			continue
		}
		deptParentId := *dept.ParentId
		if deptParentId == pNode.Id {
			dt := dto.TransDepartment2Tree(dept)
			d.ListDepartmentTree(dt, deptList)
			pNode.SubDepartment = append(pNode.SubDepartment, dt)
		}
	}

	return
}

// AddUser2Department 关联表操作::添加指定用户至指定部门
func (d *Dao) AddUser2Department(ctx context.Context, userId, deptId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Create(&model.UserDepartment{
		UserId:       userId,
		DepartmentId: deptId,
		IsOwner:      &isOwner,
	})
	return res.Error
}

// RemoveDepartmentsUser 关联表操作::移除指定部门中用户
func (d *Dao) RemoveDepartmentsUser(ctx context.Context, deptId, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserDepartment{}, "department_id=? AND user_id=?", deptId, userId)
	return res.Error
}

// RemoveUsersDepartments 关联表操作::移除指定用户所有部门
func (d *Dao) RemoveUsersDepartments(ctx context.Context, userId uint32) error {
	res := d.getDbWithCtx(ctx).Delete(model.UserDepartment{}, "user_id=?", userId)
	return res.Error
}

// SetDepartmentsOwner 关联表操作::设置用户是否是部门Owner
func (d *Dao) SetDepartmentsOwner(ctx context.Context, deptId, userId uint32, isOwner bool) error {
	res := d.getDbWithCtx(ctx).Model(&model.UserDepartment{}).
		Where("department_id=? AND user_id=?", deptId, userId).
		Update("is_owner", isOwner)
	return res.Error
}

// GetDepartmentUserCount 关联表操作::列出部门中所有用户数量
func (d *Dao) GetDepartmentUserCount(ctx context.Context, q orm.Query) (count int64, err error) {
	_db := d.getDbWithCtx(ctx).Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"ud.is_owner AS is_owner, " +
			"ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON users.id = ud.user_id")

	res := q.Where(_db).Count(&count)
	return count, res.Error
}

// IsEmptyDepartment 关联表操作::查询指定部门是否为空，不包含用户则为空
func (d *Dao) IsEmptyDepartment(ctx context.Context, deptId uint32) (bool, error) {
	count, err := d.GetDepartmentUserCount(ctx, orm.Query{"department_id=?": deptId})
	return count == 0, err
}

// ListDepartmentsUsers 关联表操作::列出部门中所有用户
func (d *Dao) ListDepartmentsUsers(
	ctx context.Context, deptId uint32, q orm.Query,
) (mbrUser []*model.MemberUser, err error) {
	res := q.Where(d.getDbWithCtx(ctx).Model(&model.User{})).
		Select("users.id AS id, "+
			"users.username AS username, "+
			"ud.is_owner AS is_owner, "+
			"ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON users.id = ud.user_id").
		Where("department_id=?", deptId).
		Find(&mbrUser)
	return mbrUser, res.Error
}

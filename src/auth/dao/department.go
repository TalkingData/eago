package dao

import (
	"eago/auth/model"
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewDepartment 新建部门
func NewDepartment(name string, parentId *int) (*model.Department, error) {
	d := model.Department{
		Name:     name,
		ParentId: parentId,
	}

	if res := db.Create(&d); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":      name,
			"parent_id": parentId,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return nil, res.Error
	}

	return &d, nil
}

// RemoveDepartment 删除部门
func RemoveDepartment(deptId int) bool {
	res := db.Delete(model.Department{}, "id=?", deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    deptId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// EmptyDepartment 清空部门
func EmptyDepartment() bool {
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.UserDepartment{})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Session.Delete.")
		return false
	}
	res = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Department{})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Session.Delete.")
		return false
	}
	return true
}

// SetDepartment 更新部门
func SetDepartment(id int, name string, parentId *int) (*model.Department, error) {
	d := model.Department{}

	ud := map[string]interface{}{"name": name}
	if parentId != nil {
		ud["department_id"] = *parentId
	}

	res := db.Model(&model.Department{}).
		Where("id=?", id).
		Updates(ud).
		First(&d)

	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":       id,
			"name":     name,
			"parentId": parentId,
			"error":    res.Error,
		}, "An error occurred while db.SetDepartment.")
		return nil, res.Error
	}

	return &d, nil
}

// GetDepartment 查询单个部门
func GetDepartment(query Query) (*model.Department, bool) {
	var (
		dept = model.Department{}
		d    = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&dept); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &dept, true
}

// GetDepartmentCount 查询部门数量
func GetDepartmentCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Department{})

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Count.")
		return count, false
	}
	return count, true
}

// ListDepartments 查询部门
func ListDepartments(query Query) (*[]model.Department, bool) {
	var d = db
	ds := make([]model.Department, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ds); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &ds, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &ds, true
}

// PagedListDepartments 查询部门-分页
func PagedListDepartments(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Department{})
	ds := make([]model.Department, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ds)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// Department2Tree 将部门转化为树结构的一个节点
func Department2Tree(dept *model.Department) *model.DepartmentTree {
	return &model.DepartmentTree{
		Id:            dept.Id,
		Name:          dept.Name,
		CreatedAt:     dept.CreatedAt,
		UpdatedAt:     dept.UpdatedAt,
		SubDepartment: make([]*model.DepartmentTree, 0),
	}
}

// ListDepartment2Tree 将部门列表转化为部门树
func ListDepartment2Tree(pNode *model.DepartmentTree, deptList *[]model.Department) {
	// 部门列表为空时直接返回
	if deptList == nil {
		return
	}

	for _, dept := range *deptList {

		// 跳过root节点
		if dept.ParentId == nil {
			continue
		}
		deptParentId := *dept.ParentId
		if deptParentId == pNode.Id {
			dt := Department2Tree(&dept)
			ListDepartment2Tree(dt, deptList)
			pNode.SubDepartment = append(pNode.SubDepartment, dt)
		}
	}

	return
}

// AddDepartmentUser 关联表操作::添加用户至部门
func AddDepartmentUser(userId, deptId int, isOwner bool) bool {
	var dp = model.UserDepartment{
		UserId:       userId,
		DepartmentId: deptId,
		IsOwner:      &isOwner,
	}

	if res := db.Create(&dp); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"is_owner":      isOwner,
			"joined_at":     dp.JoinedAt,
			"error":         res.Error,
		}, "An error occurred while db.Create.")
		return false
	}

	return true
}

// RemoveDepartmentUser 关联表操作::移除部门中用户
func RemoveDepartmentUser(userId, deptId int) bool {
	res := db.Delete(model.UserDepartment{}, "user_id=? AND department_id=?", userId, deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"error":         res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// RemoveUserDepartments 关联表操作::移除用户所有部门
func RemoveUserDepartments(userId int) bool {
	res := db.Delete(model.UserDepartment{}, "user_id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetDepartmentUserIsOwner 关联表操作::设置用户是否是部门Owner
func SetDepartmentUserIsOwner(deptId, userId int, isOwner bool) bool {
	res := db.Model(&model.UserDepartment{}).
		Where("department_id=? AND user_id=?", deptId, userId).
		Update("is_owner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"department_id": deptId,
			"user_id":       userId,
			"is_owner":      isOwner,
			"error":         res.Error,
		}, "An error occurred while db.Model.Where.Update.")
		return false
	}

	return true
}

// GetDepartmentUserCount 关联表操作::列出部门中所有用户数量
func GetDepartmentUserCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.User{}).
		Select("users.id AS id, users.username AS username, ud.is_owner AS is_owner, ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON users.id = ud.user_id")

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Count.")
		return count, false
	}

	return count, true
}

// ListDepartmentUsers 关联表操作::列出部门中所有用户
func ListDepartmentUsers(deptId int, query Query) (*[]model.MemberUser, bool) {
	var d = db.Model(&model.User{})
	mus := make([]model.MemberUser, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, users.username AS username, ud.is_owner AS is_owner, ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON users.id = ud.user_id").
		Where("department_id=?", deptId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &mus, true
}

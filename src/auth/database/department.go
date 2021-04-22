package database

import (
	"eago-common/api-suite/pagination"
	"eago-common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var DepartmentModel departmentModel

type departmentModel struct{}

type UserDepartment struct {
	Id           int    `json:"id" swaggerignore:"true"`
	UserId       int    `json:"user_id" binding:"required"`
	DepartmentId int    `json:"department_id" swaggerignore:"true"`
	IsOwner      *bool  `json:"is_owner" binding:"required"`
	JoinedAt     MyTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime" swaggerignore:"true"`
}

type Department struct {
	Id        int     `json:"id" swaggerignore:"true"`
	Name      string  `json:"name" binding:"required"`
	ParentId  *int    `json:"parent_id"`
	CreatedAt MyTime  `json:"created_at" swaggerignore:"true"`
	UpdatedAt *MyTime `json:"updated_at" swaggerignore:"true"`
}

type DepartmentTree struct {
	Id            int               `json:"id"`
	Name          string            `json:"name" binding:"required"`
	SubDepartment []*DepartmentTree `json:"sub_department"`
	CreatedAt     MyTime            `json:"created_at"`
	UpdatedAt     *MyTime           `json:"updated_at"`
}

// New 新建部门
func (dm *departmentModel) New(name string, parentId *int) *Department {
	var d = Department{
		Name:     name,
		ParentId: parentId,
	}

	if res := db.Create(&d); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":      name,
			"parent_id": parentId,
			"error":     res.Error.Error(),
		}, "Error in departmentModel.New.")
		return nil
	}

	return &d
}

// Remove 删除部门
func (dm *departmentModel) Remove(deptId int) bool {
	res := db.Delete(Department{}, "id=?", deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    deptId,
			"error": res.Error.Error(),
		}, "Error in departmentModel.Remove.")
		return false
	}

	return true
}

// Set 更新部门
func (dm *departmentModel) Set(id int, name string, parentId *int) (*Department, bool) {
	var d = Department{}

	res := db.Model(&Department{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":      name,
			"parent_id": *parentId,
		}).
		First(&d)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":        id,
			"name":      name,
			"parent_id": *parentId,
			"error":     res.Error.Error(),
		}, "Error in departmentModel.Set.")
		return nil, false
	}

	return &d, true
}

// Get 查询单个部门
func (dm *departmentModel) Get(query *Query) (*Department, bool) {
	var (
		dept = Department{}
		d    = db
	)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.First(&dept); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in departmentModel.Get.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in departmentModel.Get.")
		return nil, false
	}

	return &dept, true
}

// List 查询部门
func (dm *departmentModel) List(query *Query) (*[]Department, bool) {
	var d = db
	ds := make([]Department, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ds); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in departmentModel.List.")
			return &ds, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in departmentModel.List.")
		return nil, false
	}

	return &ds, true
}

// PagedList 查询部门-分页
func (dm *departmentModel) PagedList(query *Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Department{})
	ds := make([]Department, 0)

	for k, v := range *query {
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
			"error": err.Error(),
		}, "Error in departmentModel.PagedList.")
		return nil, false
	}

	return pg, true
}

// Department2Tree 将根部门转化为树结构
func (dm *departmentModel) Department2Tree(dept *Department) *DepartmentTree {
	var deptTree DepartmentTree

	deptTree.Id = dept.Id
	deptTree.Name = dept.Name
	deptTree.CreatedAt = dept.CreatedAt
	deptTree.UpdatedAt = dept.UpdatedAt
	deptTree.SubDepartment = make([]*DepartmentTree, 0)

	return &deptTree
}

// List2Tree 将部门列表转化为部门树
func (dm *departmentModel) List2Tree(pNode *DepartmentTree, deptList *[]Department) {
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
			dt := dm.Department2Tree(&dept)
			dm.List2Tree(dt, deptList)
			pNode.SubDepartment = append(pNode.SubDepartment, dt)
		}
	}

	return
}

// AddUser 关联表操作::添加用户至部门
func (dm *departmentModel) AddUser(userId, deptId int, isOwner bool) bool {
	var dp = UserDepartment{
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
			"error":         res.Error.Error(),
		}, "Error in departmentModel.AddUser.")
		return false
	}

	return true
}

// RemoveUser 关联表操作::移除部门中用户
func (dm *departmentModel) RemoveUser(userId, deptId int) bool {
	res := db.Delete(UserDepartment{}, "user_id=? AND department_id=?", userId, deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"error":         res.Error.Error(),
		}, "Error in departmentModel.RemoveUser.")
		return false
	}

	return true
}

// SetUserIsOwner 关联表操作::设置用户是否是部门Owner
func (dm *departmentModel) SetUserIsOwner(userId, deptId int, isOwner bool) bool {
	res := db.Model(&UserDepartment{}).
		Where("user_id=? AND department_id=?", userId, deptId).
		Update("IsOwner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"is_owner":      isOwner,
			"error":         res.Error.Error(),
		}, "Error in departmentModel.SetUserIsOwner.")
		return false
	}

	return true
}

// ListUsers 关联表操作::列出部门中所有用户
func (dm *departmentModel) ListUsers(deptId int, query *Query) (*[]memberUser, bool) {
	var d = db.Model(&User{})
	mus := make([]memberUser, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, users.username AS username, ud.is_owner AS is_owner, ud.joined_at AS joined_at").
		Joins("LEFT JOIN user_departments AS ud ON users.id = ud.user_id").
		Where("department_id=?", deptId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in departmentModel.ListUsers.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in departmentModel.ListUsers.")
		return nil, false
	}

	return &mus, true
}

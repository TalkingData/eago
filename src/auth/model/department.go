package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserDepartment struct {
	Id           int              `json:"id" swaggerignore:"true"`
	UserId       int              `json:"user_id" binding:"required"`
	DepartmentId int              `json:"department_id" swaggerignore:"true"`
	IsOwner      *bool            `json:"is_owner" binding:"required"`
	JoinedAt     *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime" swaggerignore:"true"`
}

type Department struct {
	Id        int              `json:"id" swaggerignore:"true"`
	Name      string           `json:"name" binding:"required"`
	ParentId  *int             `json:"parent_id"`
	CreatedAt *utils.LocalTime `json:"created_at" swaggerignore:"true"`
	UpdatedAt *utils.LocalTime `json:"updated_at" swaggerignore:"true"`
}

type DepartmentTree struct {
	Id            int               `json:"id"`
	Name          string            `json:"name" binding:"required"`
	SubDepartment []*DepartmentTree `json:"sub_department"`
	CreatedAt     *utils.LocalTime  `json:"created_at"`
	UpdatedAt     *utils.LocalTime  `json:"updated_at"`
}

// NewDepartment 新建部门
func NewDepartment(name string, parentId *int) *Department {
	var d = Department{
		Name:     name,
		ParentId: parentId,
	}

	if res := db.Create(&d); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":      name,
			"parent_id": parentId,
			"error":     res.Error.Error(),
		}, "Error in model.NewDepartment.")
		return nil
	}

	return &d
}

// RemoveDepartment 删除部门
func RemoveDepartment(deptId int) bool {
	res := db.Delete(Department{}, "id=?", deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    deptId,
			"error": res.Error.Error(),
		}, "Error in model.RemoveDepartment.")
		return false
	}

	return true
}

// SetDepartment 更新部门
func SetDepartment(id int, name string, parentId *int) (*Department, bool) {
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
		}, "Error in model.SetDepartment.")
		return nil, false
	}

	return &d, true
}

// GetDepartment 查询单个部门
func GetDepartment(query Query) (*Department, bool) {
	var (
		dept = Department{}
		d    = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&dept); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.GetDepartment.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.GetDepartment.")
		return nil, false
	}

	return &dept, true
}

// ListDepartments 查询部门
func ListDepartments(query Query) (*[]Department, bool) {
	var d = db
	ds := make([]Department, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ds); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.ListDepartments.")
			return &ds, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.ListDepartments.")
		return nil, false
	}

	return &ds, true
}

// PagedListDepartments 查询部门-分页
func PagedListDepartments(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Department{})
	ds := make([]Department, 0)

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
			"error": err.Error(),
		}, "Error in model.PagedListDepartments.")
		return nil, false
	}

	return pg, true
}

// Department2Tree 将部门转化为树结构的一个节点
func Department2Tree(dept *Department) *DepartmentTree {
	return &DepartmentTree{
		Id:            dept.Id,
		Name:          dept.Name,
		CreatedAt:     dept.CreatedAt,
		UpdatedAt:     dept.UpdatedAt,
		SubDepartment: make([]*DepartmentTree, 0),
	}
}

// ListDepartment2Tree 将部门列表转化为部门树
func ListDepartment2Tree(pNode *DepartmentTree, deptList *[]Department) {
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
		}, "Error in model.AddDepartmentUser.")
		return false
	}

	return true
}

// RemoveDepartmentUser 关联表操作::移除部门中用户
func RemoveDepartmentUser(userId, deptId int) bool {
	res := db.Delete(UserDepartment{}, "user_id=? AND department_id=?", userId, deptId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"error":         res.Error.Error(),
		}, "Error in model.RemoveDepartmentUser.")
		return false
	}

	return true
}

// SetDepartmentUserIsOwner 关联表操作::设置用户是否是部门Owner
func SetDepartmentUserIsOwner(userId, deptId int, isOwner bool) bool {
	res := db.Model(&UserDepartment{}).
		Where("user_id=? AND department_id=?", userId, deptId).
		Update("is_owner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"user_id":       userId,
			"department_id": deptId,
			"is_owner":      isOwner,
			"error":         res.Error.Error(),
		}, "Error in model.SetDepartmentUserIsOwner.")
		return false
	}

	return true
}

// ListDepartmentUsers 关联表操作::列出部门中所有用户
func ListDepartmentUsers(deptId int, query Query) (*[]memberUser, bool) {
	var d = db.Model(&User{})
	mus := make([]memberUser, 0)

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
				"error": res.Error.Error(),
			}, "Record not found in model.ListDepartmentUsers.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in model.ListDepartmentUsers.")
		return nil, false
	}

	return &mus, true
}

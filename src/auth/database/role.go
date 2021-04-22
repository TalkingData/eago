package database

import (
	"eago-common/api-suite/pagination"
	"eago-common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var RoleModel roleModel

type roleModel struct{}

type UserRole struct {
	Id       int    `json:"id" swaggerignore:"true"`
	UserId   int    `json:"user_id" binding:"required"`
	RoleId   int    `json:"role_id" swaggerignore:"true"`
	JoinedAt MyTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime" swaggerignore:"true"`
}

type Role struct {
	Id   int    `json:"id" swaggerignore:"true"`
	Name string `json:"name" binding:"required"`
}

type roleUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	JoinedAt MyTime `json:"joined_at"`
}

// New 新建角色
func (rm *roleModel) New(name string) *Role {
	var r = Role{Name: name}

	if res := db.Create(&r); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":  name,
			"error": res.Error.Error(),
		}, "Error in roleModel.New.")
		return nil
	}

	return &r
}

// Remove 删除角色
func (rm *roleModel) Remove(roleId int) bool {
	res := db.Delete(Role{}, "id=?", roleId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"role_id": roleId,
			"error":   res.Error.Error(),
		}, "Error in roleModel.Remove.")
		return false
	}

	return true
}

// Set 更新角色
func (rm *roleModel) Set(id int, name string) (*Role, bool) {
	var r = Role{}

	res := db.Model(&Role{}).
		Where("id=?", id).
		Updates(map[string]interface{}{"name": name}).
		First(&r)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in roleModel.Set.")
		return nil, false
	}

	return &r, true
}

// Get 查询单个角色
func (rm *roleModel) Get(query *Query) (*Role, bool) {
	var (
		r = Role{}
		d = db
	)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.First(&r); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in roleModel.Get.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in roleModel.Get.")
		return nil, false
	}

	return &r, true
}

// List 查询角色
func (rm *roleModel) List(query *Query) (*[]Role, bool) {
	var d = db
	rs := make([]Role, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}

	if res := d.Find(&rs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found on roleModel.List.")
			return &rs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in roleModel.List.")
		return nil, false
	}

	return &rs, true
}

// PagedList 查询角色-分页
func (rm *roleModel) PagedList(query *Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Role{})
	rs := make([]Role, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &rs)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in roleModel.PagedList.")
		return nil, false
	}

	return pg, true
}

// AddUser 关联表操作::添加用户至角色
func (rm *roleModel) AddUser(userId, roleId int) bool {
	var ur = UserRole{
		UserId: userId,
		RoleId: roleId,
	}

	if res := db.Create(&ur); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":   userId,
			"role_id":   roleId,
			"joined_at": ur.JoinedAt,
			"error":     res.Error.Error(),
		}, "Error in roleModel.AddUser.")
		return false
	}

	return true
}

// RemoveUser 关联表操作::移除角色中用户
func (rm *roleModel) RemoveUser(userId, roleId int) bool {
	res := db.Delete(UserRole{}, "user_id=? AND role_id=?", userId, roleId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"role_id": roleId,
			"error":   res.Error.Error(),
		}, "Error in roleModel.RemoveUser.")
		return false
	}

	return true
}

// ListUsers 关联表操作::列出角色中用户
func (rm *roleModel) ListUsers(roleId int) (*[]roleUser, bool) {
	rus := make([]roleUser, 0)

	res := db.Model(&User{}).
		Select("users.id AS id, users.username AS username, ur.joined_at AS joined_at").
		Joins("LEFT JOIN user_roles AS ur ON users.id = ur.user_id").
		Where("role_id=?", roleId).
		Find(&rus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in roleModel.ListUsers.")
			return &rus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in roleModel.ListUsers.")
		return nil, false
	}

	return &rus, true
}

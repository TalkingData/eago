package dao

import (
	"eago/auth/model"
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewRole 新建角色
func NewRole(name, description string) (*model.Role, error) {
	r := model.Role{
		Name:        name,
		Description: &description,
	}

	if res := db.Create(&r); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"description": description,
			"error":       res.Error,
		}, "An error occurred while db.Create.")
		return nil, res.Error
	}

	return &r, nil
}

// RemoveRole 删除角色
func RemoveRole(roleId int) error {
	res := db.Delete(model.Role{}, "id=?", roleId)
	if res.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}

// SetRole 更新角色
func SetRole(id int, name, description string) (*model.Role, error) {
	var r = model.Role{}

	res := db.Model(&model.Role{}).
		Where("id=?", id).
		Updates(map[string]interface{}{"name": name, "description": description}).
		First(&r)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":          id,
			"name":        name,
			"description": description,
			"error":       res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, res.Error
	}

	return &r, res.Error
}

// GetRole 查询单个角色
func GetRole(query Query) (*model.Role, bool) {
	var (
		r = model.Role{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&r); res.Error != nil {
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

	return &r, true
}

// GetRoleCount 查询角色数量
func GetRoleCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Role{})

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

// ListRoles 查询角色
func ListRoles(query Query) (*[]model.Role, bool) {
	var d = db
	rs := make([]model.Role, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Find(&rs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &rs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &rs, true
}

// PagedListRoles 查询角色-分页
func PagedListRoles(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Role{})
	rs := make([]model.Role, pageSize)

	for k, v := range query {
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
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// AddRoleUser 关联表操作::添加用户至角色
func AddRoleUser(roleId, userId int) bool {
	ur := model.UserRole{
		RoleId: roleId,
		UserId: userId,
	}

	if res := db.Create(&ur); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"role_id":   roleId,
			"user_id":   userId,
			"joined_at": ur.JoinedAt,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return false
	}

	return true
}

// RemoveRoleUser 关联表操作::移除角色中用户
func RemoveRoleUser(roleId, userId int) bool {
	res := db.Delete(model.UserRole{}, "role_id=? AND user_id=?", roleId, userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"role_id": roleId,
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// RemoveUserRoles 关联表操作::移除用户所有角色
func RemoveUserRoles(userId int) bool {
	res := db.Delete(model.UserRole{}, "user_id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// GetRoleUserCount 关联表操作::列出角色中用户数量
func GetRoleUserCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.User{}).
		Select("users.id AS id, " +
			"users.username AS username, " +
			"ur.joined_at AS joined_at").
		Joins("LEFT JOIN user_roles AS ur ON users.id = ur.user_id")

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

// ListRoleUsers 关联表操作::列出角色中用户
func ListRoleUsers(roleId int) (*[]model.RoleUser, bool) {
	rus := make([]model.RoleUser, 0)

	res := db.Model(&model.User{}).
		Select("users.id AS id, "+
			"users.username AS username, "+
			"ur.joined_at AS joined_at").
		Joins("LEFT JOIN user_roles AS ur ON users.id = ur.user_id").
		Where("role_id=?", roleId).
		Find(&rus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found.")
			return &rus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Select.Joins.Where.Find.")
		return nil, false
	}

	return &rus, true
}

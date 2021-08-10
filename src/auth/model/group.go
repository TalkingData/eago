package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserGroup struct {
	Id       int              `json:"id" swaggerignore:"true"`
	UserId   int              `json:"user_id" binding:"required"`
	GroupId  int              `json:"group_id" swaggerignore:"true"`
	IsOwner  *bool            `json:"is_owner" binding:"required"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime" swaggerignore:"true"`
}

type Group struct {
	Id          int              `json:"id" swaggerignore:"true"`
	Name        string           `json:"name" binding:"required"`
	Description *string          `json:"description" binding:"required"`
	CreatedAt   *utils.LocalTime `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *utils.LocalTime `json:"updated_at" swaggerignore:"true"`
}

// NewGroup 新建组
func NewGroup(name, description string) *Group {
	var g = Group{
		Name:        name,
		Description: &description,
	}

	if res := db.Create(&g); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"description": description,
			"error":       res.Error.Error(),
		}, "Error in model.New.")
		return nil
	}

	return &g
}

// RemoveGroup 删除组
func RemoveGroup(groupId int) bool {
	res := db.Delete(Group{}, "id=?", groupId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    groupId,
			"error": res.Error.Error(),
		}, "Error in model.RemoveGroup.")
		return false
	}

	return true
}

// SetGroup 更新组
func SetGroup(id int, name, description string) (*Group, bool) {
	var g = Group{}

	res := db.Model(&Group{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).
		First(&g)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in model.SetGroup.")
		return nil, false
	}

	return &g, true
}

// ListGroups 查询组
func ListGroups(query Query) (*[]Group, bool) {
	var d = db
	gs := make([]Group, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&gs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.ListGroups.")
			return &gs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.ListGroups.")
		return nil, false
	}

	return &gs, true
}

// PagedListGroups 查询组-分页
func PagedListGroups(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Group{})
	gs := make([]Group, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &gs)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": query,
			"error": err.Error(),
		}, "Error in model.PagedListGroups.")
		return nil, false
	}

	return pg, true
}

// AddGroupUser 关联表操作::添加用户至组
func AddGroupUser(userId, groupId int, isOwner bool) bool {
	var gp = UserGroup{
		UserId:  userId,
		GroupId: groupId,
		IsOwner: &isOwner,
	}

	if res := db.Create(&gp); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":   userId,
			"group_id":  groupId,
			"is_owner":  isOwner,
			"joined_at": gp.JoinedAt,
			"error":     res.Error.Error(),
		}, "Error in model.AddGroupUser.")
		return false
	}

	return true
}

// RemoveGroupUser 关联表操作::移除组中用户
func RemoveGroupUser(userId, groupId int) bool {
	res := db.Delete(UserGroup{}, "user_id=? AND group_id=?", userId, groupId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":  userId,
			"group_id": groupId,
			"error":    res.Error.Error(),
		}, "Error in model.RemoveGroupUser.")
		return false
	}

	return true
}

// SetGroupUserIsOwner 关联表操作::设置用户是否是组Owner
func SetGroupUserIsOwner(userId, groupId int, isOwner bool) bool {
	res := db.Model(&UserGroup{}).
		Where("user_id=? AND group_id=?", userId, groupId).
		Update("is_owner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"user_id":  userId,
			"group_id": groupId,
			"is_owner": isOwner,
			"error":    res.Error.Error(),
		}, "Error in model.SetGroupUserIsOwner.")
		return false
	}

	return true
}

// ListGroupUsers 关联表操作::列出组中所有用户
func ListGroupUsers(groupId int, query Query) (*[]memberUser, bool) {
	var d = db.Model(&User{})
	mus := make([]memberUser, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, users.username AS username, ug.is_owner AS is_owner, ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON users.id = ug.user_id").
		Where("group_id=?", groupId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in model.ListGroupUsers.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in model.ListGroupUsers.")
		return nil, false
	}

	return &mus, true
}

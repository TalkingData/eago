package database

import (
	"eago-common/api-suite/pagination"
	"eago-common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var GroupModel groupModel

type groupModel struct{}

type UserGroup struct {
	Id        int    `json:"id" swaggerignore:"true"`
	UserId    int    `json:"user_id" binding:"required"`
	GroupId   int    `json:"group_id" swaggerignore:"true"`
	IsOwner   *bool  `json:"is_owner" binding:"required"`
	CreatedAt MyTime `json:"created_at" swaggerignore:"true"`
}

type Group struct {
	Id          int     `json:"id" swaggerignore:"true"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	CreatedAt   MyTime  `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *MyTime `json:"updated_at" swaggerignore:"true"`
}

// 新建组
func (gm *groupModel) New(name string, description string) *Group {
	var g = Group{
		Name:        name,
		Description: description,
	}

	if res := db.Create(&g); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"description": description,
			"error":       res.Error.Error(),
		}, "Error in groupModel.New.")
		return nil
	}

	return &g
}

// 删除组
func (gm *groupModel) Delete(groupId int) bool {
	res := db.Delete(Group{}, "id=?", groupId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    groupId,
			"error": res.Error.Error(),
		}, "Error in groupModel.Delete.")
		return false
	}

	return true
}

// 更新组
func (gm *groupModel) Set(id int, name string, description string) (*Group, bool) {
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
		}, "Error in groupModel.Set.")
		return nil, false
	}

	return &g, true
}

// 查询组
func (gm *groupModel) List(query *Query) (*[]Group, bool) {
	var d = db
	gs := make([]Group, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	if res := d.Find(&gs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in groupModel.List.")
			return &gs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in groupModel.List.")
		return nil, false
	}

	return &gs, true
}

// 查询组-分页
func (gm *groupModel) PagedList(query *Query, page int, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Group{})
	gs := make([]Group, 0)

	for k, v := range *query {
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
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in groupModel.PagedList.")
		return nil, false
	}

	return pg, true
}

// 关联表操作::添加用户至组
func (gm *groupModel) AddUser(userId int, groupId int, isOwner bool) bool {
	var gp = UserGroup{
		UserId:  userId,
		GroupId: groupId,
		IsOwner: &isOwner,
	}

	if res := db.Create(&gp); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":    userId,
			"group_id":   groupId,
			"is_owner":   isOwner,
			"created_at": gp.CreatedAt,
			"error":      res.Error.Error(),
		}, "Error in groupModel.AddUser.")
		return false
	}

	return true
}

// 关联表操作::移除组中用户
func (gm *groupModel) RemoveUser(userId int, groupId int) bool {
	res := db.Delete(UserGroup{}, "user_id=? AND group_id=?", userId, groupId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":  userId,
			"group_id": groupId,
			"error":    res.Error.Error(),
		}, "Error in groupModel.RemoveUser.")
		return false
	}

	return true
}

// 关联表操作::设置用户是否是组Owner
func (gm *groupModel) SetUserIsOwner(userId int, groupId int, isOwner bool) bool {
	res := db.Model(&UserGroup{}).
		Where("user_id=? AND group_id=?", userId, groupId).
		Update("IsOwner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"user_id":  userId,
			"group_id": groupId,
			"is_owner": isOwner,
			"error":    res.Error.Error(),
		}, "Error in groupModel.SetUserIsOwner.")
		return false
	}

	return true
}

// 关联表操作::列出组中所有用户
func (gm *groupModel) ListUsers(groupId int, query *Query) (*[]memberUser, bool) {
	var d = db.Model(&User{})
	mus := make([]memberUser, 0)

	for k, v := range *query {
		d = d.Where(k, v)
	}
	res := d.Select("users.id AS id, users.username AS username, ug.is_owner AS is_owner, ug.created_at AS created_at").
		Joins("LEFT JOIN user_groups AS ug ON users.id = ug.user_id").
		Where("group_id=?", groupId).
		Find(&mus)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in groupModel.ListUsers.")
			return &mus, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in groupModel.ListUsers.")
		return nil, false
	}

	return &mus, true
}

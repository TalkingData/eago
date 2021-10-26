package dao

import (
	"eago/auth/model"
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewGroup 新建组
func NewGroup(name, description string) (*model.Group, error) {
	g := model.Group{
		Name:        name,
		Description: &description,
	}

	if res := db.Create(&g); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"description": description,
			"error":       res.Error,
		}, "An error occurred while db.Create.")
		return nil, res.Error
	}

	return &g, nil
}

// RemoveGroup 删除组
func RemoveGroup(groupId int) bool {
	res := db.Delete(model.Group{}, "id=?", groupId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    groupId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetGroup 更新组
func SetGroup(id int, name, description string) (*model.Group, error) {
	g := model.Group{}

	res := db.Model(&model.Group{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
		}).
		First(&g)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, res.Error
	}

	return &g, nil
}

// GetGroup 查询单个组
func GetGroup(query Query) (*model.Group, bool) {
	var (
		g = model.Group{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&g); res.Error != nil {
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

	return &g, true
}

// GetGroupCount 查询组数量
func GetGroupCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Group{})

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

// ListGroups 查询组
func ListGroups(query Query) (*[]model.Group, bool) {
	var d = db
	gs := make([]model.Group, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&gs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &gs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &gs, true
}

// PagedListGroups 查询组-分页
func PagedListGroups(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Group{})
	gs := make([]model.Group, 0)

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
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// AddGroupUser 关联表操作::添加用户至组
func AddGroupUser(groupId, userId int, isOwner bool) bool {
	gp := model.UserGroup{
		GroupId: groupId,
		UserId:  userId,
		IsOwner: &isOwner,
	}

	if res := db.Create(&gp); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"group_id":  groupId,
			"user_id":   userId,
			"is_owner":  isOwner,
			"joined_at": gp.JoinedAt,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return false
	}

	return true
}

// RemoveGroupUser 关联表操作::移除组中用户
func RemoveGroupUser(groupId, userId int) bool {
	res := db.Delete(model.UserGroup{}, "group_id=? AND user_id=?", groupId, userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"group_id": groupId,
			"user_id":  userId,
			"error":    res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// RemoveUserGroups 关联表操作::移除用户所有组
func RemoveUserGroups(userId int) bool {
	res := db.Delete(model.UserGroup{}, "user_id=?", userId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"user_id": userId,
			"error":   res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetGroupUserIsOwner 关联表操作::设置用户是否是组Owner
func SetGroupUserIsOwner(groupId, userId int, isOwner bool) bool {
	res := db.Model(&model.UserGroup{}).
		Where("group_id=? AND user_id=?", groupId, userId).
		Update("is_owner", isOwner)
	if res.Error != nil {
		log.WarnWithFields(log.Fields{
			"group_id": groupId,
			"user_id":  userId,
			"is_owner": isOwner,
			"error":    res.Error,
		}, "An error occurred while db.Model.Where.Update.")
		return false
	}

	return true
}

// GetGroupUserCount 关联表操作::列出组中所有用户数量
func GetGroupUserCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.User{}).
		Select("users.id AS id, users.username AS username, ug.is_owner AS is_owner, ug.joined_at AS joined_at").
		Joins("LEFT JOIN user_groups AS ug ON users.id = ug.user_id")

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

// ListGroupUsers 关联表操作::列出组中所有用户
func ListGroupUsers(groupId int, query Query) (*[]model.MemberUser, bool) {
	var d = db.Model(&model.User{})
	mus := make([]model.MemberUser, 0)

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

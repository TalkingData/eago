package dto

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/common/message"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type NewGroup struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0)"`
}

// Valid
func (ng *NewGroup) Valid(v *validation.Validation) {
	if ct, _ := dao.GetGroupCount(dao.Query{"name=?": ng.Name}); ct > 0 {
		_ = v.SetError("Name", "组名称已存在")
	}
}

// Validate
func (ng *NewGroup) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(ng)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveGroup struct
type RemoveGroup struct{}

// Validate
func (*RemoveGroup) Validate(gId int) *message.Message {
	// 验证组是否存在
	if ct, _ := dao.GetGroupCount(dao.Query{"id=?": gId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("组不存在")
	}

	// 验证组是否有关联存在
	if ct, _ := dao.GetGroupUserCount(dao.Query{"group_id=?": gId}); ct > 0 {
		return msg.AssociatedGroupFailed
	}

	return nil
}

type SetGroup struct {
	groupId int

	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0)"`
}

// Valid
func (sg *SetGroup) Valid(v *validation.Validation) {
	if ct, _ := dao.GetGroupCount(dao.Query{"name=?": sg.Name, "id<>?": sg.groupId}); ct > 0 {
		_ = v.SetError("Name", "组名称已存在")
	}
}

// Validate
func (sg *SetGroup) Validate(gId int) *message.Message {
	sg.groupId = gId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(sg)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListGroupsQuery struct
type ListGroupsQuery struct {
	Query *string `form:"query"`
}

// UpdateQuery
func (lgq *ListGroupsQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if lgq.Query != nil && *lgq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *lgq.Query)
		query["name LIKE @query OR id LIKE @query"] = sql.Named("query", likeQuery)
	}

	return nil
}

// AddUser2Group struct
type AddUser2Group struct {
	groupId int

	UserId  int  `json:"user_id" valid:"Required"`
	IsOwner bool `json:"is_owner" valid:"Required"`
}

// Valid
func (aug *AddUser2Group) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": aug.UserId}); ct < 1 {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该组
	if ct, _ := dao.GetGroupUserCount(dao.Query{"group_id=?": aug.groupId, "user_id": aug.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该组")
	}
}

// Validate
func (aug *AddUser2Group) Validate(gId int) *message.Message {
	// 验证组是否存在
	if ct, _ := dao.GetGroupCount(dao.Query{"id=?": gId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("组不存在")
	}

	aug.groupId = gId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(aug)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveGroupUser struct
type RemoveGroupUser struct{}

// Validate
func (*RemoveGroupUser) Validate(gId, userId int) *message.Message {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetGroupUserCount(dao.Query{"group_id=?": gId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

// SetUserIsGroupOwner struct
type SetUserIsGroupOwner struct {
	IsOwner bool `json:"is_owner" valid:"Required"`
}

// Validate
func (suo *SetUserIsGroupOwner) Validate(groupId, userId int) *message.Message {
	// 验证用户是否已经属于该组
	if ct, _ := dao.GetGroupUserCount(dao.Query{"group_id=?": groupId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(suo)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListGroupUsersQuery struct
type ListGroupUsersQuery struct {
	IsOwner *bool `form:"is_owner"`
}

func (*ListGroupUsersQuery) Validate(groupId int) *message.Message {
	// 组不存在
	if ct, _ := dao.GetGroupCount(dao.Query{"id=?": groupId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("组不存在")
	}

	return nil
}

// UpdateQuery
func (lgu *ListGroupUsersQuery) UpdateQuery(query dao.Query) error {
	if lgu.IsOwner != nil {
		query["is_owner=?"] = *lgu.IsOwner
	}

	return nil
}

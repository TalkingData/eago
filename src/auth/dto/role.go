package dto

import (
	"database/sql"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/common/message"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewRole struct
type NewRole struct {
	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (nr *NewRole) Valid(v *validation.Validation) {
	if ct, _ := dao.GetRoleCount(dao.Query{"name=?": nr.Name}); ct > 0 {
		_ = v.SetError("Name", "已有相同名称的角色存在")
	}
}

func (nr *NewRole) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(nr)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveRole struct
type RemoveRole struct{}

func (*RemoveRole) Validate(roleId int) *message.Message {
	// 验证角色是否存在
	if ct, _ := dao.GetRoleCount(dao.Query{"id=?": roleId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("角色不存在")
	}

	// 验证角色是否有关联存在
	if ct, _ := dao.GetRoleUserCount(dao.Query{"role_id=?": roleId}); ct > 0 {
		return msg.AssociatedRoleFailed
	}

	return nil
}

// SetRole struct
type SetRole struct {
	roleId int

	Name        string  `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (sr *SetRole) Valid(v *validation.Validation) {
	if ct, _ := dao.GetRoleCount(dao.Query{"name=?": sr.Name, "id<>?": sr.roleId}); ct > 0 {
		_ = v.SetError("Name", "已有相同名称的角色存在")
	}
}

func (sr *SetRole) Validate(roleId int) *message.Message {
	// 验证角色是否存在
	if ct, _ := dao.GetRoleCount(dao.Query{"id=?": roleId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("角色不存在")
	}

	sr.roleId = roleId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(sr)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// PagedListRolesQuery struct
type PagedListRolesQuery struct {
	Query *string `form:"query"`
}

func (lrq *PagedListRolesQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if lrq.Query != nil && *lrq.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *lrq.Query)
		query["(name LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}
	return nil
}

// AddUser2Role struct
type AddUser2Role struct {
	roleId int

	UserId int `json:"user_id" valid:"Required"`
}

func (aur *AddUser2Role) Valid(v *validation.Validation) {
	// 验证用户是否存在
	if ct, _ := dao.GetUserCount(dao.Query{"id=?": aur.UserId}); ct < 1 {
		_ = v.SetError("UserId", "用户不存在")
	}

	// 验证用户是否已经属于该角色
	if ct, _ := dao.GetRoleUserCount(dao.Query{"role_id=?": aur.roleId, "user_id": aur.UserId}); ct > 0 {
		_ = v.SetError("UserId", "用户已经属于该角色")
	}
}

func (aur *AddUser2Role) Validate(roleId int) *message.Message {
	// 验证角色是否存在
	if ct, _ := dao.GetRoleCount(dao.Query{"id=?": roleId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("角色不存在")
	}

	aur.roleId = roleId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(aur)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveRoleUser struct
type RemoveRoleUser struct{}

func (*RemoveRoleUser) Validate(roleId, userId int) *message.Message {
	// 验证用户是否已经属于该角色
	if ct, _ := dao.GetRoleUserCount(dao.Query{"role_id=?": roleId, "user_id": userId}); ct < 1 {
		return msg.NotFoundFailed
	}

	return nil
}

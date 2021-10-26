package model

import (
	"eago/common/utils"
	"github.com/beego/beego/v2/core/validation"
)

type Group struct {
	Id          int              `json:"id"`
	Name        string           `json:"name" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-]{1,}$/)"`
	Description *string          `json:"description" valid:"MinSize(0)"`
	CreatedAt   *utils.LocalTime `json:"created_at"`
	UpdatedAt   *utils.LocalTime `json:"updated_at"`
}

type UserGroup struct {
	Id       int              `json:"id"`
	UserId   int              `json:"user_id" valid:"Required"`
	GroupId  int              `json:"group_id"`
	IsOwner  *bool            `json:"is_owner" valid:"Required"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// Validate
func (g *Group) Validate() (err interface{}) {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(g)
	if err != nil {
		return
	}
	// 数据验证未通过
	if !ok {
		return valid.Errors
	}

	return
}

// Validate
func (ug *UserGroup) Validate() (err interface{}) {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(ug)
	if err != nil {
		return
	}
	// 数据验证未通过
	if !ok {
		return valid.Errors
	}

	return
}

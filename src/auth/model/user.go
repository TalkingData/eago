package model

import (
	"eago/common/utils"
)

// User 用户-数据库模型
type User struct {
	Id uint32 `json:"id"`

	Username    string `json:"username"`
	Password    string `json:"-"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IsSuperuser bool   `json:"is_superuser"`
	Disabled    bool   `json:"disabled"`

	LastLogin *utils.CustomTime `json:"last_login"`
	CreatedAt *utils.CustomTime `json:"created_at"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
}

// MemberUsers 成员用户-通用
type MemberUser struct {
	Id uint32 `json:"id"`

	Username string            `json:"username"`
	IsOwner  bool              `json:"is_owner"`
	JoinedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// UserGroupMember 用户所在组的信息
type UserGroupMember struct {
	Id uint32 `json:"id"`

	Name        string            `json:"username"`
	Description *string           `json:"description"`
	IsOwner     bool              `json:"is_owner"`
	JoinedAt    *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// UserProductMember 用户所在产品线的信息
type UserProductMember struct {
	Id uint32 `json:"id"`

	Name     string            `json:"name"`
	Alias    string            `json:"alias"`
	Disabled bool              `json:"disabled"`
	IsOwner  bool              `json:"is_owner"`
	JoinedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// UserDepartmentMember 用户所在部门的信息
type UserDepartmentMember struct {
	Id uint32 `json:"id"`

	Name     string            `json:"name"`
	ParentId *uint32           `json:"parent_id"`
	IsOwner  bool              `json:"is_owner"`
	JoinedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// UserDepartmentNode 用户所在部门的信息-节点模式
type UserDepartmentNode struct {
	Id uint32 `json:"id"`

	Name     string  `json:"name"`
	ParentId *uint32 `json:"parent_id"`
}

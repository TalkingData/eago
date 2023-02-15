package model

import (
	"eago/common/utils"
)

// Role 角色-数据库模型
type Role struct {
	Id uint32 `json:"id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// UserRole 用户角色关联-数据库模型
type UserRole struct {
	Id uint32 `json:"id"`

	UserId   uint32            `json:"user_id"`
	RoleId   uint32            `json:"role_id"`
	JoinedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

// RolesUser 用户所属角色的信息
type RolesUser struct {
	Id uint32 `json:"id"`

	Username string            `json:"username"`
	JoinedAt *utils.CustomTime `json:"joined_at"`
}

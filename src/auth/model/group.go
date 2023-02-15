package model

import (
	"eago/common/utils"
)

// Group 组-数据库模型
type Group struct {
	Id uint32 `json:"id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
}

// UserGroup 用户组关联-数据库模型
type UserGroup struct {
	Id uint32 `json:"id"`

	UserId   uint32            `json:"user_id"`
	GroupId  uint32            `json:"group_id"`
	IsOwner  *bool             `json:"is_owner"`
	JoinedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

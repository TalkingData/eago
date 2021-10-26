package model

import (
	"eago/common/utils"
)

type Role struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type RoleUser struct {
	Id       int              `json:"id"`
	Username string           `json:"username"`
	JoinedAt *utils.LocalTime `json:"joined_at"`
}

type UserRole struct {
	Id       int              `json:"id"`
	UserId   int              `json:"user_id"`
	RoleId   int              `json:"role_id"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

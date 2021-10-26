package model

import (
	"eago/common/utils"
)

type User struct {
	Id          int              `json:"id"`
	Username    string           `json:"username"`
	Password    string           `json:"-"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	IsSuperuser bool             `json:"is_superuser"`
	Disabled    bool             `json:"disabled"`
	LastLogin   *utils.LocalTime `json:"last_login"`
	CreatedAt   *utils.LocalTime `json:"created_at"`
	UpdatedAt   *utils.LocalTime `json:"updated_at"`
}

type UserMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type MemberUser struct {
	Id       int              `json:"id"`
	Username string           `json:"username"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type UserProductMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	Alias    string           `json:"alias"`
	Disabled bool             `json:"disabled"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

type UserDepartmentMember struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	ParentId *int             `json:"parent_id"`
	IsOwner  bool             `json:"is_owner"`
	JoinedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

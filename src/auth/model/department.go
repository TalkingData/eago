package model

import (
	"eago/common/utils"
)

// Department 部门-数据库模型
type Department struct {
	Id uint32 `json:"id"`

	Name     string  `json:"name"`
	ParentId *uint32 `json:"parent_id"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
}

// UserDepartment 用户部门关联-数据库模型
type UserDepartment struct {
	Id uint32 `json:"id"`

	UserId       uint32            `json:"user_id"`
	DepartmentId uint32            `json:"department_id"`
	IsOwner      *bool             `json:"is_owner"`
	JoinedAt     *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

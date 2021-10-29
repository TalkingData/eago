package model

import (
	"eago/common/utils"
)

type Department struct {
	Id        int              `json:"id"`
	Name      string           `json:"name" valid:"Required;MinSize(3);MaxSize(100)"`
	ParentId  *int             `json:"parent_id"`
	CreatedAt *utils.LocalTime `json:"created_at"`
	UpdatedAt *utils.LocalTime `json:"updated_at"`
}

type DepartmentTree struct {
	Id            int               `json:"id"`
	Name          string            `json:"name"`
	SubDepartment []*DepartmentTree `json:"sub_department"`
	CreatedAt     *utils.LocalTime  `json:"created_at"`
	UpdatedAt     *utils.LocalTime  `json:"updated_at"`
}

type UserDepartment struct {
	Id           int              `json:"id"`
	UserId       int              `json:"user_id" valid:"Required"`
	DepartmentId int              `json:"department_id"`
	IsOwner      *bool            `json:"is_owner" valid:"Required"`
	JoinedAt     *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
}

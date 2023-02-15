package model

import (
	"eago/common/utils"
)

type Trigger struct {
	Id uint32 `json:"id"`

	Name         string  `json:"name"`
	Description  *string `json:"description"`
	TaskCodename string  `json:"task_codename"`
	Arguments    string  `json:"arguments"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

package model

import (
	"eago/common/utils"
)

// Schedule struct
type Schedule struct {
	Id           uint32            `json:"id"`
	TaskCodename string            `json:"task_codename"`
	Expression   string            `json:"expression"`
	Timeout      *int64            `json:"timeout"`
	Arguments    string            `json:"arguments"`
	Disabled     *bool             `json:"disabled"`
	Description  *string           `json:"description"`
	CreatedAt    *utils.CustomTime `json:"created_at"`
	CreatedBy    string            `json:"created_by"`
	UpdatedAt    *utils.CustomTime `json:"updated_at"`
	UpdatedBy    *string           `json:"updated_by" gorm:"default:''"`
}

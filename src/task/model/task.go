package model

import (
	"eago/common/utils"
)

// Task struct
type Task struct {
	Id           uint32            `json:"id"`
	Category     *int32            `json:"category"`
	Codename     string            `json:"codename"`
	FormalParams string            `json:"formal_params"`
	Disabled     *bool             `json:"disabled" gorm:"default:0"`
	Description  *string           `json:"description"`
	CreatedAt    *utils.CustomTime `json:"created_at"`
	CreatedBy    string            `json:"created_by"`
	UpdatedAt    *utils.CustomTime `json:"updated_at"`
	UpdatedBy    *string           `json:"updated_by" gorm:"default:''"`
}

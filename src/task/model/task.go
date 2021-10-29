package model

import (
	"eago/common/utils"
)

// Task struct
type Task struct {
	Id           int              `json:"id"`
	Category     *int             `json:"category"`
	Codename     string           `json:"codename"`
	FormalParams string           `json:"formal_params"`
	Disabled     *bool            `json:"disabled" gorm:"default:0"`
	Description  *string          `json:"description"`
	CreatedAt    *utils.LocalTime `json:"created_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedAt    *utils.LocalTime `json:"updated_at"`
	UpdatedBy    *string          `json:"updated_by" gorm:"default:''"`
}

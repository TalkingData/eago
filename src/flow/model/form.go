package model

import (
	"eago/common/utils"
)

// Form struct
type Form struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Disabled    *bool            `json:"disabled"`
	Description *string          `json:"description"`
	Body        *string          `json:"body"`
	CreatedAt   *utils.LocalTime `json:"created_at"`
	CreatedBy   string           `json:"created_by"`
	UpdatedAt   *utils.LocalTime `json:"updated_at"`
	UpdatedBy   *string          `json:"updated_by" gorm:"default:''"`
}

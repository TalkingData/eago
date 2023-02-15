package model

import (
	"eago/common/utils"
)

type Form struct {
	Id uint32 `json:"id"`

	Name        string  `json:"name"`
	Disabled    *bool   `json:"disabled"`
	Description *string `json:"description"`
	Body        *string `json:"body"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

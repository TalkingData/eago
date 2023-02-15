package model

import (
	"eago/common/utils"
)

// Categories 类别，此模块需要注意复数形式变化
type Categories struct {
	Id uint32 `json:"id"`

	Name string `json:"name"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

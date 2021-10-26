package model

import (
	"eago/common/utils"
)

// Categories 类别，此模块需要注意复数形式变化
type Categories struct {
	Id        int              `json:"id"`
	Name      string           `json:"name"`
	CreatedAt *utils.LocalTime `json:"created_at"`
	CreatedBy string           `json:"created_by"`
	UpdatedAt *utils.LocalTime `json:"updated_at"`
	UpdatedBy *string          `json:"updated_by" gorm:"default:''"`
}

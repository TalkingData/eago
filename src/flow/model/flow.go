package model

import (
	"eago/common/utils"
)

// Flow struct
type Flow struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	CategoriesId *int             `json:"categories_id"`
	Disabled     *bool            `json:"disabled"`
	Description  *string          `json:"description"`
	FormId       int              `json:"form_id"`
	FirstNodeId  int              `json:"first_node_id"`
	CreatedAt    *utils.LocalTime `json:"created_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedAt    *utils.LocalTime `json:"updated_at"`
	UpdatedBy    *string          `json:"updated_by" gorm:"default:''"`
}

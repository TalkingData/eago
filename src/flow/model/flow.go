package model

import (
	"eago/common/utils"
)

type Flow struct {
	Id uint32 `json:"id"`

	Name          string  `json:"name"`
	InstanceTitle string  `json:"instance_title"`
	CategoriesId  *uint32 `json:"categories_id"`
	Disabled      *bool   `json:"disabled"`
	Description   *string `json:"description"`

	FormId      uint32 `json:"form_id"`
	FirstNodeId uint32 `json:"first_node_id"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

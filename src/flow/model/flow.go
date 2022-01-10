package model

import (
	"eago/common/utils"
)

// Flow struct
type Flow struct {
	Id            int              `json:"id"`
	Name          string           `json:"name"`
	InstanceTitle string           `json:"instance_title"`
	CategoriesId  *int             `json:"categories_id"`
	Disabled      *bool            `json:"disabled"`
	Description   *string          `json:"description"`
	FormId        int              `json:"form_id"`
	FirstNodeId   int              `json:"first_node_id"`
	CreatedAt     *utils.LocalTime `json:"created_at"`
	CreatedBy     string           `json:"created_by"`
	UpdatedAt     *utils.LocalTime `json:"updated_at"`
	UpdatedBy     *string          `json:"updated_by" gorm:"default:''"`
}

// ListFlows struct
type ListFlows struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	InstanceTitle  string           `json:"instance_title"`
	CategoriesId   *int             `json:"categories_id"`
	CategoriesName string           `json:"categories_name"`
	Disabled       *bool            `json:"disabled"`
	Description    *string          `json:"description"`
	FormId         int              `json:"form_id"`
	FormName       string           `json:"form_name"`
	FirstNodeId    int              `json:"first_node_id"`
	FirstNodeName  string           `json:"first_node_name"`
	CreatedAt      *utils.LocalTime `json:"created_at"`
	CreatedBy      string           `json:"created_by"`
	UpdatedAt      *utils.LocalTime `json:"updated_at"`
	UpdatedBy      *string          `json:"updated_by" gorm:"default:''"`
}

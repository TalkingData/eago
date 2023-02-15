package dto

import "eago/common/utils"

type ListFlows struct {
	Id             uint32            `json:"id"`
	Name           string            `json:"name"`
	InstanceTitle  string            `json:"instance_title"`
	CategoriesId   *uint32           `json:"categories_id"`
	CategoriesName string            `json:"categories_name"`
	Disabled       *bool             `json:"disabled"`
	Description    *string           `json:"description"`
	FormId         uint32            `json:"form_id"`
	FormName       string            `json:"form_name"`
	FirstNodeId    uint32            `json:"first_node_id"`
	FirstNodeName  string            `json:"first_node_name"`
	CreatedAt      *utils.CustomTime `json:"created_at"`
	CreatedBy      string            `json:"created_by"`
	UpdatedAt      *utils.CustomTime `json:"updated_at"`
	UpdatedBy      *string           `json:"updated_by"`
}

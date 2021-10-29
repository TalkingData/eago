package model

import (
	"eago/common/utils"
)

// Trigger struct
type Trigger struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	Description  *string          `json:"description"`
	TaskCodename string           `json:"task_codename"`
	Arguments    string           `json:"arguments"`
	CreatedAt    *utils.LocalTime `json:"created_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedAt    *utils.LocalTime `json:"updated_at"`
	UpdatedBy    *string          `json:"updated_by" gorm:"default:''"`
}

// TriggersNode struct
type TriggersNode struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id"`
}

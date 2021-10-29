package model

import (
	"eago/common/utils"
)

// Instance struct
type Instance struct {
	Id                int              `json:"id"`
	Name              string           `json:"name"`
	Status            int              `json:"status" gorm:"type:int(11) NOT NULL;index"`
	FormId            int              `json:"form_id"`
	FormData          *string          `json:"form_data" gorm:"default:'{}'"`
	FlowChain         *string          `json:"flow_chain" gorm:"default:'{}'"`
	CurrentStep       int              `json:"current_step"`
	AssigneesRequired int              `json:"assignees_required"`
	CurrentAssignees  string           `json:"current_assignees"`
	PassedAssignees   string           `json:"passed_assignees"`
	CreatedAt         *utils.LocalTime `json:"created_at"`
	CreatedBy         string           `json:"created_by"`
	UpdatedAt         *utils.LocalTime `json:"updated_at"`
	UpdatedBy         *string          `json:"updated_by" gorm:"default:''"`
}

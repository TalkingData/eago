package model

import (
	"eago/common/utils"
)

type Instance struct {
	Id uint32 `json:"id"`

	Name   string `json:"name"`
	Status int32  `json:"status" gorm:"type:int(11) NOT NULL;index"`

	FormId    uint32  `json:"form_id"`
	FormData  *string `json:"form_data" gorm:"default:'{}'"`
	FlowChain *string `json:"flow_chain" gorm:"default:'{}'"`

	CurrentStep       int32  `json:"current_step"`
	AssigneesRequired int32  `json:"assignees_required"`
	CurrentAssignees  string `json:"current_assignees"`
	PassedAssignees   string `json:"passed_assignees"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

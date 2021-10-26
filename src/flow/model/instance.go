package model

import (
	"eago/common/utils"
)

const ASSIGNEES_SPILT_TAG = ","

const (
	INSTANCE_PANIC_END_STATUS    = 200 // 系统异常
	INSTANCE_REJECTED_END_STATUS = -1  // 被驳回
	INSTANCE_APPROVED_END_STATUS = 0   // 审批通过
	INSTANCE_PENDING_STATUS      = 1   // 系统处理中
	INSTANCE_RUNNING_STATUS      = 2   // 流转中
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

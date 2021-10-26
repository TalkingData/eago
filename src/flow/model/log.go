package model

import (
	"eago/common/utils"
)

// Log 审批日志
type Log struct {
	Id         int              `json:"id"`
	InstanceId int              `json:"instance_id"`
	Result     bool             `json:"result"`
	Content    *string          `json:"content"`
	CreatedAt  *utils.LocalTime `json:"created_at"`
	CreatedBy  string           `json:"created_by"`
}

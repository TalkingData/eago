package model

import (
	"eago/common/utils"
)

// Log 审批日志
type Log struct {
	Id uint64 `json:"id"`

	InstanceId uint32  `json:"instance_id"`
	Result     bool    `json:"result"`
	Content    *string `json:"content"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
}

package model

import (
	"eago/common/utils"
)

const (
	BUTILIN_TASK_CATEGORY = 1   // 内置任务
	BASH_TASK_CATEGORY    = 100 // Bash任务
	PYTHON_TASK_CATEGORY  = 101 // Python任务
)

// Task struct
type Task struct {
	Id           int              `json:"id"`
	Category     *int             `json:"category"`
	Codename     string           `json:"codename"`
	FormalParams string           `json:"formal_params"`
	Disabled     *bool            `json:"disabled" gorm:"default:0"`
	Description  *string          `json:"description"`
	CreatedAt    *utils.LocalTime `json:"created_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedAt    *utils.LocalTime `json:"updated_at"`
	UpdatedBy    *string          `json:"updated_by" gorm:"default:''"`
}

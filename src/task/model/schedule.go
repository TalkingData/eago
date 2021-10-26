package model

import (
	"eago/common/utils"
	"github.com/beego/beego/v2/core/validation"
)

// Schedule struct
type Schedule struct {
	Id           int              `json:"id"`
	TaskCodename string           `json:"task_codename"`
	Expression   string           `json:"expression"`
	Timeout      *int64           `json:"timeout"`
	Arguments    string           `json:"arguments"`
	Disabled     *bool            `json:"disabled"`
	Description  *string          `json:"description"`
	CreatedAt    *utils.LocalTime `json:"created_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedAt    *utils.LocalTime `json:"updated_at"`
	UpdatedBy    *string          `json:"updated_by" gorm:"default:''"`
}

// Validate
func (s *Schedule) Validate() (err interface{}) {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(s)
	if err != nil {
		return
	}
	// 数据验证未通过
	if !ok {
		return valid.Errors
	}

	return
}

package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/flow/conf/msg"
	"eago/task/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

// NewSchedule struct
type NewSchedule struct {
	TaskCodename string  `json:"task_codename" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-.]{1,}$/)"`
	Expression   string  `json:"expression" valid:"Required"`
	Timeout      *int    `json:"timeout" valid:"Min(0)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (n *NewSchedule) Validate() *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(n)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// RemoveSchedule struct
type RemoveSchedule struct{}

func (*RemoveSchedule) Validate(schId int) *message.Message {
	// 验证任务是否存在
	if ct, _ := dao.GetScheduleCount(dao.Query{"id=?": schId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("计划任务不存在")
	}

	return nil
}

// SetSchedule struct
type SetSchedule struct {
	TaskCodename string  `json:"task_codename" valid:"Required;MinSize(3);MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9_-.]{1,}$/)"`
	Expression   string  `json:"expression" valid:"Required"`
	Timeout      *int64  `json:"timeout" valid:"Min(0)"`
	Arguments    string  `json:"arguments" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

func (n *SetSchedule) Validate(schId int) *message.Message {
	// 验证任务是否存在
	if ct, _ := dao.GetScheduleCount(dao.Query{"id=?": schId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("计划任务不存在")
	}

	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(n)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListSchedulesQuery struct
type ListSchedulesQuery struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

func (q *ListSchedulesQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(task_codename LIKE @query OR "+
			"description LIKE @query OR "+
			"id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Disabled != nil {
		query["disabled=?"] = *q.Disabled
	}

	return nil
}

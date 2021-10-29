package dto

import (
	"database/sql"
	"eago/common/message"
	"eago/task/conf"
	"eago/task/conf/msg"
	"eago/task/dao"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
)

type CallTask struct {
	TaskCodeName string

	Timeout   *int64 `json:"timeout" valid:"Range(0,86400000)"`
	Arguments string `json:"arguments" valid:"Required;MaxSize(2)"`
}

// Validate
func (ct *CallTask) Validate(tId int) *message.Message {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(ct)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	// 验证任务是否存在
	tObj, ok := dao.GetTask(dao.Query{"id=?": tId, "disabled=?": false})
	// 查找操作失败
	if !ok {
		return msg.NotFoundFailed
	}
	// 没找到任务或者任务是禁用状态
	if tObj == nil {
		return msg.NotFoundFailed.SetDetail("任务不存在或任务不是可用的状态")
	}

	ct.TaskCodeName = tObj.Codename

	return nil
}

type NewTask struct {
	Category     *int    `json:"category" valid:"Min(0)"`
	Codename     string  `json:"codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	FormalParams string  `json:"formal_params" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

// Valid
func (n *NewTask) Valid(v *validation.Validation) {
	if n.Category != nil && *n.Category > conf.BUTILIN_TASK_CATEGORY {
		_ = v.SetError("Category", "目前仅支持内置任务")
	}

	if ct, _ := dao.GetTaskCount(dao.Query{"codename=?": n.Codename}); ct > 0 {
		_ = v.SetError("Codename", "已有相同代号的任务存在")
	}
}

// Validate
func (n *NewTask) Validate() *message.Message {
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

// RemoveTask struct
type RemoveTask struct{}

// Validate
func (*RemoveTask) Validate(tId int) *message.Message {
	// 验证任务是否存在
	tObj, _ := dao.GetTask(dao.Query{"id=?": tId})
	// 没找到任务或者任务是禁用状态
	if tObj == nil {
		return msg.NotFoundFailed.SetDetail("任务不存在")
	}

	// 验证任务是否存在
	if ct, _ := dao.GetScheduleCount(dao.Query{"task_codename=?": tObj.Codename}); ct > 0 {
		return msg.AssociatedScheduleFailed
	}

	return nil
}

// SetTask struct
type SetTask struct {
	taskId int

	Category     *int    `json:"category" valid:"Min(0)"`
	Codename     string  `json:"codename" valid:"Required;MaxSize(100);Match(/^[a-zA-Z][a-zA-Z0-9._]{1,}$/)"`
	FormalParams string  `json:"formal_params" valid:"Required;MinSize(2)"`
	Disabled     *bool   `json:"disabled" gorm:"default:0" valid:"Required"`
	Description  *string `json:"description" valid:"MinSize(0);MaxSize(500)"`
}

// Valid
func (s *SetTask) Valid(v *validation.Validation) {
	if ct, _ := dao.GetTaskCount(dao.Query{"codename=?": s.Codename, "id<>?": s.taskId}); ct > 0 {
		_ = v.SetError("Name", "已有相同代号的任务存在")
	}
}

// Validate
func (s *SetTask) Validate(taskId int) *message.Message {
	// 验证角色是否存在
	if ct, _ := dao.GetTaskCount(dao.Query{"id=?": taskId}); ct < 1 {
		return msg.NotFoundFailed.SetDetail("任务不存在")
	}

	s.taskId = taskId
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(s)
	if err != nil {
		return msg.ValidateFailed.SetError(err)
	}
	// 数据验证未通过
	if !ok {
		return msg.ValidateFailed.SetError(valid.Errors)
	}

	return nil
}

// ListTasksQuery struct
type ListTasksQuery struct {
	Query    *string `form:"query"`
	Disabled *bool   `form:"disabled"`
}

// UpdateQuery
func (q *ListTasksQuery) UpdateQuery(query dao.Query) error {
	// 通用Query
	if q.Query != nil && *q.Query != "" {
		likeQuery := fmt.Sprintf("%%%s%%", *q.Query)
		query["(codename LIKE @query OR description LIKE @query OR id LIKE @query)"] = sql.Named("query", likeQuery)
	}

	if q.Disabled != nil {
		query["disabled=?"] = *q.Disabled
	}

	return nil
}

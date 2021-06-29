package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

const (
	BUTILIN_TASK_CATEGORY = 0   // 内置任务
	BASH_TASK_CATEGORY    = 100 // Bash任务
	PYTHON_TASK_CATEGORY  = 101 // Python任务
)

type CallTask struct {
	Timeout   *int64 `json:"timeout" binding:"required"`
	Arguments string `json:"arguments" binding:"required,json"`
}

type Task struct {
	Id          int              `json:"id" swaggerignore:"true"`
	Category    *int             `json:"category" binding:"required"`
	Codename    string           `json:"codename" binding:"required"`
	Description *string          `json:"description" binding:"required"`
	Arguments   string           `json:"arguments" binding:"required,json"`
	Disabled    *bool            `json:"disabled" gorm:"default:0" binding:"required"`
	CreatedAt   *utils.LocalTime `json:"created_at" swaggerignore:"true"`
	CreatedBy   string           `json:"created_by" swaggerignore:"true"`
	UpdatedAt   *utils.LocalTime `json:"updated_at" swaggerignore:"true"`
	UpdatedBy   *string          `json:"updated_by" gorm:"default:''" swaggerignore:"true"`
}

// NewTask 新建任务
func NewTask(disabled *bool, category *int, codename, description, arguments, createdBy string) *Task {
	var t = Task{
		Category:    category,
		Codename:    codename,
		Description: &description,
		Arguments:   arguments,
		Disabled:    disabled,
		CreatedBy:   createdBy,
	}

	if res := db.Create(&t); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"category":    category,
			"codename":    codename,
			"description": description,
			"arguments":   arguments,
			"disabled":    disabled,
			"error":       res.Error.Error(),
		}, "Error in model.NewTask.")
		return nil
	}

	return &t
}

// RemoveTask 删除任务
func RemoveTask(taskId int) bool {
	res := db.Delete(Task{}, "id=?", taskId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    taskId,
			"error": res.Error.Error(),
		}, "Error in model.RemoveTask.")
		return false
	}

	return true
}

// SetTask 更新任务
func SetTask(id int, disabled *bool, category *int, codename, description, arguments, updatedBy string) (*Task, bool) {
	var t = Task{}

	res := db.Model(&Task{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"category":    category,
			"codename":    codename,
			"description": description,
			"arguments":   arguments,
			"disabled":    disabled,
			"updated_by":  updatedBy,
		}).
		First(&t)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":          id,
			"category":    category,
			"codename":    codename,
			"description": description,
			"arguments":   arguments,
			"disabled":    disabled,
			"updated_by":  updatedBy,
			"error":       res.Error.Error(),
		}, "Error in model.SetTask.")
		return nil, false
	}

	return &t, true
}

// GetTask 查询单个任务
func GetTask(query Query) (*Task, bool) {
	var (
		t = Task{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&t); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.GetTask.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.GetTask.")
		return nil, false
	}

	return &t, true
}

// ListTasks 查询任务
func ListTasks(query Query) (*[]Task, bool) {
	var d = db
	ts := make([]Task, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ts); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.ListTasks.")
			return &ts, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.ListTasks.")
		return nil, false
	}

	return &ts, true
}

// PagedListTasks 查询任务-分页
func PagedListTasks(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&Task{})
	ts := make([]Task, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ts)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in model.PagedListTasks.")
		return nil, false
	}

	return pg, true
}

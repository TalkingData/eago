package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/task/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewTask 新建任务
func NewTask(disabled bool, category int, codename, description, fParams, createdBy string) *model.Task {
	var t = model.Task{
		Category:     &category,
		Codename:     codename,
		Description:  &description,
		FormalParams: fParams,
		Disabled:     &disabled,
		CreatedBy:    createdBy,
	}

	if res := db.Create(&t); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"category":      category,
			"codename":      codename,
			"formal_params": fParams,
			"disabled":      disabled,
			"description":   description,
			"error":         res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &t
}

// RemoveTask 删除任务
func RemoveTask(taskId int) bool {
	res := db.Delete(model.Task{}, "id=?", taskId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    taskId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetTask 更新任务
func SetTask(id int, disabled bool, category int, codename, description, fParams, updatedBy string) (*model.Task, bool) {
	var t = model.Task{}

	res := db.Model(&model.Task{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"category":      category,
			"codename":      codename,
			"formal_params": fParams,
			"disabled":      disabled,
			"description":   description,
			"updated_by":    updatedBy,
		}).
		First(&t)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":            id,
			"category":      category,
			"codename":      codename,
			"formal_params": fParams,
			"disabled":      disabled,
			"description":   description,
			"updated_by":    updatedBy,
			"error":         res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &t, true
}

// GetTask 查询单个任务
func GetTask(query Query) (*model.Task, bool) {
	var (
		d = db
		t = model.Task{}
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&t); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &t, true
}

// GetTaskCount 查询任务数量
func GetTaskCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Task{})

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Count(&count); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return count, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Count.")
		return count, false
	}
	return count, true
}

// ListTasks 查询任务
func ListTasks(query Query) (*[]model.Task, bool) {
	var d = db
	ts := make([]model.Task, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ts); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &ts, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &ts, true
}

// PagedListTasks 查询任务-分页
func PagedListTasks(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Task{})
	ts := make([]model.Task, 0)

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
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

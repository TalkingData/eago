package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/task/model"

	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewSchedule 新建计划任务
func NewSchedule(tCodeName, expr, args, description string, timeout int, disabled bool, createdBy string) *model.Schedule {
	var sch = model.Schedule{
		TaskCodename: tCodeName,
		Expression:   expr,
		Description:  &description,
		Timeout:      &timeout,
		Arguments:    args,
		Disabled:     &disabled,
		CreatedBy:    createdBy,
	}

	if res := db.Create(&sch); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"task_codename": tCodeName,
			"expression":    expr,
			"timeout":       timeout,
			"arguments":     args,
			"disabled":      disabled,
			"description":   description,
			"error":         res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &sch
}

// RemoveSchedule 删除计划任务
func RemoveSchedule(schId int) bool {
	res := db.Delete(model.Schedule{}, "id=?", schId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    schId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetSchedule 更新计划任务
func SetSchedule(id int, tCodeName, expr, args, description string, timeout int64, disabled bool, updatedBy string) (*model.Schedule, bool) {
	var sch = model.Schedule{}

	res := db.Model(&model.Schedule{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"task_codename": tCodeName,
			"expression":    expr,
			"timeout":       timeout,
			"arguments":     args,
			"disabled":      disabled,
			"description":   description,
			"updated_by":    updatedBy,
		}).
		First(&sch)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":            id,
			"task_codename": tCodeName,
			"expression":    expr,
			"timeout":       timeout,
			"arguments":     args,
			"disabled":      disabled,
			"updated_by":    updatedBy,
			"description":   description,
			"error":         res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &sch, true
}

// GetSchedule 查询单个计划任务
func GetSchedule(query Query) (*model.Schedule, bool) {
	var (
		d   = db
		sch = model.Schedule{}
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&sch); res.Error != nil {
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

	return &sch, true
}

// GetScheduleCount 查询计划任务数量
func GetScheduleCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Schedule{})

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

// ListSchedules 查询计划任务
func ListSchedules(query Query) (*[]model.Schedule, bool) {
	var d = db
	ss := make([]model.Schedule, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ss); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &ss, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &ss, true
}

// PagedListSchedules 查询计划任务-分页
func PagedListSchedules(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Schedule{})
	ss := make([]model.Schedule, pageSize)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ss)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

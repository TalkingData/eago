package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewTrigger 创建触发器
func NewTrigger(name, description, taskCodename, args, createdBy string) *model.Trigger {
	t := model.Trigger{
		Name:         name,
		Description:  &description,
		TaskCodename: taskCodename,
		Arguments:    args,
		CreatedBy:    createdBy,
	}

	if res := db.Create(&t); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":             name,
			"description":      description,
			"task_codename":    taskCodename,
			"arguments_length": len(args),
			"error":            res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &t
}

// RemoveTrigger 删除触发器
func RemoveTrigger(tId int) bool {
	res := db.Delete(model.Trigger{}, "id=?", tId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    tId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetTrigger 更新触发器
func SetTrigger(id int, name, description, taskCodename, args, updatedBy string) (*model.Trigger, bool) {
	var t model.Trigger

	res := db.Model(&model.Trigger{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":          name,
			"description":   description,
			"task_codename": taskCodename,
			"arguments":     args,
			"updated_by":    updatedBy,
		}).
		First(&t)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":               id,
			"name":             name,
			"description":      description,
			"task_codename":    taskCodename,
			"arguments_length": len(args),
			"updated_by":       updatedBy,
			"error":            res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &t, true
}

// GetTrigger 查询单个触发器
func GetTrigger(query Query) (*model.Trigger, bool) {
	var (
		t = model.Trigger{}
		d = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&t); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.First")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First")
		return nil, false
	}

	return &t, true
}

// GetTriggerCount 查询触发器数量
func GetTriggerCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Trigger{})

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

// ListTriggers 查询触发器
func ListTriggers(query Query) ([]model.Trigger, bool) {
	var d = db
	ts := make([]model.Trigger, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ts); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.Find.")
			return ts, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return ts, true
}

// PagedListTriggers 查询触发器-分页
func PagedListTriggers(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Trigger{})
	ts := make([]model.Trigger, 0)

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

// ListTriggerNodes 关联表操作::列出触发器所关联节点
func ListTriggerNodes(tId int) ([]model.TriggersNode, bool) {
	log.Info("dao.ListTriggerNodes called.")
	defer log.Info("dao.ListTriggerNodes end.")

	var d = db.Model(&model.Node{})
	tns := make([]model.TriggersNode, 0)

	res := d.Select("nodes.id AS id, "+
		"nodes.name AS name, "+
		"nodes.parent_id AS parent_id").
		Joins("LEFT JOIN node_triggers AS nt ON nodes.id = nt.node_id").
		Where("nt.trigger_id=?", tId).
		Find(&tns)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error,
			}, "Record not found")
			return tns, true
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return tns, true
}

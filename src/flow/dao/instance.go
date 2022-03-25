package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/flow/conf"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewInstance 创建流程实例
func NewInstance(formId, status int, name, formData, flowChain, createdBy string) *model.Instance {
	log.Info("dao.NewInstance called.")
	defer log.Info("dao.NewInstance end.")

	// 保证流程实例名称不超过表最大长度
	if len(name) > conf.INSTANCE_NAME_MAX_LENGTH {
		name = name[:conf.INSTANCE_NAME_MAX_LENGTH]
	}

	i := model.Instance{
		Name:      name,
		Status:    status,
		FormId:    formId,
		FormData:  &formData,
		FlowChain: &flowChain,
		CreatedBy: createdBy,
	}

	if res := db.Create(&i); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":       name,
			"status":     status,
			"form_id":    formId,
			"form_data":  formData,
			"flow_chain": flowChain,
			"created_by": createdBy,
			"error":      res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &i
}

// SetInstance 设置流程实例
func SetInstance(id, status, step, assReq int, flowChain, currAss, passedAss, updatedBy string) bool {
	log.Info("dao.SetInstance called.")
	defer log.Info("dao.SetInstance end.")

	res := db.Model(&model.Instance{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status":             status,
			"current_step":       step,
			"flow_chain":         flowChain,
			"assignees_required": assReq,
			"current_assignees":  currAss,
			"passed_assignees":   passedAss,
			"updated_by":         updatedBy,
		})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":                 id,
			"status":             status,
			"current_step":       step,
			"flow_chain":         flowChain,
			"assignees_required": assReq,
			"current_assignees":  currAss,
			"passed_assignees":   passedAss,
			"updated_by":         updatedBy,
			"error":              res.Error,
		}, "An error occurred while db.Where.Updates.")
		return false
	}

	return true
}

// SetHandleInstance 设置流程实例
func SetHandleInstance(id, status, step, assReq int, formData, currAss, passedAss, updatedBy string) bool {
	log.Info("dao.SetInstance called.")
	defer log.Info("dao.SetInstance end.")

	res := db.Model(&model.Instance{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"status":             status,
			"current_step":       step,
			"form_data":          formData,
			"assignees_required": assReq,
			"current_assignees":  currAss,
			"passed_assignees":   passedAss,
			"updated_by":         updatedBy,
		})
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":                 id,
			"status":             status,
			"current_step":       step,
			"form_data":          formData,
			"assignees_required": assReq,
			"current_assignees":  currAss,
			"passed_assignees":   passedAss,
			"updated_by":         updatedBy,
			"error":              res.Error,
		}, "An error occurred while db.Where.Updates.")
		return false
	}

	return true
}

// GetInstance 查询单个流程实例
func GetInstance(query Query) (*model.Instance, error) {
	log.Info("dao.GetInstance called.")
	defer log.Info("dao.GetInstance end.")

	var (
		ins = model.Instance{}
		d   = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&ins); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.First.")
			return nil, nil
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, res.Error
	}

	return &ins, nil
}

// GetInstancesCount 查询实例数量
func GetInstancesCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Instance{})

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

// PagedListInstances 查询流程实例-分页
func PagedListInstances(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	log.Info("dao.PagedListInstances called.")
	defer log.Info("dao.PagedListInstances end.")

	var d = db.Model(&model.Instance{})
	ins := make([]model.Instance, pageSize)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ins)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

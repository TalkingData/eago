package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewFlow 创建流程
func NewFlow(name, insTit string, catId *int, description string, disabled bool, frmID, firstNodeID int, createdBy string) *model.Flow {
	f := model.Flow{
		Name:          name,
		InstanceTitle: insTit,
		CategoriesId:  catId,
		Disabled:      &disabled,
		Description:   &description,
		FormId:        frmID,
		FirstNodeId:   firstNodeID,
		CreatedBy:     createdBy,
	}

	if res := db.Create(&f); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":           name,
			"instance_title": insTit,
			"categories_id":  catId,
			"disabled":       disabled,
			"description":    description,
			"form_id":        frmID,
			"first_node_id":  firstNodeID,
			"error":          res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &f
}

// RemoveFlow 删除流程
func RemoveFlow(fID int) bool {
	res := db.Delete(model.Flow{}, "id=?", fID)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    fID,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetFlow 更新流程
func SetFlow(id int, name, insTit string, catId *int, description string, disabled bool, frmID, firstNodeID int, updatedBy string) (*model.Flow, bool) {
	var f model.Flow

	res := db.Model(&model.Flow{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":           name,
			"instance_title": insTit,
			"categories_id":  catId,
			"disabled":       disabled,
			"description":    description,
			"form_id":        frmID,
			"first_node_id":  firstNodeID,
			"updated_by":     updatedBy,
		}).
		First(&f)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":            id,
			"name":          name,
			"categories_id": catId,
			"disabled":      disabled,
			"description":   description,
			"form_id":       frmID,
			"first_node_id": firstNodeID,
			"updated_by":    updatedBy,
			"error":         res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &f, true
}

// GetFlow 查询单个流程
func GetFlow(query Query) (*model.Flow, bool) {
	log.Info("dao.GetFlow called.")
	defer log.Info("dao.GetFlow end.")

	var (
		d = db
		f model.Flow
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&f); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.First.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &f, true
}

// GetFlowCount 查询流程数量
func GetFlowCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Flow{})

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

// ListFlows 查询流程
func ListFlows(query Query) ([]model.Flow, bool) {
	var d = db
	fs := make([]model.Flow, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&fs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.Find.")
			return fs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return fs, true
}

// PagedListFlows 查询流程-分页
func PagedListFlows(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Flow{}).
		Select("flows.id AS id, " +
			"flows.name AS name, " +
			"flows.instance_title AS instance_title, " +
			"flows.categories_id AS categories_id, " +
			"c.name AS categories_name, " +
			"flows.disabled AS disabled, " +
			"flows.description AS description, " +
			"flows.form_id AS form_id, " +
			"f.name AS form_name, " +
			"flows.first_node_id AS first_node_id, " +
			"n.name AS first_node_name, " +
			"flows.created_at AS created_at, " +
			"flows.created_by AS created_by, " +
			"flows.updated_at AS updated_at, " +
			"flows.updated_by AS updated_by").
		Joins("LEFT JOIN categories AS c ON c.id = flows.categories_id").
		Joins("LEFT JOIN forms AS f ON f.id = flows.form_id").
		Joins("LEFT JOIN nodes AS n ON n.id = flows.first_node_id")
	fs := make([]model.ListFlows, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &fs)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

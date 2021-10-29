package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewForm 创建表单
func NewForm(name string, disabled bool, description, body, createdBy string) *model.Form {
	f := model.Form{
		Name:        name,
		Disabled:    &disabled,
		Description: &description,
		Body:        &body,
		CreatedBy:   createdBy,
	}

	if res := db.Create(&f); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":        name,
			"disabled":    disabled,
			"description": description,
			"body_length": len(body),
			"error":       res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &f
}

// RemoveForm 删除表单
func RemoveForm(frmId int) bool {
	res := db.Delete(model.Form{}, "id=?", frmId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    frmId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetForm 更新表单
func SetForm(id int, name string, disabled bool, description, updatedBy string) (*model.Form, bool) {
	var f model.Form

	res := db.Model(&model.Form{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"disabled":    disabled,
			"description": description,
			"updated_by":  updatedBy,
		}).
		First(&f)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":          id,
			"name":        name,
			"disabled":    disabled,
			"description": description,
			"updated_by":  updatedBy,
			"error":       res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &f, true
}

// GetForm 查询单个表单
func GetForm(query Query) (*model.Form, bool) {
	var (
		d = db
		f = model.Form{}
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
		}, "An error occurred while b.Where.First.")
		return nil, false
	}

	return &f, true
}

// GetFormCount 查询表单数量
func GetFormCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Form{})

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

// ListForms 查询表单
func ListForms(query Query) ([]model.Form, bool) {
	var d = db
	fs := make([]model.Form, 0)

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

// PagedListForms 查询表单-分页
func PagedListForms(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.Form{})
	fs := make([]model.FormWithoutBody, 0)

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

package dao

import (
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewCategory 创建类别
func NewCategory(name, createdBy string) *model.Categories {
	c := model.Categories{
		Name:      name,
		CreatedBy: createdBy,
	}

	if res := db.Create(&c); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"name":  name,
			"error": res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &c
}

// RemoveCategory 删除类别
func RemoveCategory(tId int) bool {
	res := db.Delete(model.Categories{}, "id=?", tId)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    tId,
			"error": res.Error,
		}, "An error occurred while db.Delete.")
		return false
	}

	return true
}

// SetCategory 更新类别
func SetCategory(id int, name, updatedBy string) (*model.Categories, bool) {
	var t = model.Categories{}

	res := db.Model(&model.Categories{}).
		Where("id=?", id).
		Updates(map[string]interface{}{
			"name":       name,
			"updated_by": updatedBy,
		}).
		First(&t)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":         id,
			"name":       name,
			"updated_by": updatedBy,
			"error":      res.Error,
		}, "An error occurred while db.Model.Where.Updates.First.")
		return nil, false
	}

	return &t, true
}

// GetCategory 查询单个类别
func GetCategory(query Query) (*model.Categories, bool) {
	var (
		d = db
		t = model.Categories{}
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&t); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.First.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.First.")
		return nil, false
	}

	return &t, true
}

// GetCategoriesCount 查询类别数量
func GetCategoriesCount(query Query) (count int64, ok bool) {
	d := db.Model(&model.Categories{})

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

// ListCategories 查询类别
func ListCategories(query Query) ([]model.Categories, bool) {
	var d = db
	cs := make([]model.Categories, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&cs); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.Find.")
			return cs, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return cs, true
}

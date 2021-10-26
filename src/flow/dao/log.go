package dao

import (
	"eago/common/log"
	"eago/flow/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewLog 新增审批日志
func NewLog(insId int, result bool, content, createdBy string) *model.Log {
	l := model.Log{
		InstanceId: insId,
		Result:     result,
		Content:    &content,
		CreatedBy:  createdBy,
	}

	if res := db.Create(&l); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"instance_id": insId,
			"result":      result,
			"content":     content,
			"created_by":  createdBy,
			"error":       res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &l
}

// ListLogs 查询审批日志
func ListLogs(query Query) ([]model.Log, bool) {
	var d = db
	ls := make([]model.Log, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ls); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found in db.Where.Find.")
			return ls, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return ls, true
}

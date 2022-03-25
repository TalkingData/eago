package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/task/conf"
	"eago/task/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// NewLog 新建日志
func NewLog(partition string, resultId int, content *string) *model.Log {
	var (
		d  = db.Table(GetLogTableNameByPartition(partition))
		lg = model.Log{
			ResultId: resultId,
			Content:  *content,
		}
	)

	// 创建记录
	if res := d.Create(&lg); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &lg
}

// PagedListLogsByPartition 列出日志（需指定分区）-分页
func PagedListLogsByPartition(query Query, partition string, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var (
		tbName = GetLogTableNameByPartition(partition)
		d      = db.Table(tbName)
	)
	ls := make([]model.Log, pageSize)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &ls)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"table_name": tbName,
			"query":      fmt.Sprintf("%v", query),
			"error":      err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

// ListLogs 查询任务
func ListLogs(query Query, partition string) (*[]model.Log, bool) {
	var (
		tbName = GetLogTableNameByPartition(partition)
		d      = db.Table(tbName)
	)
	ls := make([]model.Log, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ls); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &ls, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &ls, true
}

// GetLogTableNameByPartition 按分区获得日志表名（需指定分区）
func GetLogTableNameByPartition(partition string) string {
	return fmt.Sprintf("logs_%s", partition)
}

// GetLogTableNameByTime 按时间获得日志表名（需指定分区）
func GetLogTableNameByTime(t *time.Time) string {
	return GetLogTableNameByPartition(t.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT))
}

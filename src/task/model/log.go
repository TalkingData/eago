package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"eago/task/conf"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	Id        int              `json:"id" gorm:"type:int(11) NOT NULL AUTO_INCREMENT;primaryKey" swaggerignore:"true"`
	ResultId  int              `json:"result_id" gorm:"type:int(11) NOT NULL;index"`
	Content   string           `json:"content" gorm:"type:varchar(2000) NOT NULL;default:''"`
	CreatedAt *utils.LocalTime `json:"created_at" gorm:"type:datetime NOT NULL" swaggerignore:"true"`
}

// NewLog 新建日志
func NewLog(partition string, resultId int, content *string) *Log {
	var (
		d = db.Table(GetLogTableNameByPartition(partition))
		l = Log{
			ResultId: resultId,
			Content:  *content,
		}
	)

	// 创建记录
	if res := d.Create(&l); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"error":     res.Error.Error(),
		}, "Error in model.NewLog.")
		return nil
	}

	return &l
}

// PagedListLogsByPartition 列出日志（需指定分区）-分页
func PagedListLogsByPartition(query Query, partition string, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var (
		tableName = GetLogTableNameByPartition(partition)
		d         = db.Table(tableName)
	)
	ls := make([]Log, 0)

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
			"table_name": tableName,
			"query":      fmt.Sprintf("%v", query),
			"error":      err.Error(),
		}, "Error in model.PagedListLogsByPartition.")
		return nil, false
	}

	return pg, true
}

// ListLogs 查询任务
func ListLogs(query Query, partition string) (*[]Log, bool) {
	var (
		tableName = GetLogTableNameByPartition(partition)
		d         = db.Table(tableName)
	)
	ls := make([]Log, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.Find(&ls); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.ListLogs.")
			return &ls, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.ListLogs.")
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

// TableName 获取数据库表名
func (l *Log) TableName() string {
	return "_tmp_logs"
}

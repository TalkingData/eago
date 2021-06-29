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

type Result struct {
	Id           int              `json:"id" gorm:"type:int(11) NOT NULL AUTO_INCREMENT;primaryKey" swaggerignore:"true"`
	TaskCodename string           `json:"task_codename" gorm:"type:varchar(100);not null;index"`
	Status       int              `json:"status" gorm:"type:int(11) NOT NULL;index"`
	Caller       string           `json:"caller" gorm:"type:varchar(100) NOT NULL;index"`
	Worker       string           `json:"worker" gorm:"type:varchar(100) NOT NULL;default:'';index"`
	Timeout      *int64           `json:"timeout" gorm:"type:int(11) NOT NULL"`
	Arguments    string           `json:"arguments" gorm:"type:varchar(2000) NOT NULL;default:'{}'" swaggerignore:"true"`
	StartAt      *utils.LocalTime `json:"start_at" gorm:"type:datetime NOT NULL;index" swaggerignore:"true"`
	EndAt        *utils.LocalTime `json:"end_at" gorm:"type:datetime" swaggerignore:"true"`
}

// NewResult 新建结果
func NewResult(taskCodename, caller, arguments string, timeout int64, status int) *Result {
	var (
		t         = time.Now()
		partition = t.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT)
		d         = db.Table(GetResultTableNameByPartition(partition))
		r         = Result{
			TaskCodename: taskCodename,
			Caller:       caller,
			Status:       status,
			Timeout:      &timeout,
			Arguments:    arguments,
			StartAt:      &utils.LocalTime{t},
			EndAt:        nil,
		}
	)

	// 检测分区是否存在
	p, ok := GetResultPartition(Query{"partition": partition})
	// 获取分区错误
	if !ok {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
		}, "Error in model.GetResultPartition.")
		return nil
	}

	// 如果获取不到分区，则创建
	if p == nil {
		// 检测创建是否成功
		if nil == NewResultPartitionWithCreateTables(partition) {
			log.ErrorWithFields(log.Fields{
				"partition": partition,
			}, "Error in model.NewResultPartitionWithCreateTables.")
			return nil
		}
	}

	// 创建记录
	if res := d.Create(&r); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"error":     res.Error.Error(),
		}, "Error in model.New.")
		return nil
	}

	return &r
}

// SetResultStatus 更新任务状态
func SetResultStatus(partition string, id int, status int, end bool) bool {
	var d = db.Table(GetResultTableNameByPartition(partition))
	updates := map[string]interface{}{"status": status}

	if end {
		updates["end_at"] = &utils.LocalTime{time.Now()}
	}
	res := d.Where("id=?", id).
		Updates(updates)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in model.SetResultStatus.")
		return false
	}

	return true
}

// SetResultWorker 更新执行器信息
func SetResultWorker(partition string, id int, worker string) bool {
	var d = db.Table(GetResultTableNameByPartition(partition))

	res := d.Where("id=?", id).
		Update("worker", worker)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error.Error(),
		}, "Error in model.SetResultWorker.")
		return false
	}

	return true
}

// GetResult 查询单个结果
func GetResult(partition string, id int) (*Result, bool) {
	var (
		d = db.Table(GetResultTableNameByPartition(partition))
		r = Result{}
	)

	if res := d.Where("id=?", id).First(&r); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"partition": partition,
				"id":        id,
				"error":     res.Error.Error(),
			}, "Record not found in model.GetResult.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"id":        id,
			"error":     res.Error.Error(),
		}, "Error in model.GetResult.")
		return nil, false
	}

	return &r, true
}

// PagedListResultsByPartition 列出结果（需指定分区）-分页
func PagedListResultsByPartition(query Query, partition string, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var (
		tableName = GetResultTableNameByPartition(partition)
		d         = db.Table(tableName)
	)
	rs := make([]Result, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &rs)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"table_name": tableName,
			"query":      fmt.Sprintf("%v", query),
			"error":      err.Error(),
		}, "Error in model.PagedListResultsByPartition.")
		return nil, false
	}

	return pg, true
}

// GetResultTableNameByPartition 按分区获得结果表名（需指定分区）
func GetResultTableNameByPartition(partition string) string {
	return fmt.Sprintf("results_%s", partition)
}

// GetResultTableNameByTime 按时间获得结果表名（需指定分区）
func GetResultTableNameByTime(t *time.Time) string {
	return GetResultTableNameByPartition(t.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT))
}

// TableName 获取数据库表名
func (r *Result) TableName() string {
	return "_tmp_results"
}

func (r *Result) GetPartition() (string, error) {
	if r.StartAt == nil {
		return "", fmt.Errorf("Null result object, has no partition.")
	}

	return r.StartAt.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT), nil
}

package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/common/utils"
	"eago/task/conf"
	"eago/task/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// NewResult 新建结果
func NewResult(taskCodename, caller, arguments string, timeout int64, status int) *model.Result {
	var (
		t    = time.Now()
		part = t.Format(conf.TASK_PARTITION_TIMESTAMP_FORMAT)
		d    = db.Table(GetResultTableNameByPartition(part))
		r    = model.Result{
			TaskCodename: taskCodename,
			Caller:       caller,
			Status:       status,
			Timeout:      &timeout,
			Arguments:    arguments,
			StartAt:      &utils.LocalTime{Time: t},
			EndAt:        nil,
		}
	)

	// 检测分区是否存在
	p, ok := GetResultPartition(Query{"partition": part})
	// 获取分区错误
	if !ok {
		log.ErrorWithFields(log.Fields{
			"partition": part,
		}, "An error occurred while GetResultPartition.")
		return nil
	}

	// 如果获取不到分区，则创建
	if p == nil {
		// 检测创建是否成功
		if nil == NewResultPartitionWithCreateTables(part) {
			log.ErrorWithFields(log.Fields{
				"partition": part,
			}, "An error occurred while NewResultPartitionWithCreateTables.")
			return nil
		}
	}

	// 创建记录
	if res := d.Create(&r); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": part,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &r
}

// SetResultStatus 更新任务状态
func SetResultStatus(partition string, id int, status int, end bool) bool {
	var d = db.Table(GetResultTableNameByPartition(partition))
	updates := map[string]interface{}{"status": status}

	if end {
		updates["end_at"] = &utils.LocalTime{Time: time.Now()}
	}
	res := d.Where("id=?", id).
		Updates(updates)
	if res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"id":    id,
			"error": res.Error,
		}, "An error occurred while db.Where.Updates.")
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
			"error": res.Error,
		}, "An error occurred while db.Where.Updates.")
		return false
	}

	return true
}

// GetResult 查询单个结果
func GetResult(partition string, id int) (*model.Result, bool) {
	var (
		r = model.Result{}
		d = db.Table(GetResultTableNameByPartition(partition))
	)

	if res := d.Where("id=?", id).First(&r); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"partition": partition,
				"id":        id,
				"error":     res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"id":        id,
			"error":     res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &r, true
}

// PagedListResultsByPartition 列出结果（需指定分区）-分页
func PagedListResultsByPartition(query Query, partition string, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var (
		tbName = GetResultTableNameByPartition(partition)
		d      = db.Table(tbName)
	)
	rs := make([]model.Result, 0)

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
			"table_name": tbName,
			"query":      fmt.Sprintf("%v", query),
			"error":      err,
		}, "An error occurred while pagination.GormPaging.")
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

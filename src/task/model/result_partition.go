package model

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ResultPartition struct {
	Id        int    `json:"id"`
	Partition string `json:"partition" binding:"required"`
}

// NewResultPartition 新建结果分区
func NewResultPartition(partition string) *ResultPartition {
	var rp = ResultPartition{
		Partition: partition,
	}

	// 创建记录
	if res := db.Create(&rp); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"error":     res.Error.Error(),
		}, "Error in model.NewResultPartition.")
		return nil
	}

	return &rp
}

// NewResultPartitionWithCreateTables 新建结果分区并建立结果表和日志表
func NewResultPartitionWithCreateTables(partition string) *ResultPartition {
	var (
		r  = &Result{}
		l  = &Log{}
		rp = ResultPartition{
			Partition: partition,
		}
	)
	tx := db.Begin()

	// 创建记录
	if res := tx.Create(&rp); res.Error != nil {
		tx.Rollback()
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"action":    "Create",
			"error":     res.Error.Error(),
		}, "Error in model.NewResultPartitionWithCreateTables.")
		return nil
	}

	// 1. 创建Result表
	// 判断默认名称的Result表是否存在
	if !tx.Migrator().HasTable(r) {
		// 如果不存在，则创建默认名称的Result表
		err := tx.Migrator().CreateTable(r)
		if err != nil {
			tx.Rollback()
			log.ErrorWithFields(log.Fields{
				"partition": partition,
				"action":    "CreateTable",
				"table":     "Result",
				"error":     err.Error(),
			}, "Error in model.NewResultPartitionWithCreateTables.")
			return nil
		}
	}

	// 将建默认名称的Result表改名为当前新建分区对应的表Result名
	if err := tx.Migrator().RenameTable(r.TableName(), GetResultTableNameByPartition(partition)); err != nil {
		tx.Rollback()
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"action":    "RenameTable",
			"table":     "Result",
			"error":     err.Error(),
		}, "Error in model.NewResultPartitionWithCreateTables.")
		return nil
	}

	// 2. 创建Log表
	// 判断默认名称的Result表是否存在
	if !tx.Migrator().HasTable(l) {
		// 如果不存在，则创建默认名称的Result表
		if err := tx.Migrator().CreateTable(l); err != nil {
			tx.Rollback()
			log.ErrorWithFields(log.Fields{
				"partition": partition,
				"action":    "CreateTable",
				"table":     "Logs",
				"error":     err.Error(),
			}, "Error in model.NewResultPartitionWithCreateTables.")
			return nil
		}
	}

	// 将建默认名称的Result表改名为当前新建分区对应的表Result名
	if err := tx.Migrator().RenameTable(l.TableName(), GetLogTableNameByPartition(partition)); err != nil {
		tx.Rollback()
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"action":    "RenameTable",
			"table":     "Logs",
			"error":     err.Error(),
		}, "Error in model.NewResultPartitionWithCreateTables.")
		return nil
	}

	tx.Commit()
	return &rp
}

// GetResultPartition 查询单个结果分区
func GetResultPartition(query Query) (*ResultPartition, bool) {
	var (
		rt = ResultPartition{}
		d  = db
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&rt); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found in model.GetResultPartition.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.GetResultPartition.")
		return nil, false
	}

	return &rt, true
}

// GetResultPartitionsPartition 获得结果分区名
func GetResultPartitionsPartition(id int) (string, bool) {
	rt, ok := GetResultPartition(Query{"id=?": id})
	if ok && rt != nil {

		return rt.Partition, true
	}

	return "", false
}

// ListResultPartitions 查询结果分区-分页
func ListResultPartitions(query Query) (*[]ResultPartition, bool) {
	var d = db
	rp := make([]ResultPartition, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Find(&rp); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error.Error(),
			}, "Record not found on model.ListResultPartitions.")
			return &rp, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error.Error(),
		}, "Error in model.ListResultPartitions.")
		return nil, false
	}

	return &rp, true
}

// PagedListUsers 查询结果分区-分页
func PagedListResultPartitions(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&ResultPartition{})
	rts := make([]ResultPartition, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &rts)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err.Error(),
		}, "Error in model.PagedListResultPartitions.")
		return nil, false
	}

	return pg, true
}

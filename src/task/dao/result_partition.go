package dao

import (
	"eago/common/api-suite/pagination"
	"eago/common/log"
	"eago/task/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// NewResultPartition 新建结果分区
func NewResultPartition(partition string) *model.ResultPartition {
	rPart := model.ResultPartition{
		Partition: partition,
	}

	// 创建记录
	if res := db.Create(&rPart); res.Error != nil {
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"error":     res.Error,
		}, "An error occurred while db.Create.")
		return nil
	}

	return &rPart
}

// NewResultPartitionWithCreateTables 新建结果分区并建立结果表和日志表
func NewResultPartitionWithCreateTables(partition string) *model.ResultPartition {
	var (
		r     = new(model.Result)
		lg    = new(model.Log)
		rPart = model.ResultPartition{
			Partition: partition,
		}
	)
	tx := db.Begin()

	// 创建记录
	if res := tx.Create(&rPart); res.Error != nil {
		tx.Rollback()
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"action":    "Create",
			"error":     res.Error,
		}, "An error occurred while db.Create.")
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
				"error":     err,
			}, "An error occurred while x.Migrator.CreateTable.")
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
			"error":     err,
		}, "An error occurred while tx.Migrator.RenameTable.")
		return nil
	}

	// 2. 创建Log表
	// 判断默认名称的Log表是否存在
	if !tx.Migrator().HasTable(lg) {
		// 如果不存在，则创建默认名称的Result表
		if err := tx.Migrator().CreateTable(lg); err != nil {
			tx.Rollback()
			log.ErrorWithFields(log.Fields{
				"partition": partition,
				"action":    "CreateTable",
				"table":     "Logs",
				"error":     err,
			}, "An error occurred while tx.Migrator.CreateTable.")
			return nil
		}
	}

	// 将建默认名称的Log表改名为当前新建分区对应的表Log名
	if err := tx.Migrator().RenameTable(lg.TableName(), GetLogTableNameByPartition(partition)); err != nil {
		tx.Rollback()
		log.ErrorWithFields(log.Fields{
			"partition": partition,
			"action":    "RenameTable",
			"table":     "Logs",
			"error":     err,
		}, "An error occurred while tx.Migrator.RenameTable.")
		return nil
	}

	tx.Commit()
	return &rPart
}

// GetResultPartition 查询单个结果分区
func GetResultPartition(query Query) (*model.ResultPartition, bool) {
	var (
		d     = db
		rPart model.ResultPartition
	)

	for k, v := range query {
		d = d.Where(k, v)
	}
	if res := d.First(&rPart); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return nil, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.First.")
		return nil, false
	}

	return &rPart, true
}

// GetResultPartitionsPartition 获得结果分区名
func GetResultPartitionsPartition(id int) (string, bool) {
	rPart, ok := GetResultPartition(Query{"id=?": id})
	if ok && rPart != nil {

		return rPart.Partition, true
	}

	return "", false
}

// ListResultPartitions 查询结果分区-分页
func ListResultPartitions(query Query) (*[]model.ResultPartition, bool) {
	var d = db
	rParts := make([]model.ResultPartition, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}

	if res := d.Find(&rParts); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"query": fmt.Sprintf("%v", query),
				"error": res.Error,
			}, "Record not found.")
			return &rParts, true
		}
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": res.Error,
		}, "An error occurred while db.Where.Find.")
		return nil, false
	}

	return &rParts, true
}

// PagedListResultPartitions 查询结果分区-分页
func PagedListResultPartitions(query Query, page, pageSize int, orderBy ...string) (*pagination.Paginator, bool) {
	var d = db.Model(&model.ResultPartition{})
	rParts := make([]model.ResultPartition, 0)

	for k, v := range query {
		d = d.Where(k, v)
	}
	pg, err := pagination.GormPaging(&pagination.GormParams{
		Db:       d,
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
	}, &rParts)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"query": fmt.Sprintf("%v", query),
			"error": err,
		}, "An error occurred while pagination.GormPaging.")
		return nil, false
	}

	return pg, true
}

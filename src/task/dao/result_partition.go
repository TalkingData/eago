package dao

import (
	"context"
	"eago/common/logger"
	"eago/common/orm"
	"eago/task/model"
	"errors"
)

// NewResultPartition 新建结果分区
func (d *Dao) NewResultPartition(ctx context.Context, partition string) (*model.ResultPartition, error) {
	resPart := &model.ResultPartition{
		Partition: partition,
	}

	res := d.getDbWithCtx(ctx).Create(&resPart)
	return resPart, res.Error
}

// NewResultPartitionWithCreateTables 新建结果分区并建立结果表和日志表
func (d *Dao) NewResultPartitionWithCreateTables(ctx context.Context, partition string) (*model.ResultPartition, error) {
	tx := d.getDbWithCtx(ctx).Begin()
	defer tx.Rollback()

	resPart := &model.ResultPartition{
		Partition: partition,
	}

	// 创建记录
	if res := tx.Create(&resPart); res.Error != nil {
		return nil, res.Error
	}

	modelRes := &model.Log{}
	// 创建Log表
	// 判断默认名称的Log表是否存在
	if !tx.Migrator().HasTable(modelRes) {
		// 如果不存在，则创建默认名称的Result表
		err := tx.Migrator().CreateTable(modelRes)
		if err != nil {
			d.lg.ErrorWithFields(logger.Fields{
				"partition": partition,
				"action":    "CreateTable",
				"table":     "Result",
				"error":     err,
			}, "An error occurred while tx.Migrator.CreateTable in dao.NewResultPartitionWithCreateTables.")
			return nil, err
		}
	}

	// 将建默认名称的Log表改名为当前新建分区对应的表Log名
	if err := tx.Migrator().RenameTable(modelRes.TableName(), d.getResultTableNameByPartition(partition)); err != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"partition": partition,
			"action":    "RenameTable",
			"table":     "Log",
			"error":     err,
		}, "An error occurred while tx.Migrator.RenameTable in dao.NewResultPartitionWithCreateTables.")
		return nil, err
	}

	tx.Commit()
	return resPart, nil
}

// GetResultPartition 查询单个结果分区
func (d *Dao) GetResultPartition(ctx context.Context, q orm.Query) (resPart *model.ResultPartition, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Limit(1).Find(&resPart)
	return resPart, res.Error
}

// GetResultPartitionsPartition 获得结果分区名
func (d *Dao) GetResultPartitionsPartition(ctx context.Context, id uint32) (string, error) {
	lPart, err := d.GetResultPartition(ctx, orm.Query{"id=?": id})
	if err != nil {
		return "", err
	}

	if lPart != nil {
		return lPart.Partition, nil
	}

	return "", errors.New("got an nil log partition object")
}

// ListResultPartitions 查询结果分区-分页
func (d *Dao) ListResultPartitions(ctx context.Context, q orm.Query) (resParts []*model.ResultPartition, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&resParts)
	return resParts, res.Error
}

// PagedListResultPartitions 查询结果分区-分页
func (d *Dao) PagedListResultPartitions(
	ctx context.Context, q orm.Query, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	resParts := make([]*model.ResultPartition, pageSize)
	db := q.Where(d.getDbWithCtx(ctx).Model(&model.Task{}))
	return orm.PagingQuery(db, page, pageSize, &resParts, orderBy...)
}

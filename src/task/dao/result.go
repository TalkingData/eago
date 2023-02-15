package dao

import (
	"context"
	"eago/common/logger"
	"eago/common/orm"
	"eago/common/utils"
	"eago/task/model"
	"errors"
	"fmt"
	"time"
)

// NewResult 新建结果
func (d *Dao) NewResult(
	ctx context.Context,
	taskCodename, caller, arguments string, timeout int64, status int32,
) (*model.Result, error) {
	currTime := time.Now()
	modelRes := &model.Result{
		TaskCodename: taskCodename,
		Caller:       caller,
		Status:       status,
		Timeout:      &timeout,
		Arguments:    arguments,
		StartAt:      &utils.CustomTime{Time: currTime},
		EndAt:        nil,
	}

	partitionName := currTime.Format(d.conf.Const.TaskResultPartitionTsFormat)

	// 检测分区是否存在
	resPartObj, err := d.GetResultPartition(ctx, orm.Query{"partition": partitionName})
	// 获取分区错误
	if err != nil {
		d.lg.ErrorWithFields(logger.Fields{
			"partition": partitionName,
			"error":     err,
		}, "An error occurred while dao.GetLogPartition in dao.NewLog.")
		return nil, err
	}

	// 如果获取不到分区，则创建
	if resPartObj == nil || resPartObj.Id < 1 {
		// 检测创建是否成功
		resPartObj, err = d.NewResultPartitionWithCreateTables(ctx, partitionName)
		if err != nil {
			d.lg.ErrorWithFields(logger.Fields{
				"partition": partitionName,
				"error":     err,
			}, "An error occurred while dao.NewResultPartitionWithCreateTables in dao.NewResult.")
			return nil, err
		}
		if resPartObj == nil {
			return nil, errors.New("failed NewResultPartitionWithCreateTables in dao.")
		}
	}

	// 创建结果
	res := d.getDbWithCtx(ctx).Table(d.getResultTableNameByPartition(partitionName)).Create(&modelRes)
	return modelRes, res.Error
}

// SetResultStatus 更新任务状态
func (d *Dao) SetResultStatus(ctx context.Context, partition string, id uint32, status int32, end bool) error {
	updates := map[string]interface{}{"status": status}

	if end {
		updates["end_at"] = &utils.CustomTime{Time: time.Now()}
	}

	res := d.getDbWithCtx(ctx).
		Table(d.getResultTableNameByPartition(partition)).
		Where("id=?", id).
		Updates(updates)

	return res.Error
}

// SetResultWorker 更新执行器信息
func (d *Dao) SetResultWorker(ctx context.Context, partition string, id uint32, worker string) error {
	res := d.getDbWithCtx(ctx).
		Table(d.getResultTableNameByPartition(partition)).
		Where("id=?", id).
		Update("worker", worker)

	return res.Error
}

// GetResult 查询单个结果
func (d *Dao) GetResult(ctx context.Context, partition string, id uint32) (r *model.Result, err error) {
	res := d.getDbWithCtx(ctx).
		Table(d.getResultTableNameByPartition(partition)).
		Where("id=?", id).
		Limit(1).
		Find(&r)
	return r, res.Error
}

// PagedListResultsByPartition 列出结果（需指定分区）-分页
func (d *Dao) PagedListResultsByPartition(
	ctx context.Context,
	q orm.Query, partition string, page, pageSize int, orderBy ...string,
) (*orm.Paginator, error) {
	res := make([]*model.Result, pageSize)

	db := q.Where(d.getDbWithCtx(ctx).Table(d.getResultTableNameByPartition(partition)))
	return orm.PagingQuery(db, page, pageSize, &res, orderBy...)
}

// getResultTableNameByPartition 按分区获得结果表名（需指定分区）
func (d *Dao) getResultTableNameByPartition(partition string) string {
	return fmt.Sprintf("results_%s", partition)
}

// GetResultTableNameByTime 按时间获得结果表名（需指定分区）
func (d *Dao) GetResultTableNameByTime(t *time.Time) string {
	return d.getResultTableNameByPartition(t.Format(d.conf.Const.TaskResultPartitionTsFormat))
}

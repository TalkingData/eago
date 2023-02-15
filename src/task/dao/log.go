package dao

import (
	"context"
	"eago/task/model"
	"fmt"
	"time"
)

// NewLog 新增任务日志
func (d *Dao) NewLog(ctx context.Context, partition string, resultId uint32, content *string) (*model.Log, error) {
	modelLog := &model.Log{
		ResultId: resultId,
		Content:  *content,
	}

	res := d.getDbWithCtx(ctx).Table(d.getLogTableNameByPartition(partition)).Create(&modelLog)
	return modelLog, res.Error
}

// ListLogsByPartition 查询任务日志（需指定分区）
func (d *Dao) ListLogsByPartition(ctx context.Context, partition string, resId uint32) (logs []*model.Log, err error) {
	res := d.getDbWithCtx(ctx).
		Table(d.getLogTableNameByPartition(partition)).
		Where("result_id=?", resId).
		Find(&logs)

	return logs, res.Error
}

// GetLogTableNameByPartition 按分区获得任务日志表名（需指定分区）
func (d *Dao) getLogTableNameByPartition(partition string) string {
	return fmt.Sprintf("logs_%s", partition)
}

// GetLogTableNameByTime 按时间获得任务日志表名（需指定分区）
func (d *Dao) GetLogTableNameByTime(t *time.Time) string {
	return d.getLogTableNameByPartition(t.Format(d.conf.Const.TaskResultPartitionTsFormat))
}

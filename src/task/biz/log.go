package biz

import (
	"context"
	"eago/common/logger"
)

// NewLog 新增任务日志
func (b *Biz) NewLog(ctx context.Context, taskUniqueId string, content *string) error {
	b.logger.DebugWithFields(logger.Fields{
		"task_unique_id": taskUniqueId,
	}, "biz.NewLog called.")
	defer b.logger.Debug("biz.NewLog end.")

	// 将任务唯一Id解码为任务结果Id和分区
	part, resId, err := b.TaskUniqueIdDecode(taskUniqueId)
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while biz.TaskUniqueIdDecode in biz.NewLog.")
		return err
	}

	log, err := b.dao.NewLog(ctx, part, resId, content)
	if err != nil {
		// 创建任务日志错误
		b.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"partition":      part,
			"result_id":      resId,
			"error":          err,
		}, "An error occurred while dao.NewLog in biz.NewLog.")
		return err
	}

	if log == nil || log.Id < 1 {
		// 创建日志时返回了空对象
		b.logger.WarnWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"partition":      part,
			"result_id":      resId,
		}, "An nil object is returned after calling dao.NewLog, This error will be ignored.")
		return nil
	}

	return nil
}

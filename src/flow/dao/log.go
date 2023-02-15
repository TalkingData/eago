package dao

import (
	"context"
	"eago/common/orm"
	"eago/flow/model"
)

// NewLog 新增审批日志
func (d *Dao) NewLog(ctx context.Context, insId uint32, result bool, content, createdBy string) (*model.Log, error) {
	log := &model.Log{
		InstanceId: insId,
		Result:     result,
		Content:    &content,
		CreatedBy:  createdBy,
	}

	res := d.getDbWithCtx(ctx).Create(&log)
	return log, res.Error
}

// ListLogs 查询审批日志
func (d *Dao) ListLogs(ctx context.Context, q orm.Query) (logs []*model.Log, err error) {
	res := q.Where(d.getDbWithCtx(ctx)).Find(&logs)
	return logs, res.Error
}

package main

import (
	"context"
	"eago/common/log"
	"eago/task/conf"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/srv/builtin"
	task "eago/task/srv/proto"
)

// GetResult 查看任务结果
func (ts *TaskService) GetResult(ctx context.Context, in *task.TaskUniqueId, out *task.Result) error {
	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := builtin.TaskUniqueIdDecode(in.TaskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": in.TaskUniqueId,
			"error":          err,
		}, "An error occurred while builtin.TaskUniqueIdDecode.")
		return err
	}

	// 取数据库中任务结果记录
	obj, ok := dao.GetResult(p, id)
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetResult.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.RpcError()
	}

	out.TaskCodename = obj.TaskCodename
	out.Status = int32(obj.Status)
	out.Caller = obj.Worker
	out.Worker = obj.Worker
	if obj.StartAt != nil {
		out.StartAt = obj.StartAt.Format(conf.TIMESTAMP_FORMAT)
	}
	if obj.EndAt != nil {
		out.EndAt = obj.EndAt.Format(conf.TIMESTAMP_FORMAT)
	}
	return nil
}

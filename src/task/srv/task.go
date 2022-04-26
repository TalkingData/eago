package main

import (
	"context"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/model"
	"eago/task/srv/builtin"
	task "eago/task/srv/proto"
	"eago/task/worker"
	"io"
)

// PagedListTasks 列出所有任务-分页
func (taskSrv *TaskService) PagedListTasks(ctx context.Context, in *task.QueryWithPage, out *task.PagedTasks) error {
	log.Info("taskSrv.PagedListTasks called.")
	defer log.Info("taskSrv.PagedListTasks end.")

	query := make(dao.Query)
	for k, v := range in.Query {
		query[k] = v
	}
	pagedData, ok := dao.PagedListTasks(query, int(in.Page), int(in.PageSize))
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.PagedListTasks.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	out.Tasks = make([]*task.Task, 0)
	for _, t := range *pagedData.Data.(*[]model.Task) {
		nTask := &task.Task{
			Id:           int32(t.Id),
			Codename:     t.Codename,
			FormalParams: t.FormalParams,
			Description:  *t.Description,
		}
		out.Tasks = append(out.Tasks, nTask)
	}

	out.Page = uint32(pagedData.Page)
	out.Pages = uint32(pagedData.Pages)
	out.PageSize = uint32(pagedData.PageSize)
	out.Total = uint32(pagedData.Total)
	return nil
}

// CallTask 调用任务
func (taskSrv *TaskService) CallTask(ctx context.Context, in *task.CallTaskReq, out *task.TaskUniqueId) error {
	log.Info("taskSrv.CallTask called.")
	defer log.Info("taskSrv.CallTask end.")

	tId, err := builtin.CallTask(in.TaskCodename, string(in.Arguments), in.Caller, in.Timeout)
	if err != nil {
		m := msg.UndefinedError.SetError(err, "An error occurred while builtin.CallTask.")
		log.ErrorWithFields(log.Fields{
			"task_codename": in.TaskCodename,
			"arguments":     in.Arguments,
			"timeout":       in.Timeout,
			"caller":        in.Caller,
			"error":         err,
		}, m.String())
		return m.RpcError()
	}

	log.DebugWithFields(log.Fields{
		"task_unique_id": tId,
	}, "CallTask success.")

	out.TaskUniqueId = tId
	return nil
}

// KillTask 结束任务
func (taskSrv *TaskService) KillTask(ctx context.Context, in *task.TaskUniqueId, out *task.BoolMsg) error {
	log.Info("taskSrv.CallTask KillTask.")
	defer log.Info("taskSrv.KillTask end.")

	err := builtin.KillTask(in.TaskUniqueId)
	if err != nil {
		m := msg.UndefinedError.SetError(err, "An error occurred while builtin.KillTask.")
		log.ErrorWithFields(log.Fields{
			"task_unique_id": in.TaskUniqueId,
			"error":          err,
		}, m.String())
		return m.RpcError()
	}

	log.DebugWithFields(log.Fields{
		"task_unique_id": in.TaskUniqueId,
	}, "KillTask success.")

	out.Ok = true
	return nil
}

// SetTaskStatus 设置任务状态
func (taskSrv *TaskService) SetTaskStatus(ctx context.Context, in *task.SetTaskStatusReq, out *task.BoolMsg) error {
	log.Info("taskSrv.SetTaskStatus KillTask.")
	defer log.Info("taskSrv.SetTaskStatus end.")

	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := builtin.TaskUniqueIdDecode(in.TaskUniqueId)
	if err != nil {
		m := msg.UndefinedError.SetError(err, "An error occurred while builtin.TaskUniqueIdDecode.")
		log.ErrorWithFields(log.Fields{
			"task_unique_id": in.TaskUniqueId,
			"error":          err,
		}, m.String())
		return m.RpcError()
	}

	// 取数据库中任务结果记录
	obj, ok := dao.GetResult(p, id)
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.GetResult.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.RpcError()
	}
	// 找不到数据的处理
	if obj == nil {
		m := msg.UndefinedError.SetDetail("Result object not found.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.RpcError()
	}

	// 判断任务记录，无法结束不是在执行的状态就，返回错误
	if obj.Status <= worker.TASK_SUCCESS_END_STATUS {
		m := msg.UndefinedError.SetDetail("Result object wrong state, task was ended.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.RpcError()
	}

	// Status小于等于worker.TASK_SUCCESS_END_STATUS，则说明任务已结束
	end := false
	if int(in.Status) <= worker.TASK_SUCCESS_END_STATUS {
		end = true
	}
	ok = dao.SetResultStatus(p, id, int(in.Status), end)
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.SetResultStatus.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.RpcError()
	}

	log.DebugWithFields(log.Fields{
		"partition": p,
		"result_id": id,
	}, "SetTaskStatus success.")

	out.Ok = true
	return nil
}

// AppendTaskLog 追加任务日志
func (taskSrv *TaskService) AppendTaskLog(ctx context.Context, stream task.TaskService_AppendTaskLogStream) error {
	log.Info("taskSrv.AppendTaskLog KillTask.")
	defer log.Info("taskSrv.AppendTaskLog end.")

	for {
		// 接受请求流数据
		tlq, err := stream.Recv()
		// 流结束退出
		if err == io.EOF {
			break
		}
		if err != nil {
			m := msg.UndefinedError.SetError(err, "An error occurred while stream.Recv.")
			log.ErrorWithFields(m.LogFields())
			return m.RpcError()
		}
		// 新建Log
		if err = builtin.NewLog(tlq.TaskUniqueId, &tlq.Content); err != nil {
			m := msg.UndefinedError.SetError(err, "An error occurred while local.NewLog.")
			log.ErrorWithFields(log.Fields{
				"task_unique_id": tlq.TaskUniqueId,
				"error":          err,
			}, m.String())
			return m.RpcError()
		}

		// 返回请求结果给客户端
		if err = stream.Send(&task.BoolMsg{Ok: true}); err != nil {
			m := msg.UndefinedError.SetError(err, "An error occurred while local.Send.")
			log.ErrorWithFields(log.Fields{
				"task_unique_id": tlq.TaskUniqueId,
				"error":          err,
			}, m.String())
			return m.RpcError()
		}

	}
	if err := stream.Close(); err != nil {
		m := msg.UndefinedError.SetError(err, "An error occurred while stream.Close.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	return nil
}

package main

import (
	"context"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/model"
	"eago/task/srv/local"
	task "eago/task/srv/proto"
	"eago/task/worker"
	"fmt"
	"io"
)

// ListTasks 列出所有任务
func (ts *TaskService) ListTasks(ctx context.Context, req *task.Empty, rsp *task.Tasks) error {
	log.Info("Got rpc call ListTasks.")
	defer log.Info("Rpc call ListTasks done.")

	tasks := make([]*task.Task, 0)

	// 查询所有启用的任务
	objs, ok := model.ListTasks(model.Query{"Disabled": 0})
	if !ok {
		m := msg.ErrDatabase.SetDetail("GetUser object failed.")
		log.Error(m.String())
		return m.Error()
	}

	for _, obj := range *objs {
		t := task.Task{
			Id:          int32(obj.Id),
			Codename:    obj.Codename,
			Arguments:   obj.Arguments,
			Description: *obj.Description,
		}
		tasks = append(tasks, &t)
	}

	rsp.Tasks = tasks
	return nil
}

// CallTask 调用任务
func (ts *TaskService) CallTask(ctx context.Context, req *task.CallTaskReq, rsp *task.CallTaskRsp) error {
	log.Info("Got rpc call TaskCall.")
	defer log.Info("Rpc call TaskCall done.")

	tId, err := local.CallTask(req.TaskCodename, req.Arguments, req.Caller, int64(req.Timeout))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_codename": req.TaskCodename,
			"arguments":     req.Arguments,
			"timeout":       req.Timeout,
			"caller":        req.Caller,
		}, "Error when TaskService.AppendTaskLog called, in local.CallTask")
		return err
	}

	log.DebugWithFields(log.Fields{
		"task_unique_id": tId,
	}, "TaskCall success.")

	rsp.TaskUniqueId = tId
	return nil
}

// SetTaskStatus 设置任务状态
func (ts *TaskService) SetTaskStatus(ctx context.Context, req *task.SetTaskStatusReq, rsp *task.BoolMsg) error {
	log.Info("Got rpc call TaskDone.")
	defer log.Info("Rpc call TaskDone done.")

	// 将任务唯一Id解码为任务结果Id和分区
	p, id, err := local.TaskUniqueIdDecode(req.TaskUniqueId)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": req.TaskUniqueId,
			"error":          err.Error(),
		}, "Error when TaskService.AppendTaskLog called, in local.TaskUniqueIdDecode.")
		return err
	}

	// 取数据库中任务结果记录
	obj, ok := model.GetResult(p, id)
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when TaskService.AppendTaskLog called, in model.GetResult.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.Error()
	}
	// 找不到数据的处理
	if obj == nil {
		m := msg.WarnNotFound.SetDetail("Result object not found.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.Error()
	}

	// 判断任务记录，无法结束不是在执行的状态就，返回错误
	if obj.Status <= worker.TASK_SUCCESS_END_STATUS {
		err := fmt.Errorf("Wrong state, The task that was ended, Cannot set it to DONE status.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
			"error":     err.Error(),
		}, "Error when TaskService.AppendTaskLog called, in check result status.")
		return err
	}

	// Status小于等于worker.TASK_SUCCESS_END_STATUS，则说明任务已结束
	end := false
	if int(req.Status) <= worker.TASK_SUCCESS_END_STATUS {
		end = true
	}
	ok = model.SetResultStatus(p, id, int(req.Status), end)
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when TaskService.AppendTaskLog called, in model.SetResultStatus.")
		log.ErrorWithFields(log.Fields{
			"partition": p,
			"result_id": id,
		}, m.String())
		return m.Error()
	}

	log.DebugWithFields(log.Fields{
		"partition": p,
		"result_id": id,
	}, "SetTaskStatus success.")

	rsp.Ok = true
	return nil
}

// AppendTaskLog 追加任务日志
func (ts *TaskService) AppendTaskLog(ctx context.Context, stream task.TaskService_AppendTaskLogStream) error {
	log.Info("Got rpc call TaskLog.")
	defer log.Info("Rpc call TaskLog done.")

	for {
		// 接受请求流数据
		tlq, err := stream.Recv()
		// 流结束退出
		if err == io.EOF {
			break
		}
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err.Error(),
			}, "Error when TaskService.AppendTaskLog called, in stream.Recv.")
			return err
		}
		// 新增Log
		if err := local.NewLog(tlq.TaskUniqueId, &tlq.Content); err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": tlq.TaskUniqueId,
				"error":          err.Error(),
			}, "Error when TaskService.AppendTaskLog called, in local.NewLog.")
			return err
		}

		// 返回请求结果给客户端
		if err := stream.Send(&task.BoolMsg{Ok: true}); err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": tlq.TaskUniqueId,
				"error":          err.Error(),
			}, "Error when TaskService.AppendTaskLog called, in stream.Send.")
			return err
		}

	}
	if err := stream.Close(); err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error when TaskService.AppendTaskLog called, in stream.Close.")
		return err
	}

	return nil
}

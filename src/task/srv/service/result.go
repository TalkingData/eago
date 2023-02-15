package service

import (
	"context"
	"eago/common/logger"
	"eago/task/conf/msg"
	"eago/task/dto"
	taskpb "eago/task/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (taskSrv *TaskService) SetResultStatus(
	ctx context.Context, req *taskpb.SetResultStatusReq, _ *emptypb.Empty,
) error {
	taskSrv.logger.Info("taskSrv.SetResultStatus called.")
	defer taskSrv.logger.Info("taskSrv.SetResultStatus end.")

	// 将任务唯一Id解码为任务结果Id和分区
	part, resId, err := taskSrv.biz.TaskUniqueIdDecode(req.TaskUniqueId)
	if err != nil {
		m := msg.MsgTaskUniqueIdDecodeFailed.SetError(err)
		taskSrv.logger.ErrorWithFields(
			m.ToLoggerFields().Append("task_unique_id", req.TaskUniqueId),
			"An error occurred while biz.TaskUniqueIdDecode in taskSrv.SetResultStatus.",
		)
		return m.ToMicroErr()
	}

	// 取数据库中任务结果记录
	obj, err := taskSrv.dao.GetResult(ctx, part, resId)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		taskSrv.logger.ErrorWithFields(
			m.ToLoggerFields().Append("query", req),
			"An error occurred while dao.GetResult in taskSrv.SetResultStatus.",
		)
		return m.ToMicroErr()
	}
	// 找不到数据的处理
	if obj == nil || obj.Id < 1 {
		m := msg.MsgTaskDaoErr.SetDetail("Result object not found.")
		f := m.ToLoggerFields()
		f["partition"] = part
		f["result_id"] = resId
		taskSrv.logger.ErrorWithFields(f, "An error occurred while dao.GetResult in taskSrv.SetResultStatus.")
		return m.ToMicroErr()
	}

	// 判断任务记录，无法结束不是在执行的状态，返回错误
	if obj.Status <= dto.TaskResultStatusSuccessEnd {
		m := msg.MsgSetResultStatusInvalidStatusFailed
		f := m.ToLoggerFields()
		f["partition"] = part
		f["result_id"] = resId
		taskSrv.logger.ErrorWithFields(f, "An error occurred while dao.GetResult in taskSrv.SetResultStatus.")
	}

	// Status小于等于worker.TaskStatusSuccessEnd，则说明任务已结束
	end := false
	if int(req.Status) <= dto.TaskResultStatusSuccessEnd {
		end = true
	}
	if err = taskSrv.dao.SetResultStatus(ctx, part, resId, req.Status, end); err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		f := m.ToLoggerFields()
		f["partition"] = part
		f["result_id"] = resId
		taskSrv.logger.ErrorWithFields(
			f, "An error occurred while dao.SetResultStatus in taskSrv.SetResultStatus.",
		)
		return m.ToMicroErr()
	}

	taskSrv.logger.DebugWithFields(logger.Fields{
		"partition": part,
		"result_id": resId,
	}, "SetTaskStatus success.")

	return nil
}

func (taskSrv *TaskService) GetResult(ctx context.Context, req *taskpb.TaskUniqueId, rsp *taskpb.Result) error {
	// 将任务唯一Id解码为任务结果Id和分区
	part, resId, err := taskSrv.biz.TaskUniqueIdDecode(req.TaskUniqueId)
	if err != nil {
		taskSrv.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": req.TaskUniqueId,
			"error":          err,
		}, "An error occurred while builtin.TaskUniqueIdDecode in taskSrv.GetResult.")
		return err
	}

	// 取数据库中任务结果记录
	obj, err := taskSrv.dao.GetResult(ctx, part, resId)
	if err != nil {
		cMsg := msg.MsgTaskDaoErr.SetError(err)
		f := cMsg.ToLoggerFields()
		f["partition"] = part
		f["result_id"] = resId
		taskSrv.logger.ErrorWithFields(f, "An error occurred while dao.GetResult in taskSrv.GetResult.")
		return cMsg.ToMicroErr()
	}

	rsp.TaskCodename = obj.TaskCodename
	rsp.Status = obj.Status
	rsp.Caller = obj.Caller
	rsp.Worker = obj.Worker
	if obj.StartAt != nil {
		rsp.StartAt = obj.StartAt.String()
	}
	if obj.EndAt != nil {
		rsp.EndAt = obj.EndAt.String()
	}

	return nil
}

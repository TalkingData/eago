package service

import (
	"context"
	"eago/common/logger"
	"eago/common/orm"
	commonpb "eago/common/proto"
	"eago/task/conf/msg"
	"eago/task/model"
	taskpb "eago/task/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (taskSrv *TaskService) PagedListTasks(
	ctx context.Context, req *commonpb.QueryWithPage, rsp *taskpb.PagedTasks,
) error {
	taskSrv.logger.Info("taskSrv.PagedListTasks called.")
	defer taskSrv.logger.Info("taskSrv.PagedListTasks end.")

	pagedData, err := taskSrv.dao.PagedListTasks(
		ctx, orm.NewQueryByMapStrStr(req.Query), int(req.Page), int(req.PageSize),
	)
	if err != nil {
		m := msg.MsgTaskDaoErr.SetError(err)
		taskSrv.logger.ErrorWithFields(
			m.ToLoggerFields().Append("query", req),
			"An error occurred while dao.PagedListTasks in taskSrv.PagedListTasks.",
		)
		return m.ToMicroErr()
	}

	rsp.Tasks = make([]*taskpb.Task, 0)
	for _, t := range *pagedData.Data.(*[]*model.Task) {
		rsp.Tasks = append(rsp.Tasks, &taskpb.Task{
			Id:           t.Id,
			Codename:     t.Codename,
			FormalParams: t.FormalParams,
			Description:  *t.Description,
		})
	}

	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

func (taskSrv *TaskService) CallTask(ctx context.Context, req *taskpb.CallTaskReq, rsp *taskpb.TaskUniqueId) error {
	taskSrv.logger.InfoWithFields(logger.Fields{
		"task_codename": req.TaskCodename,
		"caller":        req.Caller,
	}, "taskSrv.CallTask called.")
	defer taskSrv.logger.Info("taskSrv.CallTask end.")

	tId, err := taskSrv.biz.CallTask(ctx, req.TaskCodename, string(req.Arguments), req.Caller, req.Timeout)
	if err != nil {
		m := msg.MsgCallTaskFailed.SetError(err)
		f := m.ToLoggerFields()
		f["task_codename"] = req.TaskCodename
		f["arguments"] = req.Arguments
		f["timeout"] = req.Timeout
		f["caller"] = req.Caller
		taskSrv.logger.ErrorWithFields(f, "An error occurred while biz.CallTask in taskSrv.CallTask.")
		return m.ToMicroErr()
	}

	taskSrv.logger.DebugWithFields(logger.Fields{
		"task_unique_id": tId,
	}, "CallTask success.")

	rsp.TaskUniqueId = tId
	return nil
}

func (taskSrv *TaskService) KillTask(ctx context.Context, req *taskpb.TaskUniqueId, _ *emptypb.Empty) error {
	taskSrv.logger.InfoWithFields(logger.Fields{
		"task_unique_id": req.TaskUniqueId,
	}, "taskSrv.KillTask called.")
	defer taskSrv.logger.Info("taskSrv.KillTask end.")

	if err := taskSrv.biz.KillTask(ctx, req.TaskUniqueId); err != nil {
		m := msg.MsgKillTaskFailed.SetError(err)
		taskSrv.logger.ErrorWithFields(
			m.ToLoggerFields().Append("task_unique_id", req.TaskUniqueId),
			"An error occurred while biz.CallTask in taskSrv.KillTask.",
		)
		return m.ToMicroErr()
	}

	taskSrv.logger.DebugWithFields(logger.Fields{
		"task_unique_id": req.TaskUniqueId,
	}, "KillTask success.")

	return nil
}

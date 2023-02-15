package service

import (
	"context"
	"eago/common/orm"
	commonpb "eago/common/proto"
	"eago/task/conf/msg"
	"eago/task/model"
	taskpb "eago/task/proto"
)

// PagedListSchedules 列出所有计划任务-分页
func (taskSrv *TaskService) PagedListSchedules(
	ctx context.Context, req *commonpb.QueryWithPage, out *taskpb.PagedSchedules,
) error {
	taskSrv.logger.Info("taskSrv.PagedListSchedules called.")
	defer taskSrv.logger.Info("taskSrv.PagedListSchedules end.")

	pagedData, err := taskSrv.dao.PagedListSchedules(
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

	out.Schedules = make([]*taskpb.Schedule, 0)
	for _, s := range *pagedData.Data.(*[]*model.Schedule) {
		out.Schedules = append(out.Schedules, &taskpb.Schedule{
			Id:           s.Id,
			TaskCodename: s.TaskCodename,
			Expression:   s.Expression,
			Timeout:      *s.Timeout,
			Arguments:    s.Arguments,
			Disabled:     *s.Disabled,
		})
	}

	out.Page = uint32(pagedData.Page)
	out.Pages = uint32(pagedData.Pages)
	out.PageSize = uint32(pagedData.PageSize)
	out.Total = uint32(pagedData.Total)
	return nil
}

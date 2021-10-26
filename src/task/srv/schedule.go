package main

import (
	"context"
	"eago/common/log"
	"eago/task/conf/msg"
	"eago/task/dao"
	"eago/task/model"
	task "eago/task/srv/proto"
)

// ListSchedule 列出所有计划任务
func (ts *TaskService) ListSchedules(ctx context.Context, in *task.QueryWithPage, out *task.PagedSchedules) error {
	log.Info("Got rpc call ListSchedules.")
	defer log.Info("Rpc call ListSchedules done.")

	query := make(dao.Query)
	for k, v := range in.Query {
		query[k] = v
	}
	pagedData, ok := dao.PagedListSchedules(query, int(in.Page), int(in.PageSize))
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.PagedListSchedules.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	out.Schedules = make([]*task.Schedule, 0)
	for _, s := range *pagedData.Data.(*[]model.Schedule) {
		out.Schedules = append(out.Schedules, &task.Schedule{
			Id:           int32(s.Id),
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

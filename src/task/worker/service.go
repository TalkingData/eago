package worker

import (
	"context"
	"eago/common/log"
	worker_proto "eago/task/worker/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TaskWorkerService struct
type TaskWorkerService struct {
	wk *worker
}

// NewTaskWorkerService 创建Worker服务
func NewTaskWorkerService(wk *worker) *TaskWorkerService {
	return &TaskWorkerService{
		wk: wk,
	}
}

// CallTask 调用任务
func (tws *TaskWorkerService) CallTask(_ context.Context, req *worker_proto.CallTaskReq) (*emptypb.Empty, error) {
	log.InfoWithFields(log.Fields{
		"worker_id":      tws.wk.workerId,
		"task_codename":  req.TaskCodename,
		"task_unique_id": req.TaskUniqueId,
		"caller":         req.Caller,
	}, "Got a call task request.")
	tws.wk.callTask(
		req.TaskCodename,
		req.TaskUniqueId,
		string(req.Arguments),
		req.Timeout,
		req.Caller,
		req.Timestamp,
	)

	return new(emptypb.Empty), nil
}

// KillTask 结束任务
func (tws *TaskWorkerService) KillTask(ctx context.Context, req *worker_proto.KillTaskReq) (*emptypb.Empty, error) {
	log.InfoWithFields(log.Fields{
		"worker_id":      tws.wk.workerId,
		"task_unique_id": req.TaskUniqueId,
	}, "Got a kill task request.")
	tws.wk.killTask(req.TaskUniqueId)

	return new(emptypb.Empty), nil
}

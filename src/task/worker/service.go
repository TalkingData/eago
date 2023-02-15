package worker

import (
	"context"
	"eago/common/logger"
	workerpb "eago/task/worker/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// taskWorkerService struct
type taskWorkerService struct {
	wk *worker

	logger *logger.Logger
}

// NewTaskWorkerService 创建Worker服务
func NewTaskWorkerService(wk *worker, logger *logger.Logger) *taskWorkerService {
	return &taskWorkerService{
		wk: wk,

		logger: logger,
	}
}

// CallTask 调用任务
func (tws *taskWorkerService) CallTask(_ context.Context, req *workerpb.CallTaskReq) (*emptypb.Empty, error) {
	tws.logger.InfoWithFields(logger.Fields{
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

	return &emptypb.Empty{}, nil
}

// KillTask 结束任务
func (tws *taskWorkerService) KillTask(_ context.Context, req *workerpb.KillTaskReq) (*emptypb.Empty, error) {
	tws.logger.InfoWithFields(logger.Fields{
		"worker_id":      tws.wk.workerId,
		"task_unique_id": req.TaskUniqueId,
	}, "Got a kill task request.")
	tws.wk.killTask(req.TaskUniqueId)

	return &emptypb.Empty{}, nil
}

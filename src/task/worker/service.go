package worker

import (
	"eago/common/log"
)

type WorkerService struct {
	wk *worker
}

// KillTask
func (s *WorkerService) KillTask(req KillTaskReq, rsp *WorkerResponse) error {
	defer func() {
		rsp.Ok = true
		rsp.Message = "Success"
	}()

	log.InfoWithFields(log.Fields{
		"worker_id":      s.wk.workerId,
		"task_unique_id": req.TaskUniqueId,
	}, "Got a kill task request.")
	s.wk.killTask(req.TaskUniqueId)

	return nil
}

// CallTask
func (s *WorkerService) CallTask(req CallTaskReq, rsp *WorkerResponse) error {
	defer func() {
		rsp.Ok = true
		rsp.Message = "Success"
	}()

	log.InfoWithFields(log.Fields{
		"worker_id":      s.wk.workerId,
		"task_codename":  req.TaskCodename,
		"task_unique_id": req.TaskUniqueId,
		"caller":         req.Caller,
	}, "Got a call task request.")
	s.wk.callTask(&req)

	return nil
}

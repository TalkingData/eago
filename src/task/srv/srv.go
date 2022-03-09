package main

import (
	"context"
	"eago/common/log"
	"eago/task/srv/builtin"
	task "eago/task/srv/proto"
)

// IsAllowedSrv 判断Srv是否已经在白名单内
func (ts *TaskService) IsAllowedSrv(ctx context.Context, req *task.IsAllowedSrvQuery, rsp *task.BoolMsg) error {
	log.Info("srv.IsAllowedSrv called.")
	defer log.Info("srv.IsAllowedSrv end.")

	rsp.Ok = builtin.IsSrvAllowed(req.SrvAddr, req.WorkerAddr)
	return nil
}

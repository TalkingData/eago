package service

import (
	"context"
	taskpb "eago/task/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (taskSrv *TaskService) IsValidSrvToken(
	ctx context.Context, req *taskpb.SrvTokenQuery, rsp *wrapperspb.BoolValue,
) error {
	taskSrv.logger.Debug("taskSrv.IsValidSrvToken called.")
	defer taskSrv.logger.Info("taskSrv.IsValidSrvToken end.")

	rsp.Value = taskSrv.biz.VerifyAndUnregisterSrvToken(ctx, req.Value)
	return nil
}

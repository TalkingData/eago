package service

import (
	"context"
	cMsg "eago/common/code_msg"
	"eago/task/conf/msg"
	taskpb "eago/task/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

func (taskSrv *TaskService) AppendTaskLog(ctx context.Context, stream taskpb.TaskService_AppendTaskLogStream) error {
	taskSrv.logger.Debug("taskSrv.AppendTaskLog called.")
	defer taskSrv.logger.Debug("taskSrv.AppendTaskLog end.")

	for {
		// 接受请求流数据
		tlq, err := stream.Recv()
		// 流结束退出
		if err == io.EOF {
			break
		}
		if err != nil {
			m := cMsg.MsgUndefinedErr.SetError(err)
			taskSrv.logger.ErrorWithFields(
				m.ToLoggerFields(),
				"An error occurred while taskSrv.AppendTaskLog.",
			)
			return m.ToMicroErr()
		}
		// 新建Log
		if err = taskSrv.biz.NewLog(ctx, tlq.TaskUniqueId, &tlq.Content); err != nil {
			m := msg.MsgBizNewLogFailed.SetError(err)
			taskSrv.logger.ErrorWithFields(
				m.ToLoggerFields().Append("task_unique_id", tlq.TaskUniqueId),
				"An error occurred while biz.NewLog in taskSrv.AppendTaskLog.",
			)
			return m.ToMicroErr()
		}

		// 返回请求结果给客户端
		if err = stream.Send(&emptypb.Empty{}); err != nil {
			m := msg.MsgNewLogStreamSendFailed.SetError(err)
			taskSrv.logger.ErrorWithFields(
				m.ToLoggerFields().Append("task_unique_id", tlq.TaskUniqueId),
				"An error occurred while stream.Send in taskSrv.AppendTaskLog.",
			)
			return m.ToMicroErr()
		}
	}

	if err := stream.Close(); err != nil {
		m := msg.MsgNewLogStreamCloseFailed.SetError(err)
		taskSrv.logger.WarnWithFields(
			m.ToLoggerFields(),
			"An error occurred while stream.Close in taskSrv.AppendTaskLog, skipped it.",
		)
		return m.ToMicroErr()
	}

	return nil

}

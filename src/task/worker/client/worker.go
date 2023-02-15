package client

import (
	"context"
	"eago/common/logger"
	"eago/task/dto"
	"eago/task/worker"
	workerpb "eago/task/worker/proto"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"time"
)

type WorkerClient struct {
	etcdCli *clientv3.Client

	logger *logger.Logger
}

func NewWorkerCli(etcdAddrs []string, _logger *logger.Logger) *WorkerClient {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		_logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while clientv3.New in NewWorkerCli.")
		panic(err)
	}

	return &WorkerClient{
		etcdCli: etcdCli,

		logger: _logger,
	}
}

// List 列出所有活跃的Worker
func (wkCli *WorkerClient) List(ctx context.Context) (workers []*dto.WorkerInfo) {
	resp, err := wkCli.etcdCli.Get(ctx, worker.WorkerRegisterKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		wkCli.logger.WarnWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while etcdCli.GetUser in WorkerClient.List, skipped it.")
		return
	}

	for _, ev := range resp.Kvs {
		wk := &dto.WorkerInfo{}

		if err = json.Unmarshal(ev.Value, wk); err != nil {
			wkCli.logger.WarnWithFields(logger.Fields{
				"error": err,
			}, "An error occurred while json.Unmarshal in WorkerClient.List, skipped it.")
			continue
		}
		workers = append(workers, wk)
	}

	return
}

// ListByModular 按模块名列出所有活跃的Worker
func (wkCli *WorkerClient) ListByModular(ctx context.Context, m string) []*dto.WorkerInfo {
	workers := make([]*dto.WorkerInfo, 0)

	for _, wk := range wkCli.List(ctx) {
		if wk.Modular != m {
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// GetWorkerById 获得指定Worker
func (wkCli *WorkerClient) GetWorkerById(ctx context.Context, wkId string) *dto.WorkerInfo {
	for _, wk := range wkCli.List(ctx) {
		if wk.WorkerId == wkId {
			return wk
		}
	}

	return nil
}

// CallTask 调用任务
func (wkCli *WorkerClient) CallTask(
	ctx context.Context,
	wk *dto.WorkerInfo,
	codename, uniqueId, arguments, caller string,
	timeout int64,
	startTimestamp int64,
) error {
	// 连接Worker并获取worker grpc客户端
	cli, err := wkCli.getWorkerCli(wk)
	if err != nil {
		return err
	}

	req := &workerpb.CallTaskReq{
		TaskCodename: codename,
		TaskUniqueId: uniqueId,
		Arguments:    []byte(arguments),
		Timeout:      timeout,
		Caller:       caller,
		Timestamp:    startTimestamp,
	}
	// 调用 TaskWorkerService.CallTask
	if _, err = cli.CallTask(ctx, req); err != nil {
		return err
	}

	return nil
}

// KillTask 调用任务
func (wkCli *WorkerClient) KillTask(ctx context.Context, wk *dto.WorkerInfo, taskUniqueId string) error {
	// 连接Worker并获取worker grpc客户端
	cli, err := wkCli.getWorkerCli(wk)
	if err != nil {
		return err
	}

	req := &workerpb.KillTaskReq{
		TaskUniqueId: taskUniqueId,
		Timestamp:    time.Now().Unix(),
	}
	// 调用 TaskWorkerService.KillTask
	if _, err = cli.KillTask(ctx, req); err != nil {
		return err
	}

	return nil
}

// getWorkerCli 连接Worker并获取worker grpc客户端
func (wkCli *WorkerClient) getWorkerCli(wk *dto.WorkerInfo) (workerpb.TaskWorkerServiceClient, error) {
	conn, err := grpc.Dial(wk.Address, grpc.WithInsecure())
	if err != nil {
		wkCli.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while grpc.Dial in WorkerClient.getWorkerCli.")
		return nil, err
	}

	return workerpb.NewTaskWorkerServiceClient(conn), nil
}

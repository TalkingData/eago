package cli

import (
	"context"
	"eago/common/log"
	"eago/common/tracer"
	"eago/task/conf"
	"eago/task/worker"
	workerProto "eago/task/worker/proto"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"net"
	"time"
)

type AfterConnWrapHandler func(conn net.Conn)

var WorkerClient *workerClient

type workerClient struct {
	etcdCli *clientv3.Client
}

// InitWorkerCli 启动Worker客户端
func InitWorkerCli() {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Conf.EtcdAddresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while clientv3.New.")
		panic(err)
	}

	WorkerClient = &workerClient{
		etcdCli: etcdCli,
	}
}

// List 列出所有活跃的Worker
func (wkCli *workerClient) List(ctx context.Context) []*worker.WorkerInfo {
	sp, c := tracer.StartSpanFromContext(ctx)
	defer sp.Finish()
	workers := make([]*worker.WorkerInfo, 0)

	resp, err := wkCli.etcdCli.Get(c, worker.WORKER_REGISTER_KEY_PREFFIX, clientv3.WithPrefix())
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while etcdCli.GetUser.")
		return workers
	}

	for _, ev := range resp.Kvs {
		wk := &worker.WorkerInfo{}

		if err := json.Unmarshal(ev.Value, wk); err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err,
			}, "An error occurred while json.Unmarshal for worker info.")
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// ListByModular 按模块名列出所有活跃的Worker
func (wkCli *workerClient) ListByModular(m string) []*worker.WorkerInfo {
	workers := make([]*worker.WorkerInfo, 0)

	for _, wk := range wkCli.List(context.Background()) {
		if wk.Modular != m {
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// GetWorkerById 获得指定Worker
func (wkCli *workerClient) GetWorkerById(wkId string) *worker.WorkerInfo {
	for _, wk := range wkCli.List(context.Background()) {
		if wk.WorkerId == wkId {
			return wk
		}
	}

	return nil
}

// CallTask 调用任务
func (wkCli *workerClient) CallTask(wk *worker.WorkerInfo, codename, uniqueId, arguments, caller string, timeout int32, startTimestamp int64) error {
	// 连接Worker并获取worker grpc客户端
	cli, err := wkCli.getWorkerCli(wk)
	if err != nil {
		return err
	}

	req := &workerProto.CallTaskReq{
		TaskCodename: codename,
		TaskUniqueId: uniqueId,
		Arguments:    []byte(arguments),
		Timeout:      timeout,
		Caller:       caller,
		Timestamp:    startTimestamp,
	}
	// 调用 TaskWorkerService.CallTask
	if _, err = cli.CallTask(context.TODO(), req); err != nil {
		return err
	}

	return nil
}

// KillTask 调用任务
func (wkCli *workerClient) KillTask(wk *worker.WorkerInfo, taskUniqueId string) error {
	// 连接Worker并获取worker grpc客户端
	cli, err := wkCli.getWorkerCli(wk)
	if err != nil {
		return err
	}

	req := &workerProto.KillTaskReq{
		TaskUniqueId: taskUniqueId,
		Timestamp:    time.Now().Unix(),
	}
	// 调用 TaskWorkerService.KillTask
	if _, err = cli.KillTask(context.TODO(), req); err != nil {
		return err
	}

	return nil
}

// getWorkerCli 连接Worker并获取worker grpc客户端
func (wkCli *workerClient) getWorkerCli(wk *worker.WorkerInfo) (workerProto.TaskWorkerServiceClient, error) {
	conn, err := grpc.Dial(wk.Address, grpc.WithInsecure())
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while grpc.Dial worker address.")
		return nil, err
	}

	return workerProto.NewTaskWorkerServiceClient(conn), nil
}

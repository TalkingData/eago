package cli

import (
	"context"
	"eago/common/log"
	"eago/task/conf"
	"eago/task/worker"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var WorkerClient *workerClient

type workerClient struct {
	etcdCli *clientv3.Client
}

// InitWorkerCli 启动Worker客户端
func InitWorkerCli() {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Config.EtcdAddresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in clientv3.New.")
		panic(err)
	}

	WorkerClient = &workerClient{etcdCli}

}

// List 列出所有活跃的Worker
func (wc *workerClient) List() []*worker.WorkerInfo {
	workers := make([]*worker.WorkerInfo, 0)

	resp, err := wc.etcdCli.Get(context.Background(), worker.WORKER_REGISTER_KEY_PREFFIX, clientv3.WithPrefix())
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in etcdCli.GetUser.")
		return nil
	}

	for _, ev := range resp.Kvs {
		wk := &worker.WorkerInfo{}

		if err := json.Unmarshal(ev.Value, wk); err != nil {
			log.ErrorWithFields(log.Fields{
				"error": err.Error(),
			}, "Error in json.Unmarshal.")
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// ListByModular 按模块名列出所有活跃的Worker
func (wc *workerClient) ListByModular(m string) []*worker.WorkerInfo {
	workers := make([]*worker.WorkerInfo, 0)

	for _, wk := range wc.List() {
		if wk.Modular != m {
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// GetWorkerById 获得指定Worker
func (wc *workerClient) GetWorkerById(wkId string) *worker.WorkerInfo {
	for _, wk := range wc.List() {
		if wk.WorkerId == wkId {
			return wk
		}
	}

	return nil
}

// CallTask 调用任务
func (wc *workerClient) CallTask(wk *worker.WorkerInfo, codename, uniqueId, arguments, caller string, timeout, startTimestamp int64) error {
	req := worker.CallTaskReq{
		TaskCodename: codename,
		TaskUniqueId: uniqueId,
		Arguments:    arguments,
		Timeout:      timeout,
		Caller:       caller,
		Timestamp:    startTimestamp,
	}

	// 连接Worker
	conn, err := wc.connWorker(wk)
	if err != nil {
		return err
	}

	// 调用 WorkerService.CallTask
	wkRsp := worker.WorkerResponse{}
	err = conn.Call("WorkerService.CallTask", req, &wkRsp)
	if err != nil {
		return err
	}

	if !wkRsp.Ok {
		log.ErrorWithFields(log.Fields{
			"ok":      wkRsp.Ok,
			"message": wkRsp.Message,
		}, "Error in workerClient.CallTask, worker returned not ok status.")
		return err
	}

	return nil
}

// KillTask 调用任务
func (wc *workerClient) KillTask(wk *worker.WorkerInfo, taskUniqueId string) error {
	req := worker.KillTaskReq{
		TaskUniqueId: taskUniqueId,
		Timestamp:    time.Now().Unix(),
	}

	// 连接Worker
	conn, err := wc.connWorker(wk)
	if err != nil {
		return err
	}

	// 调用 WorkerService.KillTask
	wkRsp := worker.WorkerResponse{}
	err = conn.Call("WorkerService.KillTask", req, &wkRsp)
	if err != nil {
		return err
	}

	if !wkRsp.Ok {
		log.ErrorWithFields(log.Fields{
			"ok":      wkRsp.Ok,
			"message": wkRsp.Message,
		}, "Error in workerClient.CallTask, worker returned not ok status.")
		return err
	}

	return nil
}

// connWorker 连接Worker
func (wc *workerClient) connWorker(wk *worker.WorkerInfo) (*rpc.Client, error) {
	// 创建与WorkerService的链接
	conn, err := jsonrpc.Dial("tcp", wk.Address)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in call task jsonrpc.Dial.")
		return nil, err
	}

	if conn == nil {
		err = errors.New("Got an nil connection.")
		log.ErrorWithFields(log.Fields{
			"error": err.Error(),
		}, "Error in call task jsonrpc.Dial.")
		return nil, err
	}

	return conn, nil
}

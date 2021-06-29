package cli

import (
	"context"
	"eago/common/log"
	"eago/task/conf"
	"eago/task/worker"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
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

// List 按模块名列出所有活跃的Worker
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

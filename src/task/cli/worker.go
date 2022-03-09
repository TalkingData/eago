package cli

import (
	"context"
	"eago/common/log"
	"eago/common/tracer"
	"eago/task/conf"
	"eago/task/worker"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type AfterConnWrapHandler func(conn net.Conn)

var WorkerClient *workerClient

type workerClient struct {
	etcdCli           *clientv3.Client
	afterConnHandlers []AfterConnWrapHandler
}

// InitWorkerCli 启动Worker客户端
func InitWorkerCli(afterHandlers ...AfterConnWrapHandler) {
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
		etcdCli:           etcdCli,
		afterConnHandlers: afterHandlers,
	}
}

// List 列出所有活跃的Worker
func (wc *workerClient) List(ctx context.Context) []*worker.WorkerInfo {
	sp, c := tracer.StartSpanFromContext(ctx)
	defer sp.Finish()
	workers := make([]*worker.WorkerInfo, 0)

	resp, err := wc.etcdCli.Get(c, worker.WORKER_REGISTER_KEY_PREFFIX, clientv3.WithPrefix())
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
func (wc *workerClient) ListByModular(m string) []*worker.WorkerInfo {
	workers := make([]*worker.WorkerInfo, 0)

	for _, wk := range wc.List(context.Background()) {
		if wk.Modular != m {
			continue
		}
		workers = append(workers, wk)
	}

	return workers
}

// GetWorkerById 获得指定Worker
func (wc *workerClient) GetWorkerById(wkId string) *worker.WorkerInfo {
	for _, wk := range wc.List(context.Background()) {
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
	defer func() {
		_ = conn.Close()
	}()

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
		}, "An error occurred while workerClient.CallTask, Worker returned not ok status.")
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
	defer func() {
		_ = conn.Close()
	}()

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
		}, "An error occurred while workerClient.CallTask, Worker returned not ok status.")
		return err
	}

	return nil
}

// connWorker 连接Worker
func (wc *workerClient) connWorker(wk *worker.WorkerInfo) (*rpc.Client, error) {
	// 创建与Worker的链接
	conn, err := net.Dial("tcp", wk.Address)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while jsonrpc.Dial.")
		return nil, err
	}
	if conn == nil {
		err = errors.New("got an nil connection")
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while jsonrpc.Dial.")
		return nil, err
	}

	// 执行所有handler
	for _, h := range wc.afterConnHandlers {
		h(conn)
	}

	// 返回client
	return jsonrpc.NewClient(conn), nil
}

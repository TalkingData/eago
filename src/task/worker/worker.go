package worker

import (
	"context"
	"eago/common/log"
	proto "eago/task/srv/proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	uuid "github.com/satori/go.uuid"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Worker interface {
	RegTask(codename string, fn TaskFunc)

	// Start 启动Worker服务
	Start() error
	// Stop 关闭Worker服务
	Stop()
}

// NewWorker 创建Worker
func NewWorker(opts ...Option) Worker {
	w := newWorker(opts...)
	w.init()
	return w
}

func newWorker(opts ...Option) *worker {
	options := newOptions(opts...)
	return &worker{
		opts: options,
	}
}

type worker struct {
	workerId string

	endpoint string
	listener net.Listener

	taskList *taskList
	runList  *taskList

	etcdCli           *clientv3.Client
	etcdLease         clientv3.Lease
	taskServiceClient proto.TaskService

	mu sync.RWMutex

	opts      Options
	startTime *time.Time
}

// callTask 运行任务务
func (w *worker) callTask(req *CallTaskReq) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 查看当前Worker是否注册了调用的任务
	if !w.taskList.Exists(req.TaskCodename) {
		err := w.setTaskStatus(req.TaskUniqueId, TASK_WORKER_TASK_NOT_FOUND_ERROR_END_STATUS)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": req.TaskUniqueId,
				"error":          err,
			}, "An error occurred while taskServiceClient.SetTaskStatus.")
			return
		}
		return
	}

	// 查看调用的任务是否在运行
	if w.runList.Exists(req.TaskUniqueId) {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": req.TaskUniqueId,
		}, "An error occurred while worker.runTask, Task already running.")
		return
	}

	// 设置任务为Pending状态
	err := w.setTaskStatus(req.TaskUniqueId, TASK_PENDING_STATUS)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": req.TaskUniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		// 失败仅记录日志，不跳出
	}

	task := w.taskList.CopyGet(req.TaskCodename)
	ctx := context.Background()
	if req.Timeout > 0 {
		task.Cxt, task.Cancel = context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Second)
	} else {
		task.Cxt, task.Cancel = context.WithCancel(ctx)
	}

	task.logger = newLogger(w.opts.LogBufferSize)

	task.Param = &Param{
		req.TaskUniqueId,
		req.Caller,
		req.Timeout,
		req.Arguments,
		time.Now(),
		time.Unix(req.Timestamp, 0),
		task.logger,
	}

	w.runList.New(req.TaskUniqueId, &task)
	go func() {
		defer func() {
			// recover here for task panic end.
			if r := recover(); r != nil {
				task.logger.Error("Task panic: %s", r)
				log.ErrorWithFields(log.Fields{
					"task_unique_id": req.TaskUniqueId,
					"panic":          r,
				}, "An error occurred while task.Run, Task panic.")
				w.callback(&task, TASK_PANIC_END_STATUS)
			}
		}()

		// 设置任务为Running状态
		err := w.setTaskStatus(req.TaskUniqueId, TASK_RUNNING_STATUS)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": req.TaskUniqueId,
				"error":          err,
			}, "An error occurred while taskServiceClient.SetTaskStatus.")
			// 失败仅记录日志，不跳出
		}

		// 执行任务
		task.Run(func(ok bool) {
			if ok {
				w.callback(&task, TASK_SUCCESS_END_STATUS)
			} else {
				w.callback(&task, TASK_FAILED_END_STATUS)
			}
		})
	}()

	// 为新任务开启日志消费者
	go w.logConsumer(&task)
}

// killTask 杀任务
func (w *worker) killTask(taskUniqueId string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 查看要杀死的任务是否在运行
	if w.runList.Exists(taskUniqueId) {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
		}, "An error occurred while worker.killTask, Task is not in running state.")
		return
	}

	// 结束任务
	t := w.runList.Get(taskUniqueId)
	t.Cancel()

	// 设置任务为手动结束状态
	err := w.setTaskStatus(taskUniqueId, TASK_MANUAL_END_STATUS)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		// 失败仅记录日志，不跳出
	}

	// 关闭日志通道
	t.logger.Wg.Wait()
	close(t.logger.LogCh)

	// 删除任务
	w.runList.Del(taskUniqueId)
}

// callback 回调调度中心
func (w *worker) callback(task *Task, status int) {
	defer func() {
		// 通过context结束任务
		task.Cancel()

		// 删除任务
		w.runList.Del(task.Param.TaskUniqueId)
	}()

	// 关闭日志通道
	task.logger.Wg.Wait()
	close(task.logger.LogCh)

	err := w.setTaskStatus(task.Param.TaskUniqueId, status)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
			"status":         status,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		return
	}
}

// logConsumer 对应任务日志的消费并发送至Task模块
func (w *worker) logConsumer(task *Task) {
	log.InfoWithFields(log.Fields{
		"task_unique_id": task.Param.TaskUniqueId,
	}, "Worker log consumer started.")
	defer func() {
		log.InfoWithFields(log.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
		}, "Worker log consumer end.")
	}()

	stream, err := w.taskServiceClient.AppendTaskLog(context.Background())
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.TaskLog stream.")
		return
	}
	defer func() { _ = stream.Close() }()

	for {
		select {
		case content, ok := <-task.logger.LogCh:
			if !ok {
				log.InfoWithFields(log.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
				}, "Log channel was closed, The task may be done.")
				return
			}

			err := stream.Send(&proto.AppendTaskLogReq{TaskUniqueId: task.Param.TaskUniqueId, Content: *content})
			if err != nil {
				log.ErrorWithFields(log.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
					"error":          err,
				}, "An error occurred while stream.Send.")
				task.logger.Wg.Done()
				return
			}
			_, err = stream.Recv()
			if err != nil {
				log.ErrorWithFields(log.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
					"error":          err,
				}, "An error occurred while stream.Recv.")
				task.logger.Wg.Done()
				return
			}
			task.logger.Wg.Done()
		}
	}
}

// setTaskStatus 设置任务状态
func (w *worker) setTaskStatus(taskUniqueId string, status int) error {
	req := &proto.SetTaskStatusReq{TaskUniqueId: taskUniqueId, Status: int32(status)}
	res, err := w.taskServiceClient.SetTaskStatus(context.Background(), req)
	if err != nil {
		return err
	}

	if !res.Ok {
		return errors.New("unsuccessful rpc result")
	}

	return nil
}

// RegTask 注册任务
func (w *worker) RegTask(codename string, fn TaskFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var t = &Task{}
	t.Codename = codename
	t.fn = fn

	w.taskList.New(codename, t)
}

// Start 启动Worker
func (w *worker) Start() error {
	err := rpc.Register(&WorkerService{wk: w})
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while rpc.Register.")
		return err
	}

	// 开启RPC监听端口
	w.listener, err = net.Listen("tcp", fmt.Sprintf("%s:0", w.opts.WorkerIp))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while net.Listen.")
		return err
	}

	// 设置当前Worker的IP和端口(endpoint)信息
	w.endpoint = w.listener.Addr().String()
	log.Info(fmt.Sprintf("Worker %s starting at %s", w.workerId, w.endpoint))

	go func() {
		for {
			conn, err := w.listener.Accept()
			if err != nil {
				log.ErrorWithFields(log.Fields{
					"error": err,
				}, "An error occurred while listener.Accept.")
				continue
			}
			// 处理实际请求
			go jsonrpc.ServeConn(conn)
		}
	}()

	// 向注册中心注册Worker
	w.register()

	// 等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit
	w.Stop()
	return nil
}

// init 初始化Worker
func (w *worker) init() {
	// 生成WorkerId
	if w.opts.MultiInstance {
		w.workerId = fmt.Sprintf("%s.worker-%s", w.opts.ServiceName, uuid.NewV4().String())
	} else {
		w.workerId = fmt.Sprintf("%s.worker-unique", w.opts.ServiceName)
	}

	w.taskList = &taskList{
		data: make(map[string]*Task),
	}
	w.runList = &taskList{
		data: make(map[string]*Task),
	}

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(w.opts.EtcdAddresses...),
		etcdv3.Auth(w.opts.EtcdUsername, w.opts.EtcdPassword),
	)
	cli := micro.NewService(micro.Registry(etcdReg))
	w.taskServiceClient = proto.NewTaskService(w.opts.TaskRpcServiceName, cli.Client())

	etcdCfg := clientv3.Config{
		Endpoints:   w.opts.EtcdAddresses,
		Username:    w.opts.EtcdUsername,
		Password:    w.opts.EtcdPassword,
		DialTimeout: 5 * time.Second,
	}

	etcdCli, err := clientv3.New(etcdCfg)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while clientv3.New.")
		panic(err)
	}

	w.etcdCli = etcdCli
}

// Stop 关闭Worker服务
func (w *worker) Stop() {
	log.Info(fmt.Sprintf("Worker %s Stop called.", w.workerId))
	defer log.Info(fmt.Sprintf("Worker %s Stop end.", w.workerId))

	// 等待所有任务结束

	w.unregister()

	_ = w.etcdCli.Close()
	_ = w.listener.Close()
}

// register
func (w *worker) register() {
	log.Info(fmt.Sprintf("Worker %s register called.", w.workerId))
	defer log.Info(fmt.Sprintf("Worker %s register end.", w.workerId))

	ctx := context.TODO()

	// 设置worker开始时间
	n := time.Now()
	w.startTime = &n

	// 建立etcd租约
	w.etcdLease = clientv3.NewLease(w.etcdCli)
	leaseGrantResp, err := w.etcdLease.Grant(ctx, w.opts.RegisterTtl)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while lease.Grant.")
		panic(err)
	}

	ch, err := w.etcdLease.KeepAlive(ctx, leaseGrantResp.ID)
	// 续约应答
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				log.Info("Worker regCh may closed.")
				break
			}
		}
	}()

	// 生成注册Key和Value
	regK := fmt.Sprintf("%s/%s/%s", WORKER_REGISTER_KEY_PREFFIX, w.opts.ServiceName, w.workerId)
	regV, _ := json.Marshal(WorkerInfo{
		Modular:   w.opts.ServiceName,
		Address:   w.endpoint,
		WorkerId:  w.workerId,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	txn := clientv3.NewKV(w.etcdCli).Txn(ctx)
	// 注册到etcd
	txn.If(clientv3.Compare(clientv3.CreateRevision(regK), "=", 0)).
		Then(clientv3.OpPut(regK, string(regV), clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(regK))

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		panic(err)
	}

	if !txnResp.Succeeded {
		if !w.opts.MultiInstance {
			panic("another worker instance is already running, and worker not allowed create multi instance")
		}
		panic("another worker instance is already running using same worker id")
	}
}

// unregister
func (w *worker) unregister() {
	log.Info(fmt.Sprintf("Worker %s unregister called.", w.workerId))
	defer log.Info(fmt.Sprintf("Worker %s unregister end.", w.workerId))

	if w.etcdLease != nil {
		_ = w.etcdLease.Close()
	}

	k := fmt.Sprintf("/%s/%s/%s", WORKER_REGISTER_KEY_PREFFIX, w.opts.ServiceName, w.workerId)
	_, err := w.etcdCli.Delete(context.Background(), k)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while etcdCli.Delete, This error will be ignored.")
	}

	w.startTime = nil
}

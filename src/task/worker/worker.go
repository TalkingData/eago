package worker

import (
	"context"
	"eago/common/log"
	proto "eago/task/srv/proto"
	workerProto "eago/task/worker/proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net"
	"sync"
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

	endpoint   string
	listener   net.Listener
	grpcServer *grpc.Server

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
func (wk *worker) callTask(codename, uniqueId, args string, timeout int32, caller string, ts int64) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	// 查看当前Worker是否注册了调用的任务
	if !wk.taskList.Exists(codename) {
		err := wk.setTaskStatus(uniqueId, TASK_WORKER_TASK_NOT_FOUND_ERROR_END_STATUS)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": uniqueId,
				"error":          err,
			}, "An error occurred while taskServiceClient.SetTaskStatus.")
			return
		}
		return
	}

	// 查看调用的任务是否在运行
	if wk.runList.Exists(uniqueId) {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": uniqueId,
		}, "An error occurred while worker.runTask, Task already running.")
		return
	}

	// 设置任务为Pending状态
	err := wk.setTaskStatus(uniqueId, TASK_PENDING_STATUS)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": uniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		// 失败仅记录日志，不跳出
	}

	task := wk.taskList.CopyGet(codename)
	ctx := context.Background()
	if timeout > 0 {
		task.Cxt, task.Cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	} else {
		task.Cxt, task.Cancel = context.WithCancel(ctx)
	}

	task.logger = newLogger(wk.opts.LogBufferSize)

	task.Param = &Param{
		TaskUniqueId:    uniqueId,
		Caller:          caller,
		Timeout:         timeout,
		Arguments:       args,
		LocalStartTime:  time.Now(),
		RemoteStartTime: time.Unix(ts, 0),
		Log:             task.logger,
	}

	wk.runList.Put(uniqueId, &task)
	go func() {
		defer func() {
			// recover here for task panic end.
			if r := recover(); r != nil {
				task.logger.Error("Task panic: %s", r)
				log.ErrorWithFields(log.Fields{
					"task_unique_id": uniqueId,
					"panic":          r,
				}, "An error occurred while task.Run, Task panic.")
				wk.callback(&task, TASK_PANIC_END_STATUS)
			}
		}()

		// 设置任务为Running状态
		err := wk.setTaskStatus(uniqueId, TASK_RUNNING_STATUS)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"task_unique_id": uniqueId,
				"error":          err,
			}, "An error occurred while taskServiceClient.SetTaskStatus.")
			// 失败仅记录日志，不跳出
		}

		// 执行任务
		task.Run(func(err error) {
			switch err {
			case nil:
				// 成功结束
				wk.callback(&task, TASK_SUCCESS_END_STATUS)
			case context.DeadlineExceeded:
				// 超时结束
				wk.callback(&task, TASK_TIMEOUT_END_STATUS)
			case context.Canceled:
				// 任务取消
				wk.callback(&task, TASK_MANUAL_END_STATUS)
			default:
				// 任务失败结束
				wk.callback(&task, TASK_FAILED_END_STATUS)
			}
		})
	}()

	// 为新任务开启日志消费者
	go wk.logConsumer(&task)
}

// killTask 杀任务
func (wk *worker) killTask(taskUniqueId string) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	// 查看要杀死的任务是否在运行
	if wk.runList.Exists(taskUniqueId) {
		log.ErrorWithFields(log.Fields{
			"task_unique_id": taskUniqueId,
		}, "An error occurred while worker.killTask, Task is not in running state.")
		return
	}

	// 结束任务
	t := wk.runList.Get(taskUniqueId)
	t.Cancel()

	// 设置任务为手动结束状态
	err := wk.setTaskStatus(taskUniqueId, TASK_MANUAL_END_STATUS)
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
	wk.runList.Delete(taskUniqueId)
}

// callback 回调调度中心
func (wk *worker) callback(task *Task, status int) {
	defer func() {
		// 通过context结束任务
		task.Cancel()

		// 删除任务
		wk.runList.Delete(task.Param.TaskUniqueId)
	}()

	// 关闭日志通道
	task.logger.Wg.Wait()
	close(task.logger.LogCh)

	err := wk.setTaskStatus(task.Param.TaskUniqueId, status)
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
func (wk *worker) logConsumer(task *Task) {
	log.InfoWithFields(log.Fields{
		"task_unique_id": task.Param.TaskUniqueId,
	}, "Worker log consumer started.")
	defer func() {
		log.InfoWithFields(log.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
		}, "Worker log consumer end.")
	}()

	stream, err := wk.taskServiceClient.AppendTaskLog(context.Background())
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
func (wk *worker) setTaskStatus(taskUniqueId string, status int) error {
	req := &proto.SetTaskStatusReq{TaskUniqueId: taskUniqueId, Status: int32(status)}
	res, err := wk.taskServiceClient.SetTaskStatus(context.Background(), req)
	if err != nil {
		return err
	}

	if !res.Ok {
		return errors.New("unsuccessful rpc result")
	}

	return nil
}

// RegTask 注册任务
func (wk *worker) RegTask(codename string, fn TaskFunc) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	var t = &Task{}
	t.Codename = codename
	t.fn = fn

	wk.taskList.Put(codename, t)
}

// Start 启动Worker
func (wk *worker) Start() (err error) {
	// 开启RPC监听端口
	wk.listener, err = net.Listen("tcp", fmt.Sprintf("%s:0", wk.opts.WorkerIp))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while net.Listen.")
		return err
	}

	// 创建worker的grpc server
	wk.grpcServer = grpc.NewServer()
	// 注册GRPC服务
	workerProto.RegisterTaskWorkerServiceServer(wk.grpcServer, NewTaskWorkerService(wk))

	// 设置当前Worker的IP和端口(endpoint)信息
	wk.endpoint = wk.listener.Addr().String()
	log.Info(fmt.Sprintf("Worker %s starting at %s", wk.workerId, wk.endpoint))

	// 向注册中心注册Worker
	wk.register()

	return wk.grpcServer.Serve(wk.listener)
}

// init 初始化Worker
func (wk *worker) init() {
	// 生成WorkerId
	if wk.opts.MultiInstance {
		wk.workerId = fmt.Sprintf("%s.worker-%s", wk.opts.ServiceName, uuid.NewV4().String())
	} else {
		wk.workerId = fmt.Sprintf("%s.worker-unique", wk.opts.ServiceName)
	}

	wk.taskList = NewTaskList()
	wk.runList = NewTaskList()

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(wk.opts.EtcdAddresses...),
		etcdv3.Auth(wk.opts.EtcdUsername, wk.opts.EtcdPassword),
	)
	cli := micro.NewService(micro.Registry(etcdReg))
	wk.taskServiceClient = proto.NewTaskService(wk.opts.TaskRpcServiceName, cli.Client())

	etcdCfg := clientv3.Config{
		Endpoints:   wk.opts.EtcdAddresses,
		Username:    wk.opts.EtcdUsername,
		Password:    wk.opts.EtcdPassword,
		DialTimeout: 5 * time.Second,
	}

	etcdCli, err := clientv3.New(etcdCfg)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while clientv3.New.")
		panic(err)
	}

	wk.etcdCli = etcdCli
}

// Stop 关闭Worker服务
func (wk *worker) Stop() {
	log.Info(fmt.Sprintf("Worker %s Stop called.", wk.workerId))
	defer log.Info(fmt.Sprintf("Worker %s Stop end.", wk.workerId))

	// 等待所有任务结束

	wk.unregister()

	if wk.grpcServer != nil {
		wk.grpcServer.Stop()
	}
	if wk.etcdCli != nil {
		_ = wk.etcdCli.Close()
	}
	if wk.listener != nil {
		_ = wk.listener.Close()
	}
}

// register
func (wk *worker) register() {
	log.Info(fmt.Sprintf("Worker %s register called.", wk.workerId))
	defer log.Info(fmt.Sprintf("Worker %s register end.", wk.workerId))

	ctx := context.TODO()

	// 设置worker开始时间
	n := time.Now()
	wk.startTime = &n

	// 建立etcd租约
	wk.etcdLease = clientv3.NewLease(wk.etcdCli)
	leaseGrantResp, err := wk.etcdLease.Grant(ctx, wk.opts.RegisterTtl)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while lease.Grant.")
		panic(err)
	}

	ch, err := wk.etcdLease.KeepAlive(ctx, leaseGrantResp.ID)
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
	regK := fmt.Sprintf("%s/%s/%s", WORKER_REGISTER_KEY_PREFFIX, wk.opts.ServiceName, wk.workerId)
	regV, _ := json.Marshal(WorkerInfo{
		Modular:   wk.opts.ServiceName,
		Address:   wk.endpoint,
		WorkerId:  wk.workerId,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	txn := clientv3.NewKV(wk.etcdCli).Txn(ctx)
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
		if !wk.opts.MultiInstance {
			panic("another worker instance is already running, and worker not allowed create multi instance")
		}
		panic("another worker instance is already running using same worker id")
	}
}

// unregister
func (wk *worker) unregister() {
	log.Info(fmt.Sprintf("Worker %s unregister called.", wk.workerId))
	defer log.Info(fmt.Sprintf("Worker %s unregister end.", wk.workerId))

	if wk.etcdLease != nil {
		_ = wk.etcdLease.Close()
	}

	k := fmt.Sprintf("/%s/%s/%s", WORKER_REGISTER_KEY_PREFFIX, wk.opts.ServiceName, wk.workerId)
	_, err := wk.etcdCli.Delete(context.Background(), k)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while etcdCli.Delete, This error will be ignored.")
	}

	wk.startTime = nil
}

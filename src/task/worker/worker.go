package worker

import (
	"context"
	"eago/cli"
	"eago/common/global"
	"eago/common/logger"
	"eago/task/dto"
	taskpb "eago/task/proto"
	workerpb "eago/task/worker/proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Worker interface {
	RegTask(codename string, fn TaskFunc)

	// Start 启动Worker服务
	Start() error
	// Stop 关闭Worker服务
	Stop()
}

type worker struct {
	workerId string

	endpoint  string
	listener  net.Listener
	workerSrv *grpc.Server
	taskCli   taskpb.TaskService

	taskList    *taskList
	runningList *taskList

	etcdCli   *clientv3.Client
	etcdLease clientv3.Lease

	mu sync.RWMutex

	logger *logger.Logger

	opts        Options
	startTime   *time.Time
	runningFlag int32
}

// NewWorker 创建Worker
func NewWorker(options ...Option) Worker {
	return newWorker(options...)
}

func newWorker(options ...Option) *worker {
	opts := newOptions(options...)

	// 生成WorkerId
	wkId := ""
	if opts.MultiInstance {
		wkId = fmt.Sprintf("%s.worker-%s", opts.ServiceName, uuid.NewV4().String())
	} else {
		wkId = fmt.Sprintf("%s.worker-unique", opts.ServiceName)
	}

	etcdCfg := clientv3.Config{
		Endpoints:   opts.EtcdAddresses,
		Username:    opts.EtcdUsername,
		Password:    opts.EtcdPassword,
		DialTimeout: 5 * time.Second,
	}
	etcdCli, err := clientv3.New(etcdCfg)
	if err != nil {
		opts.Logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while clientv3.New in newWorker.")
		panic(err)
	}

	return &worker{
		workerId: wkId,

		taskCli: cli.NewTaskClient(opts.EtcdUsername, opts.EtcdPassword, opts.EtcdAddresses),

		taskList:    NewTaskList(),
		runningList: NewTaskList(),

		etcdCli: etcdCli,

		logger: opts.Logger,

		opts: opts,
	}
}

// callTask 运行任务务
func (wk *worker) callTask(codename, uniqueId, args string, timeout int64, caller string, ts int64) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	// 查看当前Worker是否注册了调用的任务
	if !wk.taskList.Exists(codename) {
		if err := wk.setTaskStatus(uniqueId, dto.TaskResultStatusWorkerTaskNotFoundErrEnd); err != nil {
			wk.logger.ErrorWithFields(logger.Fields{
				"task_unique_id": uniqueId,
				"error":          err,
			}, "An error occurred while taskServiceClient.SetTaskStatus.")
			return
		}
		return
	}

	// 查看调用的任务是否在运行
	if wk.runningList.Exists(uniqueId) {
		wk.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": uniqueId,
		}, "An error occurred while worker.runTask, Task already running.")
		return
	}

	// 设置任务为Pending状态
	if err := wk.setTaskStatus(uniqueId, dto.TaskResultStatusPending); err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
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

	task.logger = newResultLog(wk.opts.ResultLogBufferSize, wk.opts.PrintResultLog)

	task.Param = &Param{
		TaskUniqueId:    uniqueId,
		Caller:          caller,
		Timeout:         timeout,
		Arguments:       args,
		LocalStartTime:  time.Now(),
		RemoteStartTime: time.Unix(ts, 0),
		Log:             task.logger,
	}

	wk.runningList.Put(uniqueId, &task)
	go func() {
		defer func() {
			// recover here for task panic end.
			if r := recover(); r != nil {
				task.logger.Error("Task panic: %s", r)
				wk.logger.ErrorWithFields(logger.Fields{
					"task_unique_id": uniqueId,
					"panic":          r,
				}, "An error occurred while task.Run, Task panic.")
				wk.callback(&task, dto.TaskResultStatusPanicEnd)
			}
		}()

		// 设置任务为Running状态
		if err := wk.setTaskStatus(uniqueId, dto.TaskResultStatusRunning); err != nil {
			wk.logger.ErrorWithFields(logger.Fields{
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
				wk.callback(&task, dto.TaskResultStatusSuccessEnd)
			case context.DeadlineExceeded:
				// 超时结束
				wk.callback(&task, dto.TaskResultStatusTimeoutEnd)
			case context.Canceled:
				// 任务取消
				wk.callback(&task, dto.TaskResultStatusManualEnd)
			default:
				// 任务失败结束
				wk.callback(&task, dto.TaskResultStatusFailedEnd)
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
	if wk.runningList.Exists(taskUniqueId) {
		wk.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
		}, "An error occurred while worker.killTask, Task is not in running state.")
		return
	}

	// 结束任务
	t := wk.runningList.Get(taskUniqueId)
	t.Cancel()

	// 设置任务为手动结束状态
	if err := wk.setTaskStatus(taskUniqueId, dto.TaskResultStatusManualEnd); err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": taskUniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		// 失败仅记录日志，不跳出
	}

	// 关闭日志通道
	t.logger.wg.Wait()
	close(t.logger.logCh)

	// 删除任务
	wk.runningList.Delete(taskUniqueId)
}

// callback 回调调度中心
func (wk *worker) callback(task *Task, status int) {
	defer func() {
		// 通过context结束任务
		task.Cancel()

		// 删除任务
		wk.runningList.Delete(task.Param.TaskUniqueId)
	}()

	// 关闭日志通道
	task.logger.wg.Wait()
	close(task.logger.logCh)

	if err := wk.setTaskStatus(task.Param.TaskUniqueId, status); err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
			"status":         status,
			"error":          err,
		}, "An error occurred while taskServiceClient.SetTaskStatus.")
		return
	}
}

// logConsumer 对应任务日志的消费并发送至Task模块
func (wk *worker) logConsumer(task *Task) {
	wk.logger.InfoWithFields(logger.Fields{
		"task_unique_id": task.Param.TaskUniqueId,
	}, "Worker log consumer started.")
	defer func() {
		wk.logger.InfoWithFields(logger.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
		}, "Worker log consumer end.")
	}()

	stream, err := wk.taskCli.AppendTaskLog(context.Background())
	if err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"task_unique_id": task.Param.TaskUniqueId,
			"error":          err,
		}, "An error occurred while taskServiceClient.TaskLog stream.")
		return
	}
	defer func() { _ = stream.Close() }()

	for {
		select {
		case content, ok := <-task.logger.logCh:
			if !ok {
				wk.logger.InfoWithFields(logger.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
				}, "ResultLog channel was closed, The task may be done.")
				return
			}

			err = stream.Send(&taskpb.AppendTaskLogReq{TaskUniqueId: task.Param.TaskUniqueId, Content: *content})
			if err != nil {
				wk.logger.ErrorWithFields(logger.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
					"error":          err,
				}, "An error occurred while stream.Send.")
				task.logger.wg.Done()
				return
			}
			if _, err = stream.Recv(); err != nil {
				wk.logger.ErrorWithFields(logger.Fields{
					"task_unique_id": task.Param.TaskUniqueId,
					"error":          err,
				}, "An error occurred while stream.Recv.")
				task.logger.wg.Done()
				return
			}
			task.logger.wg.Done()
		}
	}
}

// setTaskStatus 设置任务状态
func (wk *worker) setTaskStatus(taskUniqueId string, status int) error {
	req := &taskpb.SetResultStatusReq{TaskUniqueId: taskUniqueId, Status: int32(status)}
	if _, err := wk.taskCli.SetResultStatus(context.Background(), req); err != nil {
		return err
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
	wk.logger.Info(fmt.Sprintf("Starting %s worker...", wk.opts.ServiceName))
	// 设置Worker为运行状态
	wk.setRunningFlag(true)

	// 开启RPC监听端口
	wk.listener, err = net.Listen("tcp", fmt.Sprintf("%s:0", wk.opts.WorkerIp))
	if err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while net.Listen.")
		return err
	}

	// 创建worker的grpc server
	wk.workerSrv = grpc.NewServer(grpc.UnaryInterceptor(wk.isValidSrvTokenInterceptor))
	// 注册GRPC服务
	workerpb.RegisterTaskWorkerServiceServer(wk.workerSrv, NewTaskWorkerService(wk, wk.logger))

	// 设置当前Worker的IP和端口(endpoint)信息
	wk.endpoint = wk.listener.Addr().String()
	wk.logger.Info(fmt.Sprintf("Worker %s listening on %s.", wk.workerId, wk.endpoint))

	// 向注册中心注册Worker
	if err = wk.register(); err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "Got critical error when worker starting.")
		panic(err)
	}

	return wk.workerSrv.Serve(wk.listener)
}

// Stop 关闭Worker服务
func (wk *worker) Stop() {
	wk.logger.Info(fmt.Sprintf("Worker %s Stop called.", wk.workerId))
	defer wk.logger.Info(fmt.Sprintf("Worker %s Stop end.", wk.workerId))

	// 设置Worker为结束状态
	wk.setRunningFlag(false)
	// 等待所有任务结束
	wk.unregister()

	if wk.workerSrv != nil {
		wk.workerSrv.Stop()
	}
	if wk.etcdCli != nil {
		_ = wk.etcdCli.Close()
	}
	if wk.listener != nil {
		_ = wk.listener.Close()
	}
}

func (wk *worker) register() error {
	wk.logger.Info(fmt.Sprintf("Worker %s register called.", wk.workerId))
	defer wk.logger.Info(fmt.Sprintf("Worker %s register end.", wk.workerId))

	ctx := context.Background()

	// 设置worker开始时间
	n := time.Now()
	wk.startTime = &n

	// 建立etcd租约
	wk.etcdLease = clientv3.NewLease(wk.etcdCli)
	leaseGrantResp, err := wk.etcdLease.Grant(ctx, wk.opts.RegisterTtl)
	if err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while lease.Grant.")
		return err
	}

	ch, err := wk.etcdLease.KeepAlive(ctx, leaseGrantResp.ID)
	// 续约应答
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				wk.logger.Error("Worker regCh may closed, it will register soon.")
				wk.unregister()
				for {
					// 若Worker不在运行，则退出
					if !wk.isRunning() {
						return
					}
					// 若Worker仍在运行，需要重新注册Worker
					if err = wk.register(); err == nil {
						return
					}
					time.Sleep(time.Second)
				}
			}
			time.Sleep(time.Second)
		}
	}()

	// 生成注册Key和Value
	regK := fmt.Sprintf("%s/%s/%s", WorkerRegisterKeyPrefix, wk.opts.ServiceName, wk.workerId)
	regV, _ := json.Marshal(dto.WorkerInfo{
		Modular:   wk.opts.ServiceName,
		Address:   wk.endpoint,
		WorkerId:  wk.workerId,
		StartTime: time.Now().Format(global.TimestampFormat),
	})

	txn := clientv3.NewKV(wk.etcdCli).Txn(ctx)
	// 注册到etcd
	txn.If(clientv3.Compare(clientv3.CreateRevision(regK), "=", 0)).
		Then(clientv3.OpPut(regK, string(regV), clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(regK))

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		if !wk.opts.MultiInstance {
			return errors.New(
				"another worker instance is already running, and worker not allowed create multi instance",
			)
		}
		return errors.New("another worker instance is already running using same worker id")
	}

	return nil
}

func (wk *worker) unregister() {
	wk.logger.Info(fmt.Sprintf("Worker %s unregister called.", wk.workerId))
	defer wk.logger.Info(fmt.Sprintf("Worker %s unregister end.", wk.workerId))

	if wk.etcdLease != nil {
		_ = wk.etcdLease.Close()
	}

	k := fmt.Sprintf("/%s/%s/%s", WorkerRegisterKeyPrefix, wk.opts.ServiceName, wk.workerId)
	_, err := wk.etcdCli.Delete(context.Background(), k)
	if err != nil {
		wk.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while etcdCli.Delete in worker.unregister, This error will be ignored.")
	}

	wk.startTime = nil
}

// isValidSrvTokenInterceptor 验证拦截器
func (wk *worker) isValidSrvTokenInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler,
) (interface{}, error) {
	wk.logger.Debug("worker.isValidSrvTokenInterceptor called.")
	defer wk.logger.Debug("worker.isValidSrvTokenInterceptor end.")

	// 记录请求者ip和请求的方法
	logF := logger.Fields{"method": info.FullMethod}
	if p, ok := peer.FromContext(ctx); ok {
		logF["from_address"] = p.Addr.String()
	}
	wk.logger.InfoWithFields(logF, "Got grpc call.")

	// 获取ctx中的metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		wk.logger.Error("Task worker get context metadata failed")
		return nil, status.Errorf(codes.Unauthenticated, "get context metadata failed")
	}

	mdVal, ok := md["srv_token"]
	if !ok {
		wk.logger.Error("Task worker got empty token")
		return nil, status.Errorf(codes.Unauthenticated, "got empty token")
	}
	srvTk := mdVal[0]

	// 调用taskCli进行认证
	boolVal, err := wk.taskCli.IsValidSrvToken(
		ctx,
		&taskpb.SrvTokenQuery{Value: srvTk},
	)
	if err != nil {
		// 验证失败，远端返回了错误值
		wk.logger.WarnWithFields(logger.Fields{
			"srv_token": srvTk,
			"error":     err,
		}, "taskServiceClient.IsValidSrvToken returned error.")
		return nil, status.Errorf(
			codes.Unauthenticated, "taskServiceClient.IsValidSrvToken returned error: %s", err.Error(),
		)
	}
	if boolVal == nil {
		// 验证失败，远端返回了错误值
		wk.logger.WarnWithFields(logger.Fields{
			"srv_token": srvTk,
		}, "taskServiceClient.IsValidSrvToken returned nil.")
		return nil, status.Errorf(codes.Unauthenticated, "taskServiceClient.IsValidSrvToken returned nil")
	}
	if !boolVal.Value {
		// 验证失败，结束
		wk.logger.WarnWithFields(logger.Fields{
			"srv_token": srvTk,
		}, "taskServiceClient.IsValidSrvToken returned false.")
		return nil, status.Errorf(codes.Unauthenticated, "taskServiceClient.IsValidSrvToken returned false")
	}

	// 验证成功，继续处理请求
	return h(ctx, req)
}

// setRunningFlag 设置worker运行标志
func (wk *worker) setRunningFlag(b bool) {
	if b {
		atomic.StoreInt32(&wk.runningFlag, 1)
		return
	}
	atomic.StoreInt32(&wk.runningFlag, 0)
}

// isRunning 判断worker是否在运行
func (wk *worker) isRunning() bool {
	return atomic.LoadInt32(&wk.runningFlag) != 0
}

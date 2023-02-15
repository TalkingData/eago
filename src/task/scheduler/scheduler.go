package main

import (
	"context"
	"eago/common/global"
	"eago/common/logger"
	commonpb "eago/common/proto"
	taskpb "eago/task/proto"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/go-basic/ipv4"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/robfig/cron"
	"time"
)

type scheduler struct {
	cron      *cron.Cron
	etcdCli   *clientv3.Client
	etcdLease clientv3.Lease
	taskCli   taskpb.TaskService

	logger *logger.Logger

	started bool
	opts    Options

	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewScheduler 创建Scheduler
func NewScheduler(ctx context.Context, options ...Option) *scheduler {
	opts := newOptions(options...)

	schCtx, schCancel := context.WithCancel(ctx)

	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(opts.EtcdAddresses...),
		etcdv3.Auth(opts.EtcdUsername, opts.EtcdPassword),
	)
	cli := micro.NewService(micro.Registry(etcdReg))
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   opts.EtcdAddresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return &scheduler{
		cron:    cron.New(),
		etcdCli: etcdCli,
		taskCli: taskpb.NewTaskService(opts.TaskRpcRegisterKey, cli.Client()),

		logger: opts.Logger,

		started: false,
		opts:    opts,

		ctx:        schCtx,
		cancelFunc: schCancel,
	}
}

// Start 启动计划任务
func (s *scheduler) Start() error {
	s.logger.Info("Starting task scheduler...")
	if s.started {
		panic("current instance is already started")
	}

	// 循环创建计划任务
	for _, sch := range s.listScheduleTasks() {
		// 创建计划任务
		err := s.cron.AddFunc(sch.Expression, func() {
			ctx := context.Background()
			req := &taskpb.CallTaskReq{
				TaskCodename: sch.TaskCodename,
				Timeout:      sch.Timeout,
				Arguments:    []byte(sch.Arguments),
				Caller:       "task.scheduler",
			}
			// 调用任务
			rsp, err := s.taskCli.CallTask(ctx, req)
			if err != nil {
				s.logger.ErrorWithFields(logger.Fields{
					"task_codename": sch.TaskCodename,
					"expression":    sch.Expression,
					"timeout":       sch.Timeout,
					"arguments":     sch.Arguments,
					"error":         err,
				}, "An error occurred while taskCli.CallTask in given task.")
				return
			}
			s.logger.InfoWithFields(logger.Fields{
				"task_codename":  sch.TaskCodename,
				"expression":     sch.Expression,
				"timeout":        sch.Timeout,
				"arguments":      sch.Arguments,
				"task_unique_id": rsp.TaskUniqueId,
			}, "Call task success.")
		})
		// 创建计划任务失败
		if err != nil {
			s.logger.ErrorWithFields(logger.Fields{
				"task_codename": sch.TaskCodename,
				"expression":    sch.Expression,
				"arguments":     sch.Arguments,
				"error":         err,
			}, "An error occurred while cron.AddFunc in scheduler.Start.")
			panic(fmt.Errorf("failed to add func, error: %w", err))
		}
		s.logger.InfoWithFields(logger.Fields{
			"task_codename": sch.TaskCodename,
			"expression":    sch.Expression,
			"timeout":       sch.Timeout,
			"arguments":     sch.Arguments,
		}, "scheduler task added.")
	}

	s.started = true
	s.register()
	s.cron.Start()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.InfoWithFields(logger.Fields{
				"context_error": s.ctx.Err(),
			}, "Scheduler stopped by context done.")
			return s.ctx.Err()
		}
	}

}

// Stop 停止计划任务
func (s *scheduler) Stop() {
	defer func() {
		s.ctx.Done()
		s.started = false
	}()

	if s.cron != nil {
		s.cron.Stop()
	}
	s.unregister()
	_ = s.etcdCli.Close()
}

// listScheduleTasks 获得已配置的计划任务
func (s *scheduler) listScheduleTasks() []*taskpb.Schedule {
	s.logger.Info("scheduler.listScheduleTasks called.")
	defer s.logger.Info("scheduler.listScheduleTasks end.")

	var maxPg uint32 = 2

	res := make([]*taskpb.Schedule, 0)
	for pg := uint32(1); pg < maxPg; pg++ {
		req := &commonpb.QueryWithPage{
			Page:     pg,
			PageSize: defaultClientMaxPageSize,
		}
		rsp, err := s.taskCli.PagedListSchedules(s.ctx, req)
		if err != nil {
			s.logger.ErrorWithFields(logger.Fields{
				"page":  pg,
				"error": err,
			}, "An error occurred while taskCli.PagedListSchedules in scheduler.listScheduleTasks.")
			panic(fmt.Errorf("failed to list schedule tasks, error: %s", err.Error()))
		}
		maxPg = rsp.Pages
		for _, r := range rsp.Schedules {
			s.logger.DebugWithFields(logger.Fields{
				"task_codename": r.TaskCodename,
				"expression":    r.Expression,
				"timeout":       r.Timeout,
				"arguments":     r.Arguments,
				"disabled":      r.Disabled,
			}, "Got a schedule.")

			// 跳过禁用的计划任务
			if r.Disabled {
				s.logger.Debug("Skip disabled schedule.")
				continue
			}

			res = append(res, &taskpb.Schedule{
				TaskCodename: r.TaskCodename,
				Expression:   r.Expression,
				Timeout:      r.Timeout,
				Arguments:    r.Arguments,
			})
		}
	}
	return res
}

func (s *scheduler) register() {
	s.logger.Info("scheduler.register called.")
	defer s.logger.Info("scheduler.register end.")

	// 建立etcd租约
	s.etcdLease = clientv3.NewLease(s.etcdCli)
	leaseGrantResp, err := s.etcdLease.Grant(s.ctx, s.opts.RegisterTtl)
	if err != nil {
		s.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while lease.Grant in scheduler.register.")
		panic(err)
	}

	ch, err := s.etcdLease.KeepAlive(s.ctx, leaseGrantResp.ID)
	// 续约应答
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				s.logger.Info("scheduler regCh may closed.")
				break
			}
		}
	}()

	// 生成注册Value
	regV, _ := json.Marshal(ScheduleInfo{
		IpAddress: ipv4.LocalIP(),
		StartTime: time.Now().Format(global.TimestampFormat),
	})

	txn := clientv3.NewKV(s.etcdCli).Txn(s.ctx)
	// 注册到etcd
	txn.If(clientv3.Compare(clientv3.CreateRevision(defaultSchedulerRegisterKey), "=", 0)).
		Then(clientv3.OpPut(defaultSchedulerRegisterKey, string(regV), clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(defaultSchedulerRegisterKey))

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		panic(err)
	}

	if !txnResp.Succeeded {
		panic("another scheduler instance is already running")
	}
}

func (s *scheduler) unregister() {
	s.logger.Info("scheduler.unregister called.")
	defer s.logger.Info("scheduler.unregister end.")

	if s.etcdLease != nil {
		_ = s.etcdLease.Close()
	}

	_, err := s.etcdCli.Delete(s.ctx, defaultSchedulerRegisterKey)
	if err != nil {
		s.logger.ErrorWithFields(logger.Fields{
			"error": err,
		}, "An error occurred while etcdCli.Delete in scheduler.unregister.")
	}
}

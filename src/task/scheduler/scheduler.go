package main

import (
	"context"
	"eago/common/log"
	proto "eago/task/srv/proto"
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

const (
	SCHEDULER_REGISTER_KEY = "/td/eago/scheduler"
	CLIENT_MAX_PAGESIZE    = 500
)

type Scheduler struct {
	cron          *cron.Cron
	etcdCli       *clientv3.Client
	etcdLease     clientv3.Lease
	taskSrvClient proto.TaskService

	started bool
	opts    Options
}

// NewScheduler 创建一个Scheduler
func NewScheduler(opts ...Option) *Scheduler {
	options := newOptions(opts...)
	s := &Scheduler{
		opts: options,
	}
	s.init()
	return s
}

// init 初始化
func (s *Scheduler) init() {
	etcdReg := etcdv3.NewRegistry(
		registry.Addrs(s.opts.EtcdAddresses...),
		etcdv3.Auth(s.opts.EtcdUsername, s.opts.EtcdPassword),
	)
	cli := micro.NewService(micro.Registry(etcdReg))
	s.taskSrvClient = proto.NewTaskService(s.opts.TaskRpcRegisterKey, cli.Client())

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.opts.EtcdAddresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while clientv3.New.")
		panic(err)
	}
	s.etcdCli = etcdCli

	s.cron = cron.New()
	s.started = false
}

// Start 启动计划任务
func (s *Scheduler) Start() {
	log.Info("Starting scheduler.")
	if s.started {
		panic("current instance is already started")
	}
	defer log.Info("Scheduler started.")

	// 循环创建计划任务
	for _, tmp := range s.getScheduleTask() {
		var sch = tmp
		// 创建计划任务
		err := s.cron.AddFunc(sch.Expression, func() {
			ctx := context.Background()
			req := &proto.CallTaskReq{
				TaskCodename: sch.TaskCodename,
				Timeout:      sch.Timeout,
				Arguments:    []byte(sch.Arguments),
				Caller:       "task.scheduler",
			}
			// 调用任务
			rsp, err := s.taskSrvClient.CallTask(ctx, req)
			if err != nil {
				log.ErrorWithFields(log.Fields{
					"task_codename": sch.TaskCodename,
					"expression":    sch.Expression,
					"timeout":       sch.Timeout,
					"arguments":     sch.Arguments,
					"error":         err,
				}, "An error occurred while Scheduler when call taskSrvClient.CallTask.")
				return
			}
			log.InfoWithFields(log.Fields{
				"task_codename":  sch.TaskCodename,
				"expression":     sch.Expression,
				"timeout":        sch.Timeout,
				"arguments":      sch.Arguments,
				"task_unique_id": rsp.TaskUniqueId,
			}, "Call task success.")
		})
		// 创建计划任务失败
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"task_codename": sch.TaskCodename,
				"expression":    sch.Expression,
				"arguments":     sch.Arguments,
				"error":         err,
			}, "An error occurred while Scheduler when call cron.AddFunc.")
			panic(fmt.Errorf("failed to add func, error: %s", err.Error()))
		}
		log.InfoWithFields(log.Fields{
			"task_codename": sch.TaskCodename,
			"expression":    sch.Expression,
			"timeout":       sch.Timeout,
			"arguments":     sch.Arguments,
		}, "Scheduler task added.")
	}

	s.started = true
	s.register()
	s.cron.Start()
}

// Stop 停止计划任务
func (s *Scheduler) Stop() {
	s.started = false
	if s.cron != nil {
		s.cron.Stop()
	}
	s.unregister()
	_ = s.etcdCli.Close()
}

// getScheduleTask 获得已配置的计划任务
func (s *Scheduler) getScheduleTask() []*Schedule {
	log.Info("Scheduler getScheduleTask called.")
	defer log.Info("Scheduler getScheduleTask end.")

	var (
		maxPg uint32 = 2
		ctx          = context.Background()
	)

	res := make([]*Schedule, 0)
	for pg := uint32(1); pg < maxPg; pg++ {
		req := &proto.QueryWithPage{
			Page:     pg,
			PageSize: CLIENT_MAX_PAGESIZE,
		}
		rsp, err := s.taskSrvClient.ListSchedules(ctx, req)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"page":  pg,
				"error": err,
			}, "An error occurred while taskSrvClient.ListSchedules.")
			panic(fmt.Errorf("failed to list schedule tasks, error: %s", err.Error()))
		}
		maxPg = rsp.Pages
		for _, r := range rsp.Schedules {
			log.DebugWithFields(log.Fields{
				"task_codename": r.TaskCodename,
				"expression":    r.Expression,
				"timeout":       r.Timeout,
				"arguments":     r.Arguments,
				"disabled":      r.Disabled,
			}, "Got a schedule.")
			// 跳过禁用的计划任务
			if r.Disabled {
				log.Debug("Skip disabled schedule.")
				continue
			}
			res = append(res, &Schedule{
				TaskCodename: r.TaskCodename,
				Expression:   r.Expression,
				Timeout:      r.Timeout,
				Arguments:    r.Arguments,
			})
		}
	}
	return res
}

// register
func (s *Scheduler) register() {
	log.Info("Scheduler register called.")
	defer log.Info("Scheduler register end.")

	ctx := context.TODO()

	// 建立etcd租约
	s.etcdLease = clientv3.NewLease(s.etcdCli)
	leaseGrantResp, err := s.etcdLease.Grant(ctx, s.opts.RegisterTtl)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while lease.Grant.")
		panic(err)
	}

	ch, err := s.etcdLease.KeepAlive(ctx, leaseGrantResp.ID)
	// 续约应答
	go func() {
		for {
			_, ok := <-ch
			if !ok {
				log.Info("Scheduler regCh may closed.")
				break
			}
		}
	}()

	// 生成注册Value
	regV, _ := json.Marshal(ScheduleInfo{
		IpAddress: ipv4.LocalIP(),
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	txn := clientv3.NewKV(s.etcdCli).Txn(ctx)
	// 注册到etcd
	txn.If(clientv3.Compare(clientv3.CreateRevision(SCHEDULER_REGISTER_KEY), "=", 0)).
		Then(clientv3.OpPut(SCHEDULER_REGISTER_KEY, string(regV), clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(SCHEDULER_REGISTER_KEY))

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		panic(err)
	}

	if !txnResp.Succeeded {
		panic("another scheduler instance is already running")
	}
}

// unregister
func (s *Scheduler) unregister() {
	log.Info("Scheduler unregister called.")
	defer log.Info("Scheduler unregister end.")

	if s.etcdLease != nil {
		_ = s.etcdLease.Close()
	}

	_, err := s.etcdCli.Delete(context.Background(), SCHEDULER_REGISTER_KEY)
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"error": err,
		}, "An error occurred while etcdCli.Delete.")
	}
}

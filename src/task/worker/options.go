package worker

import (
	"eago/common/logger"
	"github.com/go-basic/ipv4"
)

const WorkerRegisterKeyPrefix = "/td/eago/workers"

// Options struct
type Options struct {
	EtcdAddresses      []string
	EtcdUsername       string
	EtcdPassword       string
	TaskRpcServiceName string

	WorkerIp    string
	ServiceName string

	MultiInstance bool

	RegisterTtl int64

	PrintResultLog      bool
	ResultLogBufferSize uint

	Logger *logger.Logger
}

// newOptions
func newOptions(opts ...Option) Options {
	opt := Options{
		EtcdAddresses:      []string{"127.0.0.1:2379", "127.0.0.1:2379", "127.0.0.1:2379"},
		EtcdUsername:       "",
		EtcdPassword:       "",
		TaskRpcServiceName: "eago.srv.task",

		WorkerIp:    ipv4.LocalIP(),
		ServiceName: "task",

		MultiInstance: true,

		RegisterTtl: defaultWorkerRegisterTtl,

		PrintResultLog:      false,
		ResultLogBufferSize: defaultWorkerResultLogBufferSize,
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Logger == nil {
		opt.Logger = logger.NewDefaultLogger()
	}

	return opt
}

type Option func(o *Options)

// EtcdAddresses 设置注册中心地址
func EtcdAddresses(addr []string) Option {
	return func(o *Options) {
		o.EtcdAddresses = addr
	}
}

// EtcdUsername 设置注册中心用户名
func EtcdUsername(uName string) Option {
	return func(o *Options) {
		o.EtcdUsername = uName
	}
}

// EtcdPassword 设置注册中心密码
func EtcdPassword(p string) Option {
	return func(o *Options) {
		o.EtcdPassword = p
	}
}

// TaskRpcServiceName 任务模块RPC地址
func TaskRpcServiceName(srvName string) Option {
	return func(o *Options) {
		o.TaskRpcServiceName = srvName
	}
}

// Deprecated: Modular 设置Worker服务名，旧方法
func Modular(m string) Option {
	return func(o *Options) {
		o.ServiceName = m
	}
}

// ServiceName 设置Worker服务名
func ServiceName(s string) Option {
	return func(o *Options) {
		o.ServiceName = s
	}
}

// MultiInstance 设置是否允许运行多个Worker实例
func MultiInstance(b bool) Option {
	return func(o *Options) {
		o.MultiInstance = b
	}
}

// RegisterTtl 注册超时时间
func RegisterTtl(ttl int64) Option {
	return func(o *Options) {
		o.RegisterTtl = ttl
	}
}

// PrintResultLog 设置Worker在执行中是否在本地TTY打印ResultLog
func PrintResultLog(in bool) Option {
	return func(o *Options) {
		o.PrintResultLog = in
	}
}

// ResultLogBufferSize 设置Worker的ResultLogBufferSize
func ResultLogBufferSize(size uint) Option {
	return func(o *Options) {
		o.ResultLogBufferSize = size
	}
}

// Logger 设置Logger
func Logger(in *logger.Logger) Option {
	return func(o *Options) {
		o.Logger = in
	}
}

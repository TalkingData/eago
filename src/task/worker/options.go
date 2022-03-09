package worker

import (
	"github.com/go-basic/ipv4"
)

const WORKER_REGISTER_KEY_PREFFIX = "/td/eago/workers"

// Options struct
type Options struct {
	EtcdAddresses      []string
	EtcdUsername       string
	EtcdPassword       string
	TaskRpcServiceName string

	WorkerIp    string
	ServiceName string

	MultiInstance bool

	RegisterTtl   int64
	LogBufferSize uint
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

		RegisterTtl:   WORKER_REGISTER_TTL,
		LogBufferSize: WORKER_LOGGER_BUFFER_SIZE,
	}

	for _, o := range opts {
		o(&opt)
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

// LogBufferSize 设置Worker的LogBufferSize
func LogBufferSize(size uint) Option {
	return func(o *Options) {
		o.LogBufferSize = size
	}
}

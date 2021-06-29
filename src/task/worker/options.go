package worker

import (
	"github.com/go-basic/ipv4"
)

const WORKER_REGISTER_KEY_PREFFIX = "/td/eago/workers"

type Options struct {
	EtcdAddresses      []string
	EtcdUsername       string
	EtcdPassword       string
	TaskRpcServiceName string

	WorkerIp string
	Modular  string

	RegisterTtl   int64
	LogBufferSize uint
}

func newOptions(opts ...Option) Options {
	opt := Options{
		EtcdAddresses:      []string{"127.0.0.1:2379", "127.0.0.1:2379", "127.0.0.1:2379"},
		EtcdUsername:       "",
		EtcdPassword:       "",
		TaskRpcServiceName: "eago.srv.task",

		WorkerIp: ipv4.LocalIP(),
		Modular:  "task",

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

// Modular 设置Worker模块
func Modular(m string) Option {
	return func(o *Options) {
		o.Modular = m
	}
}

// LogBufferSize 设置Worker的LogBufferSize
func LogBufferSize(size uint) Option {
	return func(o *Options) {
		o.LogBufferSize = size
	}
}

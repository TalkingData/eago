package main

// Options struct
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

// newOptions
func newOptions(opts ...Option) Options {
	opt := Options{
		EtcdAddresses: []string{"127.0.0.1:2379", "127.0.0.1:2379", "127.0.0.1:2379"},
		EtcdUsername:  "",
		EtcdPassword:  "",
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

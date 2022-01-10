package main

// Options struct
type Options struct {
	EtcdAddresses []string
	EtcdUsername  string
	EtcdPassword  string

	TaskRpcRegisterKey string
	TaskRpcRetries     int

	RegisterTtl int64
}

// newOptions
func newOptions(opts ...Option) Options {
	opt := Options{
		EtcdAddresses: []string{"127.0.0.1:2379", "127.0.0.1:2379", "127.0.0.1:2379"},
		EtcdUsername:  "",
		EtcdPassword:  "",

		TaskRpcRegisterKey: "",
		TaskRpcRetries:     0,

		RegisterTtl: 10,
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

// TaskRpcRegisterKey 设置任务RPC注册Key名
func TaskRpcRegisterKey(k string) Option {
	return func(o *Options) {
		o.TaskRpcRegisterKey = k
	}
}

// TaskRpcRetries 设置任务RPC重试次数
func TaskRpcRetries(retries int) Option {
	return func(o *Options) {
		o.TaskRpcRetries = retries
	}
}

// RegisterTtl 设置Etcd注册TTL
func RegisterTtl(ttl int64) Option {
	return func(o *Options) {
		o.RegisterTtl = ttl
	}
}

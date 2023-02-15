package tracer

import "eago/common/global"

const (
	defaultRegisterKey   = "undefined.undefined"
	defaultJaegerAddress = "127.0.0.1:5775"
	defaultSamplerRate   = 1.0
)

type Option func(o *Options)

// Option struct
type Options struct {
	RegisterKey string

	JaegerHostPort string
	SamplerRate    float64

	CtxTracerKey string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		RegisterKey: defaultRegisterKey,

		JaegerHostPort: defaultJaegerAddress,
		SamplerRate:    defaultSamplerRate,

		CtxTracerKey: global.OpentracingCtxKey,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// RegisterKey 设置RegisterKey
func RegisterKey(in string) Option {
	return func(o *Options) {
		o.RegisterKey = in
	}
}

// JaegerHostPort 设置JaegerHostPort
func JaegerHostPort(in string) Option {
	return func(o *Options) {
		o.JaegerHostPort = in
	}
}

// SamplerRate 设置SamplerRate
func SamplerRate(in float64) Option {
	return func(o *Options) {
		o.SamplerRate = in
	}
}

// CtxTracerKey 设置CtxTracerKey，不建议修改
func CtxTracerKey(in string) Option {
	return func(o *Options) {
		o.CtxTracerKey = in
	}
}

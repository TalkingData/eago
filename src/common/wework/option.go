package wework

import (
	"eago/common/logger"
	"time"
)

const (
	defaultHttpTimeoutSecs            = 5
	defaultTokenExpirationAdvanceSecs = 60

	defaultTimestampFormat = "2006-01-02 15:04:05"
)

type Option func(o *Options)

// Option struct
type Options struct {
	HttpTimeoutSecs               time.Duration
	TokenExpirationAdvanceSeconds float64

	TimestampFormat string

	Logger *logger.Logger
}

func newOptions(opts ...Option) Options {
	opt := Options{
		HttpTimeoutSecs:               defaultHttpTimeoutSecs,
		TokenExpirationAdvanceSeconds: defaultTokenExpirationAdvanceSecs,

		TimestampFormat: defaultTimestampFormat,
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Logger == nil {
		opt.Logger = logger.NewDefaultLogger()
	}

	return opt
}

// HttpTimeoutSecs 设置HttpTimeoutSecs，不建议修改
func HttpTimeoutSecs(in time.Duration) Option {
	return func(o *Options) {
		o.HttpTimeoutSecs = in
	}
}

// TokenExpirationAdvanceSeconds 设置TokenExpirationAdvanceSeconds，不建议修改
func TokenExpirationAdvanceSeconds(in float64) Option {
	return func(o *Options) {
		o.TokenExpirationAdvanceSeconds = in
	}
}

// TimestampFormat 设置TimestampFormat，不建议修改
func TimestampFormat(in string) Option {
	return func(o *Options) {
		o.TimestampFormat = in
	}
}

// Logger 设置Logger
func Logger(in *logger.Logger) Option {
	return func(o *Options) {
		o.Logger = in
	}
}

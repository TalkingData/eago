package broker

import (
	"eago/common/logger"
)

const (
	defaultServiceName     = "undefined"
	defaultTopicSeparator  = "topic"
	defaultTimestampFormat = "2006-01-02 15:04:05"
)

type Option func(o *Options)

// Option struct
type Options struct {
	ServiceName    string
	TopicSeparator string

	TimestampFormat string

	Logger *logger.Logger
}

func newOptions(opts ...Option) Options {
	opt := Options{
		TopicSeparator: defaultTopicSeparator,
		ServiceName:    defaultServiceName,

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

// ServiceName 设置ServiceName
func ServiceName(in string) Option {
	return func(o *Options) {
		o.ServiceName = in
	}
}

// TopicSeparator 设置TopicSeparator，不建议修改
func TopicSeparator(in string) Option {
	return func(o *Options) {
		o.TopicSeparator = in
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

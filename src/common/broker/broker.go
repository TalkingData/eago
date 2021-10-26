package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
)

var brk broker.Broker

// NewBroker 生成一个Broker
func NewBroker(addresses []string) broker.Broker {
	brk = kafka.NewBroker(broker.Addrs(addresses...))

	if err := brk.Connect(); err != nil {
		panic(err)
	}

	return brk
}

// Close 关闭Broker
func Close() {
	if brk == nil {
		return
	}
	_ = brk.Disconnect()
}

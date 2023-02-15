package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
)

// NewKafkaBroker 生成一个KafkaBroker
func NewKafkaBroker(addresses []string) (broker.Broker, error) {
	brk := kafka.NewBroker(broker.Addrs(addresses...))
	return brk, brk.Connect()
}

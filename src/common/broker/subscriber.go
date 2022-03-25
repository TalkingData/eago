package broker

import (
	"context"
	"eago/common/log"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"github.com/mitchellh/mapstructure"
)

type HandleFunc func(ctx context.Context, m *Message) error

type Subscriber interface {
	// Subscribe 订阅
	Subscribe(service, model, event string, fn HandleFunc)
}

// subscriber struct
type subscriber struct {
	Service  string
	consumer string

	brk broker.Broker
}

// NewSubscriber 创建Subscriber
func NewSubscriber(service string) Subscriber {
	if brk == nil {
		panic("You should create a broker by NewBroker first.")
	}

	return &subscriber{
		Service:  service,
		consumer: fmt.Sprintf("%s.consumers", service),
		brk:      brk,
	}
}

// Subscribe 订阅
func (s *subscriber) Subscribe(service, model, event string, fn HandleFunc) {
	log.Info("subscriber.Subscribe called.")
	defer log.Info("subscriber.Subscribe end.")

	fullTopic := fmt.Sprintf("%s.%s.%s.%s", service, _TOPIC_SEPARATOR, model, event)

	foo := func(e broker.Event) error {
		log.InfoWithFields(log.Fields{
			"topic":    fullTopic,
			"consumer": s.consumer,
		}, "Got message from broker.")

		msg := &Message{}

		// 序列化map到Message结构体
		if err := mapstructure.Decode(e.Message().Header, &msg); err != nil {
			log.ErrorWithFields(log.Fields{
				"topic":    fullTopic,
				"consumer": s.consumer,
				"error":    err,
			}, "Error in mapstructure.Decode e.Message().Header, Ignore this message.")
		}

		bd := make(map[string]interface{})
		err := json.Unmarshal(e.Message().Body, &bd)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"topic":         fullTopic,
				"consumer":      s.consumer,
				"message_uuid":  msg.Uuid,
				"message_from":  msg.From,
				"message_event": msg.Event,
				"error":         err,
			}, "Error in json.Unmarshal event.Message().Body, Ignore this message.")
			return nil
		}

		msg.Body = bd

		err = fn(context.Background(), msg)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"topic":    fullTopic,
				"consumer": s.consumer,
				"error":    err,
			}, "Error in HandleFunc.")
			return err
		}

		return nil
	}

	if _, err := s.brk.Subscribe(fullTopic, foo, broker.Queue(s.consumer)); err != nil {
		log.ErrorWithFields(log.Fields{
			"topic":    fullTopic,
			"consumer": s.consumer,
			"error":    err,
		}, "Error in s.opts.Broker.Subscribe.")
	}
}

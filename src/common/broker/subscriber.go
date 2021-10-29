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
	RegAndSub(topic string, fn HandleFunc)
}

type subscriber struct {
	Modular  string
	consumer string

	brk broker.Broker
}

// NewSubscriber 创建Subscriber
func NewSubscriber(modular string) Subscriber {
	if brk == nil {
		panic("You should create a broker by NewBroker first.")
	}

	s := &subscriber{
		Modular: modular,
		brk:     brk,
	}
	s.consumer = fmt.Sprintf("%s.consumers", s.Modular)
	return s
}

func (s *subscriber) RegAndSub(topic string, fn HandleFunc) {
	log.Info("subscriber.RegAndSub called.")
	defer log.Info("subscriber.RegAndSub end.")

	foo := func(e broker.Event) error {
		log.InfoWithFields(log.Fields{
			"topic":    topic,
			"consumer": s.consumer,
		}, "Got message from broker.")

		msg := &Message{}

		// 序列化map到Message结构体
		if err := mapstructure.Decode(e.Message().Header, &msg); err != nil {
			log.ErrorWithFields(log.Fields{
				"topic":    topic,
				"consumer": s.consumer,
			}, "Error in mapstructure.Decode e.Message().Header, Ignore this message.")
		}

		bd := make(map[string]interface{})
		err := json.Unmarshal(e.Message().Body, &bd)
		if err != nil {
			log.ErrorWithFields(log.Fields{
				"topic":         topic,
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
				"topic":    topic,
				"consumer": s.consumer,
				"error":    err,
			}, "Error in HandleFunc.")
			return err
		}

		return nil
	}

	_, err := s.brk.Subscribe(topic, foo, broker.Queue(s.consumer))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"topic":    topic,
			"consumer": s.consumer,
			"error":    err,
		}, "Error in s.opts.Broker.Subscribe.")
	}
}
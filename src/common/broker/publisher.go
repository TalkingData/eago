package broker

import (
	"context"
	"eago/common/log"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	uuid "github.com/satori/go.uuid"
)

type Publisher interface {
	Publish(ctx context.Context, model, event string, body map[string]interface{}) error
}

type publisher struct {
	Modular     string
	topicPrefix string

	brk broker.Broker
}

// NewPublisher 创建Publisher
func NewPublisher(modular string) Publisher {
	if brk == nil {
		panic("You should create a broker by NewBroker first.")
	}

	p := &publisher{
		Modular: modular,
		brk:     brk,
	}
	p.topicPrefix = fmt.Sprintf("%s.topic", p.Modular)
	return p
}

// Publish 发布消息
func (p *publisher) Publish(ctx context.Context, modelName, event string, body map[string]interface{}) error {
	log.Info("publisher.Publish called.")
	defer log.Info("publisher.Publish end.")

	topic := fmt.Sprintf("%s.%s.%s", p.topicPrefix, modelName, event)

	msg := Message{
		Uuid:  p.genMessageUuid(topic),
		From:  fmt.Sprintf("%s.%s", p.topicPrefix, modelName),
		Event: event,
		Body:  body,
	}

	// 发送消息
	log.Info("Prepare call Broker.Publish.")
	err := p.brk.Publish(topic, msg.ToBrokerMessage(), broker.PublishContext(ctx))
	if err != nil {
		log.InfoWithFields(log.Fields{
			"message_uuid":  msg.Uuid,
			"message_from":  msg.From,
			"message_event": msg.Event,
			"topic":         topic,
		}, "Error in Broker.Publish.")
		return err
	}

	return nil
}

// genMessageUuid 生成消息UUID
func (p *publisher) genMessageUuid(topic string) string {
	return uuid.NewV5(uuid.NewV1(), topic).String()
}

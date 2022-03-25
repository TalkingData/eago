package broker

import (
	"context"
	"eago/common/log"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	uuid "github.com/satori/go.uuid"
)

type Publisher interface {
	// Publish 发布消息
	Publish(ctx context.Context, model, event string, body map[string]interface{}) error
	// SpecificPublish 以指定服务发布消息
	SpecificPublish(ctx context.Context, service, modelName, event string, body map[string]interface{}) error
}

type publisher struct {
	Service string
	brk     broker.Broker
}

// NewPublisher 创建Publisher
func NewPublisher(service string) Publisher {
	if brk == nil {
		panic("You should create a broker by NewBroker first.")
	}

	return &publisher{
		Service: service,
		brk:     brk,
	}
}

// Publish 发布消息
func (p *publisher) Publish(ctx context.Context, modelName, event string, body map[string]interface{}) error {
	log.Info("publisher.Publish called.")
	defer log.Info("publisher.Publish end.")

	return p.pub(ctx, p.Service, modelName, event, body)
}

// SpecificPublish 以指定服务发布消息
func (p *publisher) SpecificPublish(ctx context.Context, service, modelName, event string, body map[string]interface{}) error {
	log.Info("publisher.SpecificPublish called.")
	defer log.Info("publisher.SpecificPublish end.")

	return p.pub(ctx, service, modelName, event, body)
}

// pub 具体发布消息的操作方法
func (p *publisher) pub(ctx context.Context, service, modelName, event string, body map[string]interface{}) error {
	topic := fmt.Sprintf("%s.%s.%s.%s", service, _TOPIC_SEPARATOR, modelName, event)

	msg := Message{
		Uuid:  p.genMessageUuid(topic),
		From:  fmt.Sprintf("%s.%s.%s", service, _TOPIC_SEPARATOR, modelName),
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

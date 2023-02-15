package broker

import (
	"context"
	"eago/common/logger"
	"github.com/micro/go-micro/v2/broker"
	uuid "github.com/satori/go.uuid"
)

type Publisher interface {
	Publish(ctx context.Context, currModel, tgtServiceName, tgtModel, event string, body map[string]interface{}) error
	GetServiceName() string
}

type publisher struct {
	ServiceName string

	broker broker.Broker
	logger *logger.Logger

	opts Options
}

// NewPublisher 创建Publisher
func NewPublisher(broker broker.Broker, options ...Option) Publisher {
	opts := newOptions(options...)

	return &publisher{
		ServiceName: opts.ServiceName,

		broker: broker,
		logger: opts.Logger,

		opts: opts,
	}
}

// Publish 发布消息
func (p *publisher) Publish(
	ctx context.Context, fromModel, tgtServiceName, tgtModel, event string, body map[string]interface{},
) error {
	p.logger.DebugWithFields(logger.Fields{
		"from_model":          fromModel,
		"target_service_name": tgtServiceName,
		"target_model":        tgtModel,
		"target_event":        event,
	}, "publisher.Publish called.")
	defer p.logger.Debug("publisher.Publish end.")

	tgtTopicName := genFullTopicName(tgtServiceName, p.opts.TopicSeparator, tgtModel, event)

	msg := Message{
		Uuid:  p.genMessageUuid(tgtTopicName),
		From:  genFromName(p.ServiceName, fromModel),
		Event: event,
		Body:  body,
	}

	// 发送消息
	p.logger.Debug("Prepare call publisher.broker.Publish.")
	err := p.broker.Publish(tgtTopicName, msg.ToBrokerMessage(p.opts.TimestampFormat), broker.PublishContext(ctx))
	if err != nil {
		p.logger.InfoWithFields(logger.Fields{
			"message_uuid":  msg.Uuid,
			"message_from":  msg.From,
			"message_event": msg.Event,
			"topic":         tgtTopicName,
		}, "An error occurred while publisher.broker.Publish.")
		return err
	}

	return nil
}

func (p *publisher) GetServiceName() string {
	return p.ServiceName
}

// genMessageUuid 生成消息UUID
func (p *publisher) genMessageUuid(topic string) string {
	return uuid.NewV5(uuid.NewV1(), topic).String()
}

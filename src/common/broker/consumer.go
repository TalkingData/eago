package broker

import (
	"context"
	"eago/common/logger"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"github.com/mitchellh/mapstructure"
)

type SubHandleFunc func(ctx context.Context, m *Message) error

type consumer struct {
	serviceName  string
	consumerName string

	broker broker.Broker

	logger *logger.Logger

	opts Options

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewConsumer(ctx context.Context, broker broker.Broker, options ...Option) *consumer {
	opts := newOptions(options...)

	consumerCtx, consumerCancel := context.WithCancel(ctx)

	return &consumer{
		serviceName: opts.ServiceName,

		consumerName: fmt.Sprintf("%s.consumers", opts.ServiceName),

		broker: broker,
		logger: opts.Logger,

		opts: opts,

		ctx:        consumerCtx,
		cancelFunc: consumerCancel,
	}
}

func (c *consumer) Start() error {
	c.logger.Info(fmt.Sprintf("Starting %s consumer...", c.serviceName))

	for {
		select {
		case <-c.ctx.Done():
			c.logger.InfoWithFields(logger.Fields{
				"context_error": c.ctx.Err(),
			}, "Consumer stopped by context done.")
			return c.ctx.Err()
		}
	}
}

func (c *consumer) Stop() {
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
}

func (c *consumer) Subscribe(model, event string, h SubHandleFunc) {
	c.logger.Debug("consumer.Subscribe called.")
	defer c.logger.Debug("consumer.Subscribe end.")

	subTopicName := genFullTopicName(c.serviceName, c.opts.TopicSeparator, model, event)
	foo := func(e broker.Event) error {
		c.logger.InfoWithFields(logger.Fields{
			"topic":         subTopicName,
			"consumer_name": c.consumerName,
		}, "Got message from broker.")

		msg := &Message{}

		// 序列化map到Message结构体
		if err := mapstructure.Decode(e.Message().Header, &msg); err != nil {
			c.logger.WarnWithFields(logger.Fields{
				"topic":          subTopicName,
				"consumer_name":  c.consumerName,
				"message_header": e.Message().Header,
				"error":          err,
			}, "An error occurred while mapstructure.Decode e.Message().Header, This message will be discarded.")
		}

		bd := make(map[string]interface{})
		if err := json.Unmarshal(e.Message().Body, &bd); err != nil {
			c.logger.WarnWithFields(logger.Fields{
				"topic":         subTopicName,
				"consumer_name": c.consumerName,
				"message_uuid":  msg.Uuid,
				"message_from":  msg.From,
				"message_event": msg.Event,
				"message_body":  e.Message().Body,
				"error":         err,
			}, "An error occurred while run json.Unmarshal event.Message().Body, This message will be discarded.")
			return nil
		}

		msg.Body = bd

		if err := h(c.ctx, msg); err != nil {
			c.logger.ErrorWithFields(logger.Fields{
				"topic":         subTopicName,
				"consumer_name": c.consumerName,
				"error":         err,
			}, "An error occurred while run SubHandleFunc.")
			return err
		}

		return nil
	}

	if _, err := c.broker.Subscribe(subTopicName, foo, broker.Queue(c.consumerName)); err != nil {
		c.logger.ErrorWithFields(logger.Fields{
			"topic":         subTopicName,
			"consumer_name": c.consumerName,
			"error":         err,
		}, "An error occurred while consumer.broker.Subscribe.")
	}
}

package broker

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"time"
)

type Message struct {
	Uuid     string                 `json:"uuid"`
	From     string                 `json:"from"`
	Event    string                 `json:"event"`
	Datetime string                 `json:"datetime"`
	Body     map[string]interface{} `json:"body"`
}

// ToBrokerMessage 转换为broker.Message
func (m *Message) ToBrokerMessage(tsFmt string) *broker.Message {
	h := map[string]string{
		"uuid":     m.Uuid,
		"from":     m.From,
		"event":    m.Event,
		"datetime": time.Now().Format(tsFmt),
	}

	bd, _ := json.Marshal(m.Body)

	return &broker.Message{Header: h, Body: bd}
}

// genFullTopicName 生成Topic名
func genFullTopicName(serviceName, topicSeparator, model, event string) string {
	return fmt.Sprintf("%s.%s.%s.%s", serviceName, topicSeparator, model, event)
}

// genFromName 生成消息From
func genFromName(serviceName, model string) string {
	return fmt.Sprintf("%s.%s", serviceName, model)
}

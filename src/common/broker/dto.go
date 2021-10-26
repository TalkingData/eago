package broker

import (
	"encoding/json"
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
func (m *Message) ToBrokerMessage() *broker.Message {
	h := map[string]string{
		"uuid":     m.Uuid,
		"from":     m.From,
		"event":    m.Event,
		"datetime": time.Now().Format("2006-01-02 15:04:05"),
	}

	bd, _ := json.Marshal(m.Body)

	return &broker.Message{Header: h, Body: bd}
}

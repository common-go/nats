package nats

import (
	"context"
	"net/http"

	"github.com/nats-io/nats.go"
)

type Producer struct {
	Conn *nats.Conn
	Subject string
}

func NewProducer(conn *nats.Conn, subject string) *Producer {
	return &Producer{conn, subject}
}

func (p *Producer) Produce(ctx context.Context, data []byte, messageAttributes *map[string]string) (string, error) {
	if messageAttributes == nil {
		err := p.Conn.Publish(p.Subject, data)
		return "", err
	} else {
		header := MapToHeader(messageAttributes)
		var msg = &nats.Msg{
			Subject: p.Subject,
			Data: data,
			Reply: "",
			Header: *header,
		}
		err := p.Conn.PublishMsg(msg)
		return "", err
	}
}

func MapToHeader(messageAttributes *map[string]string) *http.Header {
	if messageAttributes == nil || len (*messageAttributes) == 0 {
		return nil
	}
	header := &http.Header{}
	for k, v := range *messageAttributes {
		header.Add(k, v)
	}
	return header
}
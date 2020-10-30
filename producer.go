package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"net/http"
)

type Producer struct {
	Conn    *nats.Conn
	Subject string
}

func NewProducer(conn *nats.Conn, subject string) *Producer {
	return &Producer{conn, subject}
}
func NewProducerByConfig(p ProducerConfig) (*Producer, error) {
	if p.Connection.Retry.Retry1 <= 0 {
		conn, err := nats.Connect(p.Connection.Url, p.Connection.Options)
		if err != nil {
			return nil, err
		}
		return NewProducer(conn, p.Subject), nil
	} else {
		durations := DurationsFromValue(p.Connection.Retry, "Retry", 9)
		conn, err := NewConn(durations, p.Connection.Url, p.Connection.Options)
		if err != nil {
			return nil, err
		}
		return NewProducer(conn, p.Subject), nil
	}
}
func (p *Producer) Produce(ctx context.Context, data []byte, messageAttributes *map[string]string) (string, error) {
	defer p.Conn.Flush()
	if messageAttributes == nil {
		err := p.Conn.Publish(p.Subject, data)
		return "", err
	} else {
		header := MapToHeader(messageAttributes)
		var msg = &nats.Msg{
			Subject: p.Subject,
			Data:    data,
			Reply:   "",
			Header:  *header,
		}
		err := p.Conn.PublishMsg(msg)
		return "", err
	}
}

func MapToHeader(messageAttributes *map[string]string) *http.Header {
	if messageAttributes == nil || len(*messageAttributes) == 0 {
		return nil
	}
	header := &http.Header{}
	for k, v := range *messageAttributes {
		header.Add(k, v)
	}
	return header
}

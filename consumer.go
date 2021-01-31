package nats

import (
	"context"
	"github.com/common-go/mq"
	"github.com/nats-io/nats.go"
	"net/http"
	"runtime"
)

type Consumer struct {
	Conn    *nats.Conn
	Subject string
	Header  bool
}

func NewConsumer(conn *nats.Conn, subject string, header bool) *Consumer {
	return &Consumer{conn, subject, header}
}

func NewConsumerByConfig(c ConsumerConfig) (*Consumer, error) {
	if c.Connection.Retry.Retry1 <= 0 {
		conn, err := nats.Connect(c.Connection.Url, c.Connection.Options)
		if err != nil {
			return nil, err
		}
		return NewConsumer(conn, c.Subject, c.Header), nil
	} else {
		durations := DurationsFromValue(c.Connection.Retry, "Retry", 9)
		conn, err := NewConn(durations, c.Connection.Url, c.Connection.Options)
		if err != nil {
			return nil, err
		}
		return NewConsumer(conn, c.Subject, c.Header), nil
	}
}

func (c *Consumer) Consume(ctx context.Context, handle func(context.Context, *mq.Message, error) error) {
	if c.Header {
		c.Conn.Subscribe(c.Subject, func(msg *nats.Msg) {
			attrs := HeaderToMap(msg.Header)
			message := &mq.Message{
				Data:       msg.Data,
				Attributes: attrs,
				Raw:        msg,
			}
			handle(ctx, message, nil)
		})
		c.Conn.Flush()
		runtime.Goexit()
	} else {
		c.Conn.Subscribe(c.Subject, func(msg *nats.Msg) {
			message := &mq.Message{
				Data: msg.Data,
				Raw:  msg,
			}
			handle(ctx, message, nil)
		})
		c.Conn.Flush()
		runtime.Goexit()
	}
}

func HeaderToMap(header http.Header) map[string]string {
	attributes := make(map[string]string, 0)
	for name, values := range header {
		for _, value := range values {
			attributes[name] = value
		}
	}
	return attributes
}

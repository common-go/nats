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

func (c *Consumer) Consume(ctx context.Context, caller mq.ConsumerCaller) {
	if c.Header {
		c.Conn.Subscribe(c.Subject, func(msg *nats.Msg) {
			attrs := HeaderToMap(msg.Header)
			message := &mq.Message{
				Data:       msg.Data,
				Attributes: attrs,
				Raw:        msg,
			}
			caller.Call(ctx, message, nil)
		})
		c.Conn.Flush()
		runtime.Goexit()
	} else {
		c.Conn.Subscribe(c.Subject, func(msg *nats.Msg) {
			message := &mq.Message{
				Data: msg.Data,
				Raw:  msg,
			}
			caller.Call(ctx, message, nil)
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

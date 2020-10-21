package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"net"
	"time"
)

type HttpHealthCheck struct {
	name    string
	url     string
	timeout time.Duration
}

func NewHttpHealthCheck(name, url string, timeout time.Duration) *HttpHealthCheck {
	return &HttpHealthCheck{name, url, timeout}
}

func NewDefaultHttpHealthCheck(name, url string) *HttpHealthCheck {
	return &HttpHealthCheck{name, url, 5 * time.Second}
}

func (s *HttpHealthCheck) Name() string {
	return s.name
}

func (s *HttpHealthCheck) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	dialer := &net.Dialer{Timeout: s.timeout,DualStack: true}
	opts := &nats.Options{
		Servers: []string{s.url},
		Dialer:  dialer,
	}
	conn, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	conn.Close()
	res["status"] = "success"
	return res, nil
}

func (s *HttpHealthCheck) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	data["error"] = err.Error()
	return data
}

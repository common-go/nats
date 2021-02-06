package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"net"
	"time"
)

type HttpHealthChecker struct {
	name    string
	url     string
	timeout time.Duration
}

func NewHealthChecker(name, url string, timeout time.Duration) *HttpHealthChecker {
	return &HttpHealthChecker{name: name, url: url, timeout: timeout}
}

func NewHttpHealthChecker(name, url string) *HttpHealthChecker {
	return &HttpHealthChecker{name: name, url: url, timeout: 4 * time.Second}
}

func (s *HttpHealthChecker) Name() string {
	return s.name
}

func (s *HttpHealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	dialer := &net.Dialer{Timeout: s.timeout, DualStack: true}
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

func (s *HttpHealthChecker) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	data["error"] = err.Error()
	return data
}

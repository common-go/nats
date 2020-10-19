package nats

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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
	client := http.Client{
		Timeout: s.timeout,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(s.url)
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return res, fmt.Errorf("time out: %w", e)
	} else if err != nil {
		return res, err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return res, nil
	} else {
		return res, fmt.Errorf("status code is: %d", resp.StatusCode)
	}
}

func (s *HttpHealthCheck) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	data["error"] = err.Error()
	return data
}

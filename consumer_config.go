package nats

import "github.com/nats-io/nats.go"

type ConsumerConfig struct {
	Subject    string     `mapstructure:"subject"`
	Header     bool       `mapstructure:"header"`
	Connection ConnConfig `mapstructure:"connection"`
}

type ConnConfig struct {
	Url     string      `mapstructure:"url"`
	Options nats.Option `mapstructure:"option"`
}

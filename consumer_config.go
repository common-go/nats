package nats

type ConsumerConfig struct {
	Subject    string     `mapstructure:"subject"`
	Header     bool       `mapstructure:"header"`
	Connection ConnConfig `mapstructure:"connection"`
}

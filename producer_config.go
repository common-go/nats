package nats

type ProducerConfig struct {
	Subject    string     `mapstructure:"subject"`
	Connection ConnConfig `mapstructure:"connection"`
}

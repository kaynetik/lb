package client

import (
	"flag"
)

type Config struct {
	QueueName string
	InputFile string
}

func ParseConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.QueueName, "queue", "commandQueue", "Name of the RabbitMQ queue")
	flag.StringVar(&config.InputFile, "input", "", "File containing commands (optional)")
	flag.Parse()
	return config
}

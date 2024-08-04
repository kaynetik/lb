package client

import (
	"flag"
)

type Config struct {
	RabbitMQConfig
	ServiceConfig
}

type RabbitMQConfig struct {
	// DialTarget is the RabbitMQ dial target.
	DialTarget string
	// QueueName is the name of the RabbitMQ queue.
	QueueName string
	// InputFile is the file containing commands (optional).
	InputFile string
}

type ServiceConfig struct {
	Name            string
	Environment     string
	OTELDigesterURL string
}

func ParseConfig() *Config {
	config := &Config{}

	// Populate ServiceConfig
	flag.StringVar(&config.Name, "name", "client", "Name of the service")
	flag.StringVar(&config.Environment, "env", "dev", "Environment of the service")
	flag.StringVar(&config.OTELDigesterURL, "oteldigester", "otel-collector:4317", "URL of the OTEL digester") // TODO: Digester still not enabled.

	// Populate RabbitMQConfig
	flag.StringVar(&config.DialTarget, "dial", "amqp://guest:guest@rabbitmq:5672/", "RabbitMQ dial target")
	flag.StringVar(&config.QueueName, "queue", "commandQueue", "Name of the RabbitMQ queue")
	flag.StringVar(&config.InputFile, "input", "", "File containing commands (optional)")
	flag.Parse()

	return config
}

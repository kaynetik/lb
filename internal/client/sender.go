package client

import (
	"lightblocks/internal/server/queue"
)

type Queue interface {
	Publish(message string) error
	Consume(handler func(string)) error
	Close()
}

type Sender struct {
	queue Queue
}

func NewSender(queueName string) (*Sender, error) {
	rabbitMQ, err := queue.NewRabbitMQ(queueName)
	if err != nil {
		return nil, err
	}

	return &Sender{
		queue: rabbitMQ,
	}, nil
}

func (s *Sender) Send(message string) error {
	return s.queue.Publish(message)
}

func (s *Sender) Close() {
	s.queue.Close()
}

package client

import (
	"lightblocks/internal/server/queue"
)

type Sender struct {
	rabbitMQ *queue.RabbitMQ
}

func NewSender(queueName string) (*Sender, error) {
	rabbitMQ, err := queue.NewRabbitMQ(queueName)
	if err != nil {
		return nil, err
	}
	return &Sender{rabbitMQ: rabbitMQ}, nil
}

func (s *Sender) Send(command string) error {
	return s.rabbitMQ.Publish(command)
}

func (s *Sender) Close() {
	s.rabbitMQ.Close()
}

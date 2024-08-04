package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"lightblocks/internal/server/queue"
)

func TestSender(t *testing.T) {
	mockRabbitMQ := queue.NewMockRabbitMQ()

	sender := &Sender{
		queue: mockRabbitMQ,
	}
	defer sender.Close()

	messageReceived := make(chan string)
	err := mockRabbitMQ.Consume(func(message string) {
		messageReceived <- message
	})
	assert.NoError(t, err)

	err = sender.Send("test message")
	assert.NoError(t, err)

	select {
	case msg := <-messageReceived:
		assert.Equal(t, "test message", msg)
	case <-time.After(1 * time.Second):
		t.Fatal("Message was not received")
	}
}

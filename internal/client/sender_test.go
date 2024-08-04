package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"lightblocks/internal/server/queue"
)

func TestSender(t *testing.T) {
	rabbitMQ, err := queue.NewRabbitMQ("testQueue")
	assert.NoError(t, err)
	defer rabbitMQ.Close()

	sender, err := NewSender("testQueue")
	assert.NoError(t, err)
	defer sender.Close()

	messageReceived := make(chan string)
	err = rabbitMQ.Consume(func(message string) {
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

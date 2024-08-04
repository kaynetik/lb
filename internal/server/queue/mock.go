package queue

// MockRabbitMQ is a mock implementation of the RabbitMQ interface.
// TODO: Realistically depending on the needs `mockery` might be a better choice.
type MockRabbitMQ struct {
	messages chan string
}

func NewMockRabbitMQ() *MockRabbitMQ {
	return &MockRabbitMQ{
		messages: make(chan string, 100), // buffer to avoid blocking
	}
}

func (mq *MockRabbitMQ) Publish(message string) error {
	mq.messages <- message
	return nil
}

func (mq *MockRabbitMQ) Consume(handler func(string)) error {
	go func() {
		for msg := range mq.messages {
			handler(msg)
		}
	}()
	return nil
}

func (mq *MockRabbitMQ) Close() {
	close(mq.messages)
}

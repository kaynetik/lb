package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	os.Args = []string{"cmd", "-queue", "testQueue", "-input", "test.txt"}
	config := ParseConfig()
	assert.Equal(t, "testQueue", config.QueueName)
	assert.Equal(t, "test.txt", config.InputFile)
}

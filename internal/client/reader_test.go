package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCommandsFromStdin(t *testing.T) {
	input := "command1\ncommand2\ncommand3\n"
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	go func() {
		w.WriteString(input)
		w.Close()
	}()

	commands, err := readCommandsFromStdin()
	assert.NoError(t, err)
	assert.Equal(t, []string{"command1", "command2", "command3"}, commands)

	os.Stdin = oldStdin
}

func TestReadCommandsFromFile(t *testing.T) {
	input := "command1\ncommand2\ncommand3\n"
	file, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString(input)
	assert.NoError(t, err)
	file.Close()

	commands, err := readCommandsFromFile(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, []string{"command1", "command2", "command3"}, commands)
}

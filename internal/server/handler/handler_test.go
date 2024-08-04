package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"lightblocks/internal/observer"
	orderedmap "lightblocks/internal/server/map"
	"os"
	"testing"
)

func TestHandlers(t *testing.T) {
	om := orderedmap.NewOrderedMap()
	opChan := make(chan orderedmap.Operation)
	go om.Run(opChan)

	// TODO: Create a no-op observer for testing
	// I'm not certain I'll have time to properly cleanup this part.
	observer.InitObserver("server/handler/test", "", "test")
	obs, _ := observer.Action(context.Background(), tracer)

	AddItemHandler(opChan, "key1", "value1")
	AddItemHandler(opChan, "key2", "value2")
	DeleteItemHandler(opChan, "key1")

	// Check key1
	getResultChan := make(chan interface{})
	opChan <- orderedmap.Operation{
		Action: orderedmap.Get,
		Key:    "key1",
		Result: getResultChan,
	}
	result := (<-getResultChan).(struct {
		Value  string
		Exists bool
	})
	assert.False(t, result.Exists)
	assert.Equal(t, "", result.Value)

	// Check key2
	getResultChan = make(chan interface{})
	opChan <- orderedmap.Operation{
		Action: orderedmap.Get,
		Key:    "key2",
		Result: getResultChan,
	}
	result = (<-getResultChan).(struct {
		Value  string
		Exists bool
	})
	assert.True(t, result.Exists)
	assert.Equal(t, "value2", result.Value)

	// Clean up output file if it exists
	if err := os.Remove("output.txt"); err != nil {
		return
	}

	GetItemHandler(obs, opChan, "key2")
	GetAllItemsHandler(obs, opChan)

	content, err := os.ReadFile("output.txt")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "key: key2, value: value2")
}

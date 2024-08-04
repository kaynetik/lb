package handler

import (
	orderedmap "lightblocks/internal/server/map"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	om := orderedmap.NewOrderedMap()

	AddItemHandler(om, "key1", "value1")
	AddItemHandler(om, "key2", "value2")
	DeleteItemHandler(om, "key1")

	value, exists := om.Get("key1")
	assert.False(t, exists)
	assert.Equal(t, "", value)

	value, exists = om.Get("key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", value)

	os.Remove("output.txt")
	GetItemHandler(om, "key2")
	GetAllItemsHandler(om)

	content, _ := os.ReadFile("output.txt")
	assert.Contains(t, string(content), "key: key2, value: value2")
}

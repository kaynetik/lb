package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	om := NewOrderedMap()

	om.Add("key1", "value1")
	om.Add("key2", "value2")

	value, exists := om.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	om.Delete("key1")
	_, exists = om.Get("key1")
	assert.False(t, exists)

	om.Add("key3", "value3")
	allItems := om.GetAll()
	assert.Equal(t, 2, len(allItems))
	assert.Equal(t, "key2", allItems[0].Key)
	assert.Equal(t, "key3", allItems[1].Key)
}

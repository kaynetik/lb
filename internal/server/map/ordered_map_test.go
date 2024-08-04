package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	tests := []struct {
		name       string
		operations []Operation
		expect     interface{}
	}{
		{
			name: "Add and Get Item",
			operations: []Operation{
				{Action: Add, Key: "key1", Value: "value1"},
				{Action: Get, Key: "key1", Result: make(chan interface{})},
			},
			expect: struct {
				Value  string
				Exists bool
			}{"value1", true},
		},
		{
			name: "Add, Delete, and Get Item",
			operations: []Operation{
				{Action: Add, Key: "key2", Value: "value2"},
				{Action: Delete, Key: "key2"},
				{Action: Get, Key: "key2", Result: make(chan interface{})},
			},
			expect: struct {
				Value  string
				Exists bool
			}{"", false},
		},
		{
			name: "Get All Items",
			operations: []Operation{
				{Action: Add, Key: "key3", Value: "value3"},
				{Action: Add, Key: "key4", Value: "value4"},
				{Action: GetAll, Result: make(chan interface{})},
			},
			expect: []KeyValuePair{
				{"key3", "value3"},
				{"key4", "value4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := NewOrderedMap()
			opChan := make(chan Operation)

			go om.Run(opChan)

			var result interface{}
			for _, op := range tt.operations {
				if op.Result != nil {
					opChan <- op
					result = <-op.Result
				} else {
					opChan <- op
				}
			}

			switch expected := tt.expect.(type) {
			case struct {
				Value  string
				Exists bool
			}:
				actual := result.(struct {
					Value  string
					Exists bool
				})
				assert.Equal(t, expected.Exists, actual.Exists)
				assert.Equal(t, expected.Value, actual.Value)
			case []KeyValuePair:
				actual := result.([]KeyValuePair)
				assert.Equal(t, len(expected), len(actual))
				for i, kv := range expected {
					assert.Equal(t, kv.Key, actual[i].Key)
					assert.Equal(t, kv.Value, actual[i].Value)
				}
			}
		})
	}
}

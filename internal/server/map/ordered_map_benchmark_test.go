package orderedmap

import (
	"fmt"
	"testing"
)

func BenchmarkOrderedMap_Add(b *testing.B) {
	om := NewOrderedMap()
	opChan := make(chan Operation)

	go om.Run(opChan)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}
}

func BenchmarkOrderedMap_Get(b *testing.B) {
	om := NewOrderedMap()
	opChan := make(chan Operation)

	go om.Run(opChan)

	// Prepopulate the map
	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resultChan := make(chan interface{})
		opChan <- Operation{
			Action: Get,
			Key:    fmt.Sprintf("key%d", i),
			Result: resultChan,
		}
		<-resultChan
	}
}

func BenchmarkOrderedMap_Delete(b *testing.B) {
	om := NewOrderedMap()
	opChan := make(chan Operation)

	go om.Run(opChan)

	// Prepopulate the map
	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Delete,
			Key:    fmt.Sprintf("key%d", i),
		}
	}
}

func BenchmarkOrderedMap_GetAll(b *testing.B) {
	om := NewOrderedMap()
	opChan := make(chan Operation)

	go om.Run(opChan)

	// Prepopulate the map
	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resultChan := make(chan interface{})
		opChan <- Operation{
			Action: GetAll,
			Result: resultChan,
		}
		<-resultChan
	}
}

// Verify O(1) complexity requirement.

func benchmarkAdd(b *testing.B, n int) {
	om := NewOrderedMap()
	opChan := make(chan Operation)
	go om.Run(opChan)

	for i := 0; i < n; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i+n),
			Value:  fmt.Sprintf("value%d", i+n),
		}
	}
}

func BenchmarkOrderedMap_Add_1K(b *testing.B)   { benchmarkAdd(b, 1000) }
func BenchmarkOrderedMap_Add_10K(b *testing.B)  { benchmarkAdd(b, 10000) }
func BenchmarkOrderedMap_Add_100K(b *testing.B) { benchmarkAdd(b, 100000) }

func benchmarkGet(b *testing.B, n int) {
	om := NewOrderedMap()
	opChan := make(chan Operation)
	go om.Run(opChan)

	for i := 0; i < n; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resultChan := make(chan interface{})
		opChan <- Operation{
			Action: Get,
			Key:    fmt.Sprintf("key%d", i%n),
			Result: resultChan,
		}
		<-resultChan
	}
}

func BenchmarkOrderedMap_Get_1K(b *testing.B)   { benchmarkGet(b, 1000) }
func BenchmarkOrderedMap_Get_10K(b *testing.B)  { benchmarkGet(b, 10000) }
func BenchmarkOrderedMap_Get_100K(b *testing.B) { benchmarkGet(b, 100000) }

func benchmarkDelete(b *testing.B, n int) {
	om := NewOrderedMap()
	opChan := make(chan Operation)
	go om.Run(opChan)

	for i := 0; i < n; i++ {
		opChan <- Operation{
			Action: Add,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		opChan <- Operation{
			Action: Delete,
			Key:    fmt.Sprintf("key%d", i%n),
		}
	}
}

func BenchmarkOrderedMap_Delete_1K(b *testing.B)   { benchmarkDelete(b, 1000) }
func BenchmarkOrderedMap_Delete_10K(b *testing.B)  { benchmarkDelete(b, 10000) }
func BenchmarkOrderedMap_Delete_100K(b *testing.B) { benchmarkDelete(b, 100000) }

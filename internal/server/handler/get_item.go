package handler

import (
	"fmt"
	orderedmap "lightblocks/internal/server/map"
	"os"
)

func GetItemHandler(om *orderedmap.OrderedMap, key string) {
	value, exists := om.Get(key)
	if exists {
		f, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(fmt.Sprintf("key: %s, value: %s\n", key, value))
	}
}

func GetAllItemsHandler(om *orderedmap.OrderedMap) {
	f, _ := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	for _, kv := range om.GetAll() {
		f.WriteString(fmt.Sprintf("key: %s, value: %s\n", kv.Key, kv.Value))
	}
}

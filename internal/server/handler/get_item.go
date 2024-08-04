package handler

import (
	"fmt"
	"lightblocks/internal/observer"
	orderedmap "lightblocks/internal/server/map"
	"os"
	"path/filepath"
)

func GetItemHandler(obs observer.Observer, opChan chan<- orderedmap.Operation, key string) {
	resultChan := make(chan interface{})
	opChan <- orderedmap.Operation{
		Action: orderedmap.Get,
		Key:    key,
		Result: resultChan,
	}

	result := (<-resultChan).(struct {
		Value  string
		Exists bool
	})

	obs = obs.Str("key", key)

	if !result.Exists {
		obs.Warn("key not found: ", key)
		return
	}

	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		obs.Err(err).Error("failed to open output file")
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("key: %s, value: %s\n", key, result.Value)); err != nil {
		obs.Err(err).Error("failed to write to output file")
	} else {
		obs.Info("successfully wrote to output file")
	}
}

func GetAllItemsHandler(obs observer.Observer, opChan chan<- orderedmap.Operation) {
	resultChan := make(chan interface{})
	opChan <- orderedmap.Operation{
		Action: orderedmap.GetAll,
		Result: resultChan,
	}

	result := (<-resultChan).([]orderedmap.KeyValuePair)

	f, err := os.OpenFile("./output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		obs.Err(err).Error("failed to open output file")
		return
	}
	defer f.Close()

	fullPath, err := filepath.Abs(f.Name())
	if err != nil {
		obs.Err(err).Error("failed to get absolute path of output file")
	} else {
		obs = obs.Str("output_file_path", fullPath)
	}

	for _, kv := range result {
		if _, err = f.WriteString(fmt.Sprintf("key: %s, value: %s\n", kv.Key, kv.Value)); err != nil {
			obs.Err(err).Error("failed to write to output file")
		}
	}

	obs.Info("successfully wrote all items to output file")
}

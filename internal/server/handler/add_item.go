package handler

import (
	orderedmap "lightblocks/internal/server/map"
)

func AddItemHandler(opChan chan<- orderedmap.Operation, key, value string) {
	opChan <- orderedmap.Operation{
		Action: "add",
		Key:    key,
		Value:  value,
	}
}

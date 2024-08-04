package handler

import (
	orderedmap "lightblocks/internal/server/map"
)

func DeleteItemHandler(opChan chan<- orderedmap.Operation, key string) {
	opChan <- orderedmap.Operation{
		Action: "delete",
		Key:    key,
	}
}

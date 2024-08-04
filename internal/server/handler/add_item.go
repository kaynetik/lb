package handler

import (
	orderedmap "lightblocks/internal/server/map"
)

func AddItemHandler(om *orderedmap.OrderedMap, key, value string) {
	om.Add(key, value)
}

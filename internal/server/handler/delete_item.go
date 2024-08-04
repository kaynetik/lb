package handler

import (
	orderedmap "lightblocks/internal/server/map"
)

func DeleteItemHandler(om *orderedmap.OrderedMap, key string) {
	om.Delete(key)
}

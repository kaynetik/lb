package orderedmap

type KeyValuePair struct {
	Key   string
	Value string
}

type node struct {
	kv   KeyValuePair
	prev *node
	next *node
}

type OrderedMap struct {
	data map[string]*node
	head *node
	tail *node
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		data: make(map[string]*node),
	}
}

type Operation struct {
	Action string
	Key    string
	Value  string
	Result chan interface{}
}

func (om *OrderedMap) Run(opChan <-chan Operation) {
	for op := range opChan {
		switch op.Action {
		case "add":
			om.add(op.Key, op.Value)
		case "delete":
			om.delete(op.Key)
		case "get":
			value, exists := om.get(op.Key)
			op.Result <- struct {
				Value  string
				Exists bool
			}{value, exists}
		case "getAll":
			op.Result <- om.getAll()
		}
	}
}

func (om *OrderedMap) add(key, value string) {
	if n, exists := om.data[key]; exists {
		n.kv.Value = value
		return
	}
	newNode := &node{
		kv: KeyValuePair{
			Key:   key,
			Value: value,
		},
	}
	if om.tail == nil {
		om.head = newNode
		om.tail = newNode
	} else {
		om.tail.next = newNode
		newNode.prev = om.tail
		om.tail = newNode
	}
	om.data[key] = newNode
}

func (om *OrderedMap) delete(key string) {
	if n, exists := om.data[key]; exists {
		if n.prev != nil {
			n.prev.next = n.next
		} else {
			om.head = n.next
		}
		if n.next != nil {
			n.next.prev = n.prev
		} else {
			om.tail = n.prev
		}
		delete(om.data, key)
	}
}

func (om *OrderedMap) get(key string) (string, bool) {
	if n, exists := om.data[key]; exists {
		return n.kv.Value, true
	}
	return "", false
}

func (om *OrderedMap) getAll() []KeyValuePair {
	var result []KeyValuePair
	for n := om.head; n != nil; n = n.next {
		result = append(result, n.kv)
	}
	return result
}

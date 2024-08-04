package main

import (
	"lightblocks/internal/server/handler"
	orderedmap "lightblocks/internal/server/map"
	"lightblocks/internal/server/queue"
	"log"
	"strings"
)

func handleCommand(command string, om *orderedmap.OrderedMap) {
	parts := strings.Fields(command)
	switch parts[0] {
	case "addItem":
		if len(parts) == 3 {
			handler.AddItemHandler(om, parts[1], parts[2])
		}
	case "deleteItem":
		if len(parts) == 2 {
			handler.DeleteItemHandler(om, parts[1])
		}
	case "getItem":
		if len(parts) == 2 {
			handler.GetItemHandler(om, parts[1])
		}
	case "getAllItems":
		handler.GetAllItemsHandler(om)
	default:
		log.Printf("Unknown command: %s", command)
	}
}

func main() {
	om := orderedmap.NewOrderedMap()
	rabbitMQ, err := queue.NewRabbitMQ("commandQueue")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	err = rabbitMQ.Consume(func(message string) {
		go handleCommand(message, om)
	})
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	// Block forever
	select {}
}

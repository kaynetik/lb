package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"lightblocks/internal/observer"
	"lightblocks/internal/server/handler"
	"lightblocks/internal/server/map"
	"lightblocks/internal/server/queue"
	"log"
	"strings"
)

var tracer trace.Tracer

func init() {
	observer.InitObserver("server", "", "dev")
	tracer = otel.Tracer("server")
}

func handleCommand(command string, opChan chan<- orderedmap.Operation) {
	obs, _ := observer.Action(context.Background(), tracer, "handleCommand")
	defer obs.Close()

	parts := strings.Fields(command)
	switch parts[0] {
	case "addItem":
		if len(parts) == 3 {
			handler.AddItemHandler(opChan, parts[1], parts[2])
		}
	case "deleteItem":
		if len(parts) == 2 {
			handler.DeleteItemHandler(opChan, parts[1])
		}
	case "getItem":
		if len(parts) == 2 {
			handler.GetItemHandler(obs, opChan, parts[1])
		}
	case "getAllItems":
		handler.GetAllItemsHandler(obs, opChan)
	default:
		obs.Warn("unknown command: ", command)
	}
}

func main() {
	om := orderedmap.NewOrderedMap()
	opChan := make(chan orderedmap.Operation)

	go om.Run(opChan)

	rabbitMQ, err := queue.NewRabbitMQ("commandQueue")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	err = rabbitMQ.Consume(func(message string) {
		go handleCommand(message, opChan)
	})
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	// Block forever
	select {}
}

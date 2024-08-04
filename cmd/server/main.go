package main

import (
	"context"
	"lightblocks/internal/server"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"lightblocks/internal/observer"
	"lightblocks/internal/server/handler"
	"lightblocks/internal/server/map"
	"lightblocks/internal/server/queue"
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
	cfg := server.ParseConfig()
	observer.InitObserver(cfg.Name, cfg.OTELDigesterURL, cfg.Environment)
	obs, _ := observer.Action(context.Background(), tracer)

	om := orderedmap.NewOrderedMap()
	opChan := make(chan orderedmap.Operation)

	go om.Run(opChan)

	rabbitMQ, err := queue.NewRabbitMQ(cfg.DialTarget, cfg.QueueName)
	if err != nil {
		obs.Err(err).Fatal("failed to connect to RabbitMQ")
	}
	defer rabbitMQ.Close()

	err = rabbitMQ.Consume(func(message string) {
		go handleCommand(message, opChan) // If bigger lag were to occur this would be a bottleneck.
		// To handle such scenario we could go the easy route, via a worker-pool [POND](github.com/alitto/pond) with an `Eager` strategy.
	})
	if err != nil {
		obs.Err(err).Fatal("failed to start consuming messages")
	}

	select {}
}

package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"lightblocks/internal/client"
	"lightblocks/internal/observer"
)

func main() {
	cfg := client.ParseConfig()
	observer.InitObserver(cfg.Name, cfg.OTELDigesterURL, cfg.Environment)
	obs, ctx := observer.Action(context.Background(), tracer)

	commands, err := client.ReadCommands(ctx, cfg.InputFile)
	if err != nil {
		obs.Err(err).Fatal("failed to read commands")
	}

	sender, err := client.NewSender(cfg.DialTarget, cfg.QueueName)
	if err != nil {
		obs.Err(err).Fatal("failed to connect to RabbitMQ")
	}
	defer sender.Close()

	for _, command := range commands {
		obs = obs.Str("command", command)
		obs.Debug("sending command")

		if err = sender.Send(command); err != nil {
			obs.Err(err).Error("failed to send command")
		}

		obs.Debug("command sent")
	}
}

var tracer trace.Tracer

func init() {
	tracer = otel.GetTracerProvider().Tracer("cmd/client")
}

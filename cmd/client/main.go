package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"lightblocks/internal/client"
	"lightblocks/internal/observer"
)

func main() {
	config := client.ParseConfig()

	// TODO: Read digester and environment from the env vars.
	observer.InitObserver("client", "TBD", "dev")

	obs, ctx := observer.Action(context.Background(), tracer)

	commands, err := client.ReadCommands(ctx, config.InputFile)
	if err != nil {
		obs.Err(err).Fatal("failed to read commands")
	}

	sender, err := client.NewSender(config.QueueName)
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

var (
	tracer trace.Tracer
)

func init() {
	tracer = otel.GetTracerProvider().Tracer("cmd/client")
}

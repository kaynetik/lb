package client

import (
	"bufio"
	"context"
	"lightblocks/internal/observer"
	"os"
)

// NOTE: Observer that can be observed here is intended purely as a demonstration of a different ways to gather
// tracing data (primarily intended for the `Tempo` system).

// ReadCommands reads commands from the input file or stdin.
func ReadCommands(ctx context.Context, inputFile string) ([]string, error) {
	obs, ctx := observer.Action(ctx, tracer)
	if inputFile == "" {
		return readCommandsFromStdin(ctx, obs)
	}

	return readCommandsFromFile(ctx, obs, inputFile)
}

func readCommandsFromStdin(ctx context.Context, obs observer.Observer) ([]string, error) {
	obs, _ = obs.ChildAction(ctx, tracer)

	var commands []string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		obs.Err(err).Error("failed to read commands from stdin")

		return nil, err
	}

	obs.Debug("successfully read commands from stdin")

	return commands, nil
}

func readCommandsFromFile(ctx context.Context, obs observer.Observer, inputFile string) ([]string, error) {
	obs, _ = obs.ChildAction(ctx, tracer)
	obs = obs.Str("input_file", inputFile)

	file, err := os.Open(inputFile)
	if err != nil {
		obs.Err(err).Error("failed to open input file")

		return nil, err
	}
	defer file.Close()

	var commands []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		obs.Err(err).Error("failed to read commands from file")

		return nil, err
	}

	obs.Debug("successfully read commands from file")

	return commands, nil
}

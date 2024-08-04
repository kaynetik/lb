package main

import (
	"lightblocks/internal/client"
	"log"
)

func main() {
	config := client.ParseConfig()

	commands, err := client.ReadCommands(config.InputFile)
	if err != nil {
		log.Fatalf("Failed to read commands: %v", err)
	}

	sender, err := client.NewSender(config.QueueName)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer sender.Close()

	for _, command := range commands {
		log.Println("Sending command", command)
		err := sender.Send(command)
		if err != nil {
			log.Printf("Failed to send command: %v", err)
		}
		log.Println("Command sent")
	}
}

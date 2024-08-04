package client

import (
	"bufio"
	"os"
)

func ReadCommands(inputFile string) ([]string, error) {
	if inputFile == "" {
		return readCommandsFromStdin()
	}
	return readCommandsFromFile(inputFile)
}

func readCommandsFromStdin() ([]string, error) {
	var commands []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return commands, nil
}

func readCommandsFromFile(inputFile string) ([]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return commands, nil
}

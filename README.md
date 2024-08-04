# Mapper

## Description

This project demonstrates a Golang-based client-server application using RabbitMQ for message queuing. The server
processes commands to manage an in-memory ordered map, while the client sends commands from the command line or a file.

## Features

Features

- Server
    - In-memory ordered map
    - Parallel command execution
    - Supports add, delete, get, and get all items operations
- Client:
    - Sends commands to RabbitMQ
    - Configurable via command line or file

### TODOs

- [ ] Add one benchmark test
- [ ] Initialize `promtail` and `loki` for logging
- [ ] Initialize `tempo` digester and compactor

## Usage

### Requirements

Only pre-requisite is Docker with Docker Compose.

### Running the Application

1. Clone the repo:

```bash
git clone git@github.com:kaynetik/lb.git
cd lb
```

2. Run with Docker Compose:

```bash
docker-compose up --build
```

### Commands

Client Commands:

- `addItem('key', 'value')`
- `deleteItem('key')`
- `getItem('key')`
- `getAllItems()`

## Testing

To run the tests, execute the following command:

```bash
go test ./...
```

Or if you have available `tparse` :rocket: :

```bash
go test -json -race ./... | tparse -all
```
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

- [x] Add one benchmark test
- [x] Add multiple clients into `docker-compose` (primarily to demonstrate _"Clients can be added / removed / started
  while not inteferring to the server or other clients"_)
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

### Benchmarking

The benchmarks confirm that the `Add`, `Get`, and `Delete` operations have O(1) complexity as their time per operation
remains
relatively constant with the increasing number of items in the map.
The `GetAll` operation, as expected, exhibits `O(n)`complexity. Was your requirement to have it as `O(1)` a mistake?
IMHO `GetAll` operation is expected to be `O(n)` since it needs to iterate over all elements, and the benchmark reflects
this with a higher time per operation.

### Benchmark Results

| Benchmark                       | Runs    | ns/op  | B/op    | allocs/op |
|---------------------------------|---------|--------|---------|-----------|
| BenchmarkOrderedMap_Add         | 1710042 | 717.8  | 240     | 5         |
| BenchmarkOrderedMap_Get         | 1794475 | 789.6  | 143     | 3         |
| BenchmarkOrderedMap_Delete      | 2568500 | 449.8  | 23      | 1         |
| BenchmarkOrderedMap_GetAll      | 10000   | 218295 | 1471192 | 21        |
| BenchmarkOrderedMap_Add_1K      | 1638547 | 722.2  | 174     | 5         |
| BenchmarkOrderedMap_Add_10K     | 1663308 | 691.6  | 173     | 5         |
| BenchmarkOrderedMap_Add_100K    | 1618531 | 731.4  | 244     | 5         |
| BenchmarkOrderedMap_Get_1K      | 2344520 | 504.3  | 133     | 3         |
| BenchmarkOrderedMap_Get_10K     | 2380174 | 514.9  | 135     | 3         |
| BenchmarkOrderedMap_Get_100K    | 2206395 | 562.9  | 135     | 3         |
| BenchmarkOrderedMap_Delete_1K   | 5105319 | 224.8  | 13      | 1         |
| BenchmarkOrderedMap_Delete_10K  | 5138942 | 232.8  | 15      | 1         |
| BenchmarkOrderedMap_Delete_100K | 4704482 | 238.5  | 15      | 1         |


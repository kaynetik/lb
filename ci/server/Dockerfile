# Stage 1: Build the application
FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Stage 2: Create the runtime container
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]

#!/bin/bash

echo "Starting gRPC Server..."
go run server/main.go &

echo "Starting gRPC Client..."
go run main.go &

echo "Starting Gin API Server..."
go run gin/main.go &

wait

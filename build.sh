#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Running tests..."
go test ./... -v

echo "Building binary..."
# Note: On Linux, binaries typically don't have an .exe extension
go build -o golab ./cmd/golab

echo "Building Docker image..."
docker build -t golab:latest .

echo "Done."

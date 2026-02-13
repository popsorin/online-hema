#!/bin/bash

set -e

echo "Running tests..."

docker run --rm \
  -v "$(pwd)":/app \
  -w /app \
  golang:1.22-alpine \
  sh -c "go test -v ./..."

echo "Tests completed!"

#!/bin/bash

set -e

# Run tests inside the API container
# This script assumes the containers are running and environment variables are set

echo "Running tests..."

# Check that required environment variables are set
required_vars=("TEST_DATABASE_HOST" "TEST_DATABASE_PORT" "TEST_DATABASE_USER" "TEST_DATABASE_PASSWORD" "TEST_DATABASE_DBNAME")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        echo "Please set all TEST_DATABASE_* environment variables before running tests"
        exit 1
    fi
done

# We need to use the golang builder image to run tests with all dependencies
docker run --rm \
  --network hema-lessons_default \
  -v "$(pwd)":/app \
  -w /app \
  -e TEST_DATABASE_HOST \
  -e TEST_DATABASE_PORT \
  -e TEST_DATABASE_USER \
  -e TEST_DATABASE_PASSWORD \
  -e TEST_DATABASE_DBNAME \
  -e TEST_DATABASE_SSLMODE="${TEST_DATABASE_SSLMODE:-disable}" \
  -e TEST_APP_ENVIRONMENT="${TEST_APP_ENVIRONMENT:-testing}" \
  golang:1.22-alpine \
  sh -c "apk add --no-cache git && cd /app && go mod tidy && go test -v ./internal/handlers/..."

echo "Tests completed!"

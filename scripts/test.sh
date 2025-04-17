#!/bin/bash

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if the container is running
if ! docker ps | grep -q go_app_db; then
    echo "PostgreSQL container is not running. Starting it..."
    docker-compose up -d
    # Wait for the container to be ready
    sleep 5
fi

# Run the tests
go test -v ./services/player/... 
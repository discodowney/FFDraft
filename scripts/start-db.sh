#!/bin/bash

# Start the PostgreSQL container
docker-compose up -d

# Wait for the container to be ready
sleep 5

echo "PostgreSQL container is running" 
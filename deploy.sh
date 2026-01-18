#!/bin/bash

echo "Starting deployment..."

# 1. Pull latest code
# echo "Pulling latest code..."
# git pull origin main

# 2. Build the application
echo "Building binary..."
if go build -v -o flowgo-server ./cmd/server/main.go; then
    echo "Build success!"
else
    echo "Build failed!"
    exit 1
fi

# 3. Simulate Restart (In production you would use systemctl restart flowgo)
echo "Simulating restart..."
# systemctl restart flowgo
echo "Deployment completed successfully!"

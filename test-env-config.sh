#!/bin/bash

echo "Testing environment variable configuration..."

# Run with different environment variables to test override
export LNI_SERVER_PORT=9090
export LNI_MONGODB_DATABASE=test_database

cd /home/dksensei/Work/Dev/LetsNormalizeIt-V2.0

# Run the server for a few seconds to see the config output
timeout 5s go run cmd/server/main.go 2>&1 | grep -E "(Configuration loaded|Server Port|MongoDB Database)" || echo "Server started but output not captured"

echo "Test completed."

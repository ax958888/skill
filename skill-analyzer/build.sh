#!/bin/bash
set -e

echo "Building skill-analyzer..."

# Download dependencies
go mod download

# Build static binary
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-analyzer cmd/main.go

# Verify
echo "Verifying binary..."
ls -lh skill-analyzer

if command -v ldd &> /dev/null; then
    ldd skill-analyzer || echo "Static binary confirmed"
fi

echo "Build complete: skill-analyzer"

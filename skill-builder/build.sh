#!/bin/bash
set -e

echo "Building skill-builder..."

# Download dependencies
go mod download

# Build static binary
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-builder cmd/main.go

# Verify
echo "Verifying binary..."
ls -lh skill-builder

if command -v ldd &> /dev/null; then
    ldd skill-builder || echo "Static binary confirmed"
fi

echo "Build complete: skill-builder"

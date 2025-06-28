#!/bin/bash

# == Build ContainDB for Linux (Bin File) ==

set -e

BINARY_PATH="./bin/ContainDB"
BUILD_OUTPUT_DIR=$(dirname "$BINARY_PATH")

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo "Go not found. Installing Go via snap..."
  sudo snap install go --classic
else
  echo "Go is already installed."
fi

# Create output directory if it doesn't exist
mkdir -p "$BUILD_OUTPUT_DIR"

# Build the Go project with -o flag
echo "Building the project..."
go build -o "$BINARY_PATH" ./src/Core

echo "Build complete. Binary available at $BINARY_PATH"

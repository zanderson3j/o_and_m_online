#!/bin/bash
set -e

echo "Building server..."
cd server
go build -o server .

echo "Build complete!"
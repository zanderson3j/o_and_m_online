#!/bin/bash
set -e

echo "Building server..."
cd server
go build -o server .

echo "Building web client..."
cd ..
mkdir -p server/web
GOOS=js GOARCH=wasm go build -o server/web/game.wasm .
cp index.html server/web/
cp wasm_exec.js server/web/

echo "Build complete!"
#!/bin/bash

echo "Building for Render deployment..."

# Create web directory in server folder
mkdir -p server/web

# Build WASM
echo "Building WASM..."
GOOS=js GOARCH=wasm go build -o server/web/game.wasm .

# Copy web files
echo "Copying web files..."
cp index.html server/web/
cp wasm_exec.js server/web/

# Build server
echo "Building server..."
cd server
go build -o server .
cd ..

echo ""
echo "Build complete! The server directory is ready for deployment."
echo ""
echo "To deploy to Render:"
echo "1. Commit and push these changes to your git repository"
echo "2. Render will automatically rebuild and deploy"
echo ""
echo "Your game will be available at:"
echo "- WebSocket: wss://o-and-m-online.onrender.com/ws"
echo "- Web Client: https://o-and-m-online.onrender.com"
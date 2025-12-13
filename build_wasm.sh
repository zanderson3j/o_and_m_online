#!/bin/bash

echo "Building WASM version..."

# Build the WASM binary
GOOS=js GOARCH=wasm go build -o game.wasm .

if [ $? -eq 0 ]; then
    echo "WASM build successful!"
    echo ""
    echo "To run the web version locally:"
    echo "1. Start a local web server:"
    echo "   python3 -m http.server 8000"
    echo "   or"
    echo "   npx http-server -p 8000"
    echo ""
    echo "2. Open http://localhost:8000 in your browser"
else
    echo "WASM build failed!"
    exit 1
fi
#!/bin/bash

echo "Starting web server on http://localhost:8000"
echo "Press Ctrl+C to stop"
echo ""

# Check if python3 is available
if command -v python3 &> /dev/null; then
    python3 -m http.server 8000
elif command -v python &> /dev/null; then
    python -m http.server 8000
else
    echo "Python not found. Please install Python or use another web server."
    echo "Alternative: npx http-server -p 8000"
    exit 1
fi
#!/bin/bash

# O&M Game Room Installation Helper
echo "O&M Game Room - Installation Helper"
echo "===================================="
echo ""
echo "This script will help install O&M Game Room on your Mac."
echo ""

# Check if app exists
if [ -d "/Applications/O&M Game Room.app" ]; then
    echo "Found O&M Game Room in Applications folder."
    echo "Removing quarantine flags..."
    xattr -cr "/Applications/O&M Game Room.app"
    echo "✅ Done! You can now open O&M Game Room normally."
else
    echo "❌ O&M Game Room not found in Applications folder."
    echo "Please drag it from the DMG to Applications first."
fi

echo ""
echo "Press any key to close..."
read -n 1
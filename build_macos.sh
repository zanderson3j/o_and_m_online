#!/bin/bash

# Build script for macOS
set -e

APP_NAME="O&M Game Room"
BUNDLE_ID="com.oliveandmillie.gameroom"
VERSION="1.0.0"

echo "Building macOS app..."

# Create app bundle structure
APP_BUNDLE="build/darwin/${APP_NAME}.app"
rm -rf "$APP_BUNDLE"
mkdir -p "$APP_BUNDLE/Contents/MacOS"
mkdir -p "$APP_BUNDLE/Contents/Resources"

# Build for Intel
echo "Building for Intel (amd64)..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "build/darwin/game_client_amd64" .

# Build for Apple Silicon
echo "Building for Apple Silicon (arm64)..."
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "build/darwin/game_client_arm64" .

# Create universal binary
echo "Creating universal binary..."
lipo -create "build/darwin/game_client_amd64" "build/darwin/game_client_arm64" -output "$APP_BUNDLE/Contents/MacOS/game_client"

# Copy resources
cp resources/Info.plist "$APP_BUNDLE/Contents/"
cp resources/app.icns "$APP_BUNDLE/Contents/Resources/"

# Set executable permissions
chmod +x "$APP_BUNDLE/Contents/MacOS/game_client"

# Create DMG for distribution
echo "Creating DMG..."
DMG_NAME="OandM_Game_Room_${VERSION}.dmg"
DMG_PATH="build/darwin/$DMG_NAME"

# Create temporary DMG directory
DMG_DIR="build/darwin/dmg"
rm -rf "$DMG_DIR"
mkdir -p "$DMG_DIR"

# Copy app to DMG directory
cp -R "$APP_BUNDLE" "$DMG_DIR/"

# Create Applications symlink
ln -s /Applications "$DMG_DIR/Applications"

# Create DMG
hdiutil create -volname "$APP_NAME" -srcfolder "$DMG_DIR" -ov -format UDZO "$DMG_PATH"

# Clean up
rm -rf "$DMG_DIR"
rm -f "build/darwin/game_client_amd64" "build/darwin/game_client_arm64"

echo "Build complete! DMG created at: $DMG_PATH"
echo ""
echo "To install:"
echo "1. Open $DMG_NAME"
echo "2. Drag '$APP_NAME' to Applications folder"
echo "3. Launch from Applications or Spotlight"
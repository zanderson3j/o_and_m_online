#!/bin/bash

# Build script for macOS - current architecture only
set -e

APP_NAME="O&M Game Room"
BUNDLE_ID="com.oliveandmillie.gameroom"
VERSION="1.0.4"

echo "Building macOS app for current architecture..."

# Create app bundle structure
APP_BUNDLE="build/darwin/${APP_NAME}.app"
rm -rf "$APP_BUNDLE"
mkdir -p "$APP_BUNDLE/Contents/MacOS"
mkdir -p "$APP_BUNDLE/Contents/Resources"

# Build for current architecture
echo "Building binary..."
go build -o "$APP_BUNDLE/Contents/MacOS/game_client" .

# Copy resources
cp resources/Info.plist "$APP_BUNDLE/Contents/"
cp resources/app.icns "$APP_BUNDLE/Contents/Resources/"

# Set executable permissions
chmod +x "$APP_BUNDLE/Contents/MacOS/game_client"

# Ad-hoc sign the app (no Apple Developer account needed)
echo "Signing app..."
codesign --force --deep --sign - "$APP_BUNDLE"

# Create DMG for distribution
echo "Creating DMG..."
DMG_NAME="OandM_Game_Room_${VERSION}_$(uname -m).dmg"
DMG_PATH="build/darwin/$DMG_NAME"

# Create temporary DMG directory
DMG_DIR="build/darwin/dmg"
rm -rf "$DMG_DIR"
mkdir -p "$DMG_DIR"

# Copy app to DMG directory
cp -R "$APP_BUNDLE" "$DMG_DIR/"

# Copy install helper
cp install_helper.command "$DMG_DIR/"

# Create Applications symlink
ln -s /Applications "$DMG_DIR/Applications"

# Create DMG with internet-enable to reduce quarantine issues
hdiutil create -volname "$APP_NAME" -srcfolder "$DMG_DIR" -ov -format UDZO "$DMG_PATH"
hdiutil internet-enable -yes "$DMG_PATH" || true

# Clean up
rm -rf "$DMG_DIR"

echo "Build complete! DMG created at: $DMG_PATH"
echo ""
echo "Architecture: $(uname -m)"
echo ""
echo "To install:"
echo "1. Open $DMG_NAME"
echo "2. Drag '$APP_NAME' to Applications folder"
echo "3. Launch from Applications or Spotlight"
echo ""
echo "Note: This build is for your current architecture only."
echo "For universal builds, you'll need to build on both Intel and Apple Silicon machines."
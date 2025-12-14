# O&M Game Room - Application Packaging Guide

## Overview

The game is packaged as a native macOS application with:
- Custom app icon
- DMG installer for easy distribution
- Support for both Intel and Apple Silicon Macs
- Auto-update notifications with one-click download

## Building for macOS

### Prerequisites
- Go 1.21 or later
- Xcode Command Line Tools
- macOS 10.13 or later

### Build for Current Architecture
```bash
./build_macos_current.sh
```

This creates a DMG file in `build/darwin/` for your current Mac architecture.

### Build Universal Binary (requires both architectures)
To create a universal binary that runs on both Intel and Apple Silicon:
1. Build on Intel Mac: `./build_macos.sh`
2. Build on Apple Silicon Mac: `./build_macos.sh`
3. Combine using the universal build script

## Application Structure

```
O&M Game Room.app/
├── Contents/
│   ├── Info.plist          # App metadata and configuration
│   ├── MacOS/
│   │   └── game_client     # The executable
│   └── Resources/
│       └── app.icns        # App icon
```

## Distribution

### Manual Distribution
1. Build the DMG using the build script
2. Upload to your website or file hosting
3. Users download and install by dragging to Applications

### GitHub Releases (Automated)
1. Update version in `updater.go`
2. Commit and push changes
3. Create a git tag: `git tag v1.0.1 && git push --tags`
4. GitHub Actions will automatically build and create a release

## Auto-Updates

The app checks for updates on startup by querying GitHub releases.

### Current Implementation
- Checks GitHub releases API for new versions on startup
- Shows persistent notification in lobby (bottom-left corner)
- Clickable notification opens browser to download page
- Users can manually download and install updates

### Future Improvements
1. **Code Signing**: Sign the app for Gatekeeper approval ($99/year Apple Developer account)
2. **Automatic Installation**: Once signed, enable seamless background updates
3. **Sparkle Integration**: Use native macOS update framework for professional update experience

## Icon Requirements

The app uses `oam_icon.png` converted to `.icns` format with these sizes:
- 16x16, 32x32 (1x and 2x)
- 128x128, 256x256, 512x512 (1x and 2x)
- 1024x1024 (for App Store, if needed)

## Version Management

1. Update `currentVersion` in `updater.go`
2. Update `CFBundleShortVersionString` in `Info.plist`
3. Tag the release in git

## Troubleshooting

### "App is damaged" Error
If users see this error on first launch (expected for unsigned apps):
1. Right-click the app and select "Open"
2. Click "Open" in the security dialog
3. Or use Security & Privacy settings to "Open Anyway"

See INSTALL.txt (included in DMG) for detailed instructions.

### Auto-Update Not Working
Check:
- GitHub repository is public or API token is configured
- Version string follows semantic versioning (v1.0.0)
- DMG filename matches expected pattern

## Next Steps

1. **Code Signing**: Get an Apple Developer certificate ($99/year) for seamless updates
2. **Notarization**: Submit app to Apple for malware scanning
3. **Windows/Linux**: Extend packaging to other platforms
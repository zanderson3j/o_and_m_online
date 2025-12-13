# Building O&M Game Room for Intel Macs

If you have an Intel Mac and want to build from source:

## Prerequisites
1. Install Go: https://go.dev/dl/
2. Install Xcode Command Line Tools: `xcode-select --install`

## Build Steps

```bash
# Clone the repository
git clone https://github.com/zanderson3j/o_and_m_online.git
cd o_and_m_online

# Build the app
./build_macos_current.sh
```

The DMG will be created in `build/darwin/`

## Alternative: Request an Intel Build

If you can't build it yourself, please:
1. Open an issue on GitHub
2. Request an Intel build
3. We'll upload one for you

## Note
The next release (v1.0.2) will include automatic Intel builds via GitHub Actions.
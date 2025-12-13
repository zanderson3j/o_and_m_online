# Deploying to Render

This guide explains how to deploy both the game server and web client to Render.

## Setup

Your Render service is already configured to:
1. Run the WebSocket game server
2. Serve the web client files

## Deployment Steps

### Method 1: Automatic (via Git)

1. Run the build script locally to test:
   ```bash
   ./build_for_render.sh
   ```

2. Commit all changes:
   ```bash
   git add .
   git commit -m "Add web client support"
   git push
   ```

3. Render will automatically:
   - Detect the push
   - Run the build commands
   - Deploy the new version

### Method 2: Manual Deploy

If you need to trigger a manual deploy:

1. Go to your Render dashboard
2. Navigate to your service (o-and-m-online)
3. Click "Manual Deploy" > "Deploy latest commit"

## What Gets Deployed

The server now serves:
- **WebSocket endpoint**: `wss://o-and-m-online.onrender.com/ws`
- **Web client**: `https://o-and-m-online.onrender.com`

## Build Configuration

Render runs these commands (defined in render.yaml):

```bash
# Build server
cd server && go build -o server .

# Build web client
cd .. && mkdir -p server/web
GOOS=js GOARCH=wasm go build -o server/web/game.wasm .
cp index.html server/web/
cp wasm_exec.js server/web/
```

## Testing After Deployment

1. **Test the web client**:
   - Visit https://o-and-m-online.onrender.com
   - The game should load in your browser

2. **Test multiplayer**:
   - Open the URL in multiple browser tabs
   - Or use one desktop client and one web client
   - Both should connect and play together

## Troubleshooting

**Web client shows "Failed to load game"**:
- Check Render logs for build errors
- Ensure all three files (index.html, game.wasm, wasm_exec.js) are in server/web/

**WebSocket connection fails**:
- The web client connects to the same URL as the page (automatically uses wss:// on HTTPS)
- Check that the server is running by visiting the base URL

**Build fails on Render**:
- Check that Go version on Render supports WASM (Go 1.11+)
- Ensure wasm_exec.js is committed to the repository

## Important Files

- `server/main.go` - Updated to serve static files from `./web` directory
- `render.yaml` - Build configuration for Render
- `build_for_render.sh` - Local build script for testing
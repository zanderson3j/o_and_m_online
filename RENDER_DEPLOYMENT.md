# Deploying to Render

This guide explains how to deploy the game server to Render.

## Setup

Your Render service is configured to run the WebSocket game server.

## Deployment Steps

### Method 1: Automatic (via Git)

1. Commit all changes:
   ```bash
   git add .
   git commit -m "Your commit message"
   git push
   ```

2. Render will automatically:
   - Detect the push
   - Run the build commands
   - Deploy the new version

### Method 2: Manual Deploy

If you need to trigger a manual deploy:

1. Go to your Render dashboard
2. Navigate to your service (o-and-m-online)
3. Click "Manual Deploy" > "Deploy latest commit"

## What Gets Deployed

The server provides:
- **WebSocket endpoint**: `wss://o-and-m-online.onrender.com/ws`

## Build Configuration

Render runs these commands (defined in render.yaml):

```bash
# Build server
cd server && go build -o server .
```

## Testing After Deployment

1. **Test with desktop client**:
   - Run the desktop client locally
   - It should connect to `wss://o-and-m-online.onrender.com/ws`

2. **Test multiplayer**:
   - Run multiple desktop clients
   - They should be able to create/join rooms and play together

## Troubleshooting

**WebSocket connection fails**:
- Check that the server is running in Render logs
- Ensure the desktop client is using the correct server URL

**Build fails on Render**:
- Check Go version compatibility in render.yaml
- Review build logs in Render dashboard

## Important Files

- `server/main.go` - WebSocket server implementation
- `render.yaml` - Build configuration for Render
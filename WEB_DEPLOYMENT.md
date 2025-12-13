# Web Deployment Guide

This guide explains how to deploy the web version of Olive & Millie's Game Room.

## Local Development

1. Build the WASM version:
   ```bash
   ./build_wasm.sh
   ```

2. Start a local web server:
   ```bash
   ./serve_web.sh
   ```

3. Open http://localhost:8000 in your browser

## Production Deployment

### Option 1: Static Hosting (Netlify, Vercel, GitHub Pages)

The web version consists of static files that can be hosted anywhere:

**Required files:**
- `index.html`
- `game.wasm`
- `wasm_exec.js`

**Deployment steps:**

1. Build the WASM file:
   ```bash
   GOOS=js GOARCH=wasm go build -o game.wasm .
   ```

2. Upload these three files to your static hosting service

3. Ensure your hosting service serves `.wasm` files with the correct MIME type: `application/wasm`

### Option 2: Deploy to Netlify

1. Create a `netlify.toml` file:
   ```toml
   [[headers]]
     for = "*.wasm"
     [headers.values]
       Content-Type = "application/wasm"
   ```

2. Build and deploy:
   ```bash
   ./build_wasm.sh
   netlify deploy --prod --dir .
   ```

### Option 3: Deploy to Vercel

1. Create a `vercel.json` file:
   ```json
   {
     "headers": [
       {
         "source": "/(.*).wasm",
         "headers": [
           {
             "key": "Content-Type",
             "value": "application/wasm"
           }
         ]
       }
     ]
   }
   ```

2. Deploy:
   ```bash
   ./build_wasm.sh
   vercel --prod
   ```

### Option 4: Self-Hosted with Nginx

1. Build the WASM file
2. Copy files to your web root
3. Add to nginx config:
   ```nginx
   location ~ \.wasm$ {
     add_header Content-Type application/wasm;
   }
   ```

## Important Notes

- The web client connects to the same WebSocket server as the desktop client
- Ensure your server URL in `main.go` is accessible from browsers (wss:// for HTTPS sites)
- The game requires WebAssembly support (all modern browsers support this)
- For production, ensure your WebSocket server has proper CORS headers

## Troubleshooting

**Game won't load:**
- Check browser console for errors
- Ensure all three files (index.html, game.wasm, wasm_exec.js) are present
- Verify .wasm files are served with correct MIME type

**WebSocket connection fails:**
- Check if server URL is correct in main.go
- Ensure server allows connections from your domain
- For HTTPS sites, WebSocket must use wss:// (not ws://)

**Performance issues:**
- WASM performance is generally good but slightly slower than native
- Ensure browser hardware acceleration is enabled
- Some browsers perform better than others (Chrome/Edge typically fastest)
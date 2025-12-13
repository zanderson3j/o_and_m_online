# Olive & Millie's Game Room - Online Multiplayer

A collection of 2-player games with online multiplayer support using websockets.

## Features

- **4 Games**: Yahtzee, Santorini, Connect Four, Memory Match
- **Online Multiplayer**: Play with friends over the network
- **Lobby System**: Create or join rooms for each game type
- **Studio Ghibli Theme**: Whimsical forest aesthetic with kodama spirits
- **Cross-Platform**: Desktop app and web browser support

## How to Run

### 1. Start the Server

```bash
cd server
./server
```

The server will start on port 8080.

### 2. Start the Client(s)

#### Desktop Client
In separate terminals (one for each player):

```bash
./game_room_online
```

#### Web Client
1. Build the web version:
   ```bash
   ./build_wasm.sh
   ```

2. Serve the files:
   ```bash
   ./serve_web.sh
   ```

3. Open http://localhost:8000 in your browser

Both clients will automatically connect to `ws://localhost:8080/ws`.

## How to Play

1. **Select a Game**: Click on any game type button in the lobby
2. **Join or Create Room**:
   - If rooms exist, you'll see a list to join
   - If no rooms exist, a new one is created automatically
3. **Wait for Player**: The room needs 2 players to start
4. **Start Game**: Once both players are in, click "START GAME"
5. **Play**: The game will begin for both players

## Architecture

### Server (`server/main.go`)
- WebSocket server managing rooms and players
- Handles lobby operations (create/join/leave rooms)
- Broadcasts game state between players
- Runs on port 8080

### Client
- **network.go**: WebSocket client with message handling
- **lobby.go**: Lobby UI for creating/joining rooms
- **main.go**: Main game loop with network integration

### Game Flow

```
Client connects to server
  ↓
Lobby screen shows available games
  ↓
Player selects game type
  ↓
Create/join room for that game
  ↓
Wait for 2 players
  ↓
Both players click "START GAME"
  ↓
Server broadcasts start message
  ↓
Game begins on both clients
```

## Network Protocol

Messages are JSON-formatted with the following types:

- `create_room`: Create a new game room
- `join_room`: Join an existing room
- `leave_room`: Leave current room
- `start_game`: Begin the game (requires 2 players)
- `game_move`: Send a game action to opponent
- `room_list`: Server sends list of available rooms
- `player_joined/left`: Room status updates

## Next Steps (TODO)

The current implementation handles:
- ✅ Server with lobby system
- ✅ Client with network layer
- ✅ Lobby UI
- ✅ Room creation/joining
- ✅ Game starting

Still needed for full multiplayer gameplay:
- ⏳ Games emitting move events
- ⏳ Games receiving opponent moves
- ⏳ Synchronizing game state between players
- ⏳ Handling disconnections gracefully
- ⏳ Turn validation on server

## Development

To modify the server URL, edit `serverURL` constant in `main.go`:

```go
const serverURL = "ws://localhost:8080/ws"
```

For remote play, replace `localhost` with the server's IP address.

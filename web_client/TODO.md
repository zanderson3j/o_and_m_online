# Web Client TODO List - Match Desktop Client Behavior

## 1. Home Screen Fixes
- [ ] Remove automatic room creation when clicking game buttons
- [ ] Add "GO ONLINE" button that only shows when offline
- [ ] Update status text to show "Offline Mode - Click 'Go Online' to connect"
- [ ] Auto-connect on startup (like desktop does after 100ms)
- [ ] Show "ONLINE" status when connected

## 2. Create Proper Lobby Screen
- [ ] Create a proper lobby screen that shows after connecting
- [ ] Show "ONLINE MULTIPLAYER LOBBY" title
- [ ] Show 4 game buttons (same as home screen)
- [ ] Show current avatar in bottom-right with "Click avatar to change" text
- [ ] Clicking game button should show room list, not create room

## 3. Implement Room List View
- [ ] Show "AVAILABLE ROOMS" title
- [ ] List existing rooms for selected game type
- [ ] Filter out started rooms
- [ ] Show room name and player count
- [ ] Format player count differently for 2-player vs multi-player games
- [ ] Add "CREATE NEW ROOM" button
- [ ] Add "BACK" button to return to lobby

## 4. Implement "In Room" View (Waiting Room)
- [ ] Show "IN GAME ROOM" title
- [ ] Display room name
- [ ] Show player list with avatars
- [ ] Show appropriate waiting messages based on game type and player count
- [ ] Enable "START GAME" button when minimum players reached
- [ ] Add "BACK" button to leave room

## 5. Fix Avatar Selection
- [ ] Create separate avatar selection screen
- [ ] Show all 15 avatars in grid layout
- [ ] Highlight currently selected avatar
- [ ] Send avatar update when selection changes
- [ ] Add "BACK" button

## 6. Fix Network Message Handling
- [ ] Handle all message types with correct field names (capitalized)
- [ ] Implement proper room creation flow (create → server responds → auto-join)
- [ ] Handle player_joined messages correctly
- [ ] Handle player_left messages
- [ ] Handle game_ended messages
- [ ] Update player lists when receiving player_update messages

## 7. Implement Game State Management
- [ ] Store selected game type when viewing rooms
- [ ] Clear room info when leaving a room
- [ ] Properly transition between screens
- [ ] Handle disconnection gracefully
- [ ] Implement reconnection logic

## 8. Fix Screen Navigation
- [ ] Home (offline) → Click "Go Online" → Lobby
- [ ] Lobby → Click game → Room List
- [ ] Room List → Join/Create room → In Room
- [ ] In Room → Start game → Game Screen
- [ ] Any screen → Click avatar → Avatar Selection

## 9. UI/UX Improvements
- [ ] Add proper button hover states
- [ ] Show connection status consistently
- [ ] Add loading states during network operations
- [ ] Show error messages appropriately
- [ ] Match the desktop client's visual style

## 10. Game Integration
- [ ] Ensure games receive correct player data on start
- [ ] Handle game state updates from server
- [ ] Send game moves to server
- [ ] Handle player disconnection during games

## Current State vs Expected State

### Current (Wrong):
- Home Screen → Click Game → Attempts to create room immediately

### Expected (Correct):
- Home Screen → Click "Go Online" → Online Lobby → Click Game → Room List → Create/Join Room → Waiting Room → Start Game

## Implementation Order
1. Fix home screen to add "GO ONLINE" button
2. Create proper lobby screen
3. Implement room list view
4. Fix network message handling
5. Implement waiting room view
6. Add avatar selection screen
7. Complete game integration
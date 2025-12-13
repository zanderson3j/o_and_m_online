# Release Notes

## Version 1.0.0

First public release of O&M Game Room!

### Features
- 4 multiplayer games: Yahtzee, Santorini, Connect Four, and Memory Match
- Online multiplayer support
- Beautiful forest theme with Studio Ghibli-inspired visuals
- Support for up to 20 players (game dependent)
- Custom avatars for each player

### Installation

1. Download the DMG for your Mac:
   - **Apple Silicon** (M1/M2/M3): `OandM_Game_Room_1.0.0_apple_silicon.dmg`
   - **Intel**: `OandM_Game_Room_1.0.0_intel.dmg`

2. Open the downloaded DMG file

3. Drag "O&M Game Room" to your Applications folder

4. **IMPORTANT - First time opening:**
   
   **If you see "damaged and can't be opened":**
   - Open Terminal (find it in Applications > Utilities)
   - Copy and paste this command:
     ```
     xattr -cr "/Applications/O&M Game Room.app"
     ```
   - Press Enter
   - Now open the app normally
   
   **Alternative method:**
   - Right-click on "O&M Game Room" in Applications
   - Select "Open" from the menu
   - Click "Open" in the security dialog

5. Enjoy playing!

### Known Issues
- First launch requires right-click → Open due to no code signing (this is normal)
- Auto-update checks for new versions but requires manual download for now

### System Requirements
- macOS 10.13 or later
- Internet connection for multiplayer
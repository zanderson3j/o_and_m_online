package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type LobbyScreen struct {
	networkClient       *NetworkClient
	createButtons       []*Button // One button per game type
	roomButtons         []*Button // Buttons for joinable rooms
	createRoomButton    *Button   // Button to create new room
	backButton          *Button
	startButton         *Button
	avatarButtons       []*Button // Avatar selection buttons
	randomAvatarButton  *Button   // Random avatar selection button
	selectedGame        string
	selectedAvatar      AvatarType
	showingRooms        bool
	inRoom              bool
	waitingForGame      bool
	showAvatarSelect    bool
	updateMessageHovered bool
}

func NewLobbyScreen(nc *NetworkClient) *LobbyScreen {
	ls := &LobbyScreen{
		networkClient:  nc,
		createButtons:  make([]*Button, 4),
		roomButtons:    make([]*Button, 0),
		avatarButtons:  make([]*Button, int(AvatarNumTypes)),
		selectedAvatar: AvatarHuman,
		showingRooms:   false,
		inRoom:         false,
	}

	// Create buttons for each game type (4 games)
	games := []string{"YAHTZEE", "SANTORINI", "CONNECT FOUR", "MEMORY MATCH"}
	buttonWidth := 280.0
	buttonHeight := 90.0
	startX := float64(screenWidth/2) - buttonWidth - 20
	startY := 150.0  // Moved up from 250 to show more background
	spacingX := buttonWidth + 40.0
	spacingY := buttonHeight + 20.0

	for i, game := range games {
		row := i / 2
		col := i % 2
		x := startX + float64(col)*spacingX
		y := startY + float64(row)*spacingY

		ls.createButtons[i] = &Button{
			x:       x,
			y:       y,
			width:   buttonWidth,
			height:  buttonHeight,
			text:    game,
			enabled: true,
		}
	}

	// Back button
	ls.backButton = &Button{
		x:       20,
		y:       float64(screenHeight - 70),
		width:   150,
		height:  50,
		text:    "BACK",
		enabled: true,
	}

	// Start game button (when in room)
	ls.startButton = &Button{
		x:       float64(screenWidth/2) - 100,
		y:       400,  // Changed from screenHeight-100 to avoid potential overlap
		width:   200,
		height:  60,
		text:    "START GAME",
		enabled: true,
	}

	// Create new room button (when viewing room list)
	ls.createRoomButton = &Button{
		x:       float64(screenWidth/2) - 150,
		y:       float64(screenHeight - 100),
		width:   300,
		height:  60,
		text:    "CREATE NEW ROOM",
		enabled: true,
	}

	// Avatar selection buttons - arranged in grid
	avatarSize := 80.0
	avatarSpacing := 100.0
	avatarsPerRow := 8  // 8 avatars per row
	totalButtons := int(AvatarNumTypes) + 1  // +1 for random button
	numRows := (totalButtons + avatarsPerRow - 1) / avatarsPerRow

	// Center the grid
	totalWidth := float64(avatarsPerRow) * avatarSpacing - (avatarSpacing - avatarSize)
	avatarStartX := (float64(screenWidth) - totalWidth) / 2
	avatarStartY := 200.0

	for i := 0; i < int(AvatarNumTypes); i++ {
		row := i / avatarsPerRow
		col := i % avatarsPerRow

		// Center last row if it has fewer buttons (including random)
		xOffset := 0.0
		if row == numRows-1 {
			buttonsInLastRow := totalButtons % avatarsPerRow
			if buttonsInLastRow == 0 {
				buttonsInLastRow = avatarsPerRow
			}
			lastRowWidth := float64(buttonsInLastRow) * avatarSpacing - (avatarSpacing - avatarSize)
			xOffset = (totalWidth - lastRowWidth) / 2
		}

		ls.avatarButtons[i] = &Button{
			x:       avatarStartX + float64(col)*avatarSpacing + xOffset,
			y:       avatarStartY + float64(row)*120,  // 120 pixels between rows (avatar + name)
			width:   avatarSize,
			height:  avatarSize,
			text:    "", // No text, we'll draw avatars instead
			enabled: true,
		}
	}

	// Random avatar button - positioned as the next item in the grid
	randomIndex := int(AvatarNumTypes)
	randomRow := randomIndex / avatarsPerRow
	randomCol := randomIndex % avatarsPerRow

	// Apply same centering logic for last row
	randomXOffset := 0.0
	if randomRow == numRows-1 {
		buttonsInLastRow := totalButtons % avatarsPerRow
		if buttonsInLastRow == 0 {
			buttonsInLastRow = avatarsPerRow
		}
		lastRowWidth := float64(buttonsInLastRow) * avatarSpacing - (avatarSpacing - avatarSize)
		randomXOffset = (totalWidth - lastRowWidth) / 2
	}

	ls.randomAvatarButton = &Button{
		x:       avatarStartX + float64(randomCol)*avatarSpacing + randomXOffset,
		y:       avatarStartY + float64(randomRow)*120,
		width:   avatarSize,
		height:  avatarSize,
		text:    "", // No text, we'll draw a special icon
		enabled: true,
	}

	// Register network handlers
	nc.RegisterHandler(MsgStartGame, func(msg Message) {
		ls.waitingForGame = false
		// Game will be started by the handler in main.go
	})

	nc.RegisterHandler(MessageType("player_joined"), func(msg Message) {
		// Only set inRoom if we're the one who joined
		if msg.PlayerID == nc.GetPlayerID() {
			ls.inRoom = true
			ls.waitingForGame = false
			ls.showingRooms = false
			// Update network client's current room
			nc.mu.Lock()
			nc.currentRoom = msg.RoomID
			nc.mu.Unlock()
		}
	})

	nc.RegisterHandler(MessageType("room_created"), func(msg Message) {
		ls.inRoom = true
		ls.waitingForGame = false  // Don't set to true until we actually start the game
		ls.showingRooms = false
		// Update network client's current room
		nc.mu.Lock()
		nc.currentRoom = msg.RoomID
		nc.mu.Unlock()
	})

	nc.RegisterHandler(MessageType("player_left"), func(msg Message) {
		// If we're the one who left, reset state
		if msg.PlayerID == nc.GetPlayerID() {
			ls.inRoom = false
			ls.showingRooms = false
			ls.waitingForGame = false
			// Clear network client's current room
			nc.mu.Lock()
			nc.currentRoom = ""
			nc.mu.Unlock()
		}
	})

	nc.RegisterHandler(MessageType("game_ended"), func(msg Message) {
		// Reset lobby state when game ends
		ls.inRoom = false
		ls.showingRooms = false
		ls.waitingForGame = false
		// Clear network client's current room
		nc.mu.Lock()
		nc.currentRoom = ""
		nc.mu.Unlock()
	})
	
	return ls
}

func (ls *LobbyScreen) Reset() {
	ls.inRoom = false
	ls.showingRooms = false
	ls.waitingForGame = false
	ls.selectedGame = ""
}

// ShowAvatarSelection opens the avatar selection screen
func (ls *LobbyScreen) ShowAvatarSelection() {
	ls.showAvatarSelect = true
}

func (ls *LobbyScreen) Update(gr *GameRoom) error {
	mx, my := ebiten.CursorPosition()

	// Update hover state for update message (bottom left)
	if gr.updateAvailable {
		updateMsgX := float64(20)
		updateMsgY := float64(screenHeight - 40)
		updateMsgWidth := float64(300)
		updateMsgHeight := float64(30)
		ls.updateMessageHovered = float64(mx) >= updateMsgX && float64(mx) <= updateMsgX+updateMsgWidth &&
			float64(my) >= updateMsgY && float64(my) <= updateMsgY+updateMsgHeight
	}

	// Handle update message click
	if gr.updateAvailable && ls.updateMessageHovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Printf("Opening update URL: %s", gr.updateURL)
		OpenBrowser(gr.updateURL)
		// Exit the app so it can be overwritten by the update
		os.Exit(0)
	}

	if ls.showAvatarSelect {
		// Avatar selection mode
		for _, btn := range ls.avatarButtons {
			btn.hovered = btn.Contains(mx, my)
		}
		ls.randomAvatarButton.hovered = ls.randomAvatarButton.Contains(mx, my)
		ls.backButton.hovered = ls.backButton.Contains(mx, my)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i, btn := range ls.avatarButtons {
				if btn.hovered {
					ls.selectedAvatar = AvatarType(i)
					ls.networkClient.SetAvatar(i)
					ls.showAvatarSelect = false
				}
			}
			// Handle random avatar button
			if ls.randomAvatarButton.hovered {
				// Pick a random avatar (from 0 to AvatarNumTypes-1)
				randomAvatar := rand.Intn(int(AvatarNumTypes))
				ls.selectedAvatar = AvatarType(randomAvatar)
				ls.networkClient.SetAvatar(randomAvatar)
				ls.showAvatarSelect = false
			}
			if ls.backButton.hovered {
				ls.showAvatarSelect = false
			}
		}
		return nil
	}

	if ls.inRoom {
		// In a room - update start button
		ls.startButton.hovered = ls.startButton.Contains(mx, my)
		ls.backButton.hovered = ls.backButton.Contains(mx, my)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if ls.startButton.hovered {
				ls.networkClient.StartGame()
				ls.waitingForGame = true
			} else if ls.backButton.hovered {
				ls.networkClient.LeaveRoom()
				ls.inRoom = false
				ls.showingRooms = false
			}
		}
	} else if ls.showingRooms {
		// Showing available rooms - update room buttons
		ls.updateRoomButtons()
		for _, btn := range ls.roomButtons {
			btn.hovered = btn.Contains(mx, my)
		}
		ls.backButton.hovered = ls.backButton.Contains(mx, my)
		ls.createRoomButton.hovered = ls.createRoomButton.Contains(mx, my)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i, btn := range ls.roomButtons {
				if btn.hovered {
					rooms := ls.networkClient.GetRooms()
					availableRooms := make([]RoomInfo, 0)
					for _, room := range rooms {
						if room.GameType == ls.selectedGame && !room.Started {
							// For multi-player games, always show if not started
							// For 2-player games, only show if not full
							if room.MaxPlayers > 2 || room.Players < room.MaxPlayers {
								availableRooms = append(availableRooms, room)
							}
						}
					}
					if i < len(availableRooms) {
						log.Printf("Joining room %s", availableRooms[i].ID)
						ls.networkClient.JoinRoom(availableRooms[i].ID)
						// Don't set inRoom here - wait for the player_joined message
					}
				}
			}
			if ls.createRoomButton.hovered {
				roomName := fmt.Sprintf("%s Room", ls.selectedGame)
				ls.networkClient.CreateRoom(ls.selectedGame, roomName)
				// Don't set inRoom here - wait for the room_created message
			}
			if ls.backButton.hovered {
				ls.showingRooms = false
			}
		}
	} else {
		// Main lobby - update create game buttons
		for _, btn := range ls.createButtons {
			btn.hovered = btn.Contains(mx, my)
		}

		// Check if clicked on current avatar (to change it)
		avatarX := float64(screenWidth) - 100
		avatarY := float64(screenHeight) - 100
		if mx >= int(avatarX) && mx <= int(avatarX)+50 && 
		   my >= int(avatarY) && my <= int(avatarY)+50 {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				ls.showAvatarSelect = true
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i, btn := range ls.createButtons {
				if btn.hovered {
					games := []string{"yahtzee", "santorini", "connect_four", "memory"}
					ls.selectedGame = games[i]
					// Try to join existing room first, or create new one
					ls.showRoomsForGame(games[i])
				}
			}
		}
	}

	return nil
}

func (ls *LobbyScreen) showRoomsForGame(gameType string) {
	// Always show room list so players can choose which room to join
	// or create a new one
	ls.showingRooms = true
}

func (ls *LobbyScreen) getRoomDisplayText(room RoomInfo) string {
	// For games that support many players, show range
	if room.MaxPlayers > 2 {
		if room.GameType == "yahtzee" || room.GameType == "memory" {
			return fmt.Sprintf("%s (%d players, 1-%d)", room.Name, room.Players, room.MaxPlayers)
		}
	}
	// For 2-player games, show traditional format
	return fmt.Sprintf("%s (%d/%d)", room.Name, room.Players, room.MaxPlayers)
}

func (ls *LobbyScreen) updateRoomButtons() {
	rooms := ls.networkClient.GetRooms()
	availableRooms := make([]RoomInfo, 0)

	for _, room := range rooms {
		if room.GameType == ls.selectedGame && !room.Started {
			// For multi-player games, always show if not started
			// For 2-player games, only show if not full
			if room.MaxPlayers > 2 || room.Players < room.MaxPlayers {
				availableRooms = append(availableRooms, room)
			}
		}
	}

	ls.roomButtons = make([]*Button, len(availableRooms))
	buttonWidth := 400.0
	buttonHeight := 60.0
	startY := 150.0

	for i := range availableRooms {
		ls.roomButtons[i] = &Button{
			x:       float64(screenWidth/2) - buttonWidth/2,
			y:       startY + float64(i)*70,
			width:   buttonWidth,
			height:  buttonHeight,
			text:    ls.getRoomDisplayText(availableRooms[i]),
			enabled: true,
		}
	}
}

func (ls *LobbyScreen) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)

	if ls.showAvatarSelect {
		ls.drawAvatarSelection(screen)
	} else if ls.inRoom {
		ls.drawRoomWaiting(screen)
	} else if ls.showingRooms {
		ls.drawRoomList(screen)
	} else {
		ls.drawGameSelection(screen)
	}

	// Draw update notification (bottom left, always visible in lobby)
	if gr.updateAvailable {
		ls.drawUpdateMessage(screen, gr)
	}
}

func (ls *LobbyScreen) drawGameSelection(screen *ebiten.Image) {
	// Title
	titleWidth := float32(400)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "ONLINE MULTIPLAYER LOBBY"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	// Info box
	infoWidth := float32(500)
	infoX := float32(screenWidth/2) - infoWidth/2
	vector.DrawFilledRect(screen, infoX, 70, infoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, infoX, 70, infoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	status := "Select a game to play online"
	if !ls.networkClient.IsConnected() {
		status = "Connecting to server..."
	}

	statusX := int(infoX + (infoWidth-float32(len(status)*6))/2)
	ebitenutil.DebugPrintAt(screen, status, statusX, 90)

	// Draw game selection buttons
	for _, btn := range ls.createButtons {
		ls.drawButton(screen, btn)
	}

	// Draw current avatar in bottom right
	avatarX := float64(screenWidth) - 100
	avatarY := float64(screenHeight) - 100
	DrawAvatar(screen, ls.selectedAvatar, float32(avatarX), float32(avatarY), 1)
	
	// Draw "Click to change" text
	changeText := "Click avatar to change"
	textX := int(avatarX - float64(len(changeText)*3) + 25)
	ebitenutil.DebugPrintAt(screen, changeText, textX, int(avatarY)-20)
}

func (ls *LobbyScreen) drawRoomList(screen *ebiten.Image) {
	// Title
	titleWidth := float32(400)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "AVAILABLE ROOMS"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	// Room list
	if len(ls.roomButtons) == 0 {
		infoText := "No rooms available"
		infoX := screenWidth/2 - len(infoText)*3
		ebitenutil.DebugPrintAt(screen, infoText, infoX, 250)
	} else {
		for _, btn := range ls.roomButtons {
			ls.drawButton(screen, btn)
		}
	}

	// Create new room button
	ls.drawButton(screen, ls.createRoomButton)

	// Back button
	ls.drawButton(screen, ls.backButton)
}

func (ls *LobbyScreen) drawRoomWaiting(screen *ebiten.Image) {
	// Title
	titleWidth := float32(400)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "IN GAME ROOM"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	// Status
	rooms := ls.networkClient.GetRooms()
	currentRoom := ls.networkClient.GetCurrentRoom()

	var roomInfo *RoomInfo
	for _, room := range rooms {
		if room.ID == currentRoom {
			roomInfo = &room
			break
		}
	}

	statusY := 200
	if roomInfo == nil {
		errorText := "Room not found - please go back"
		ebitenutil.DebugPrintAt(screen, errorText, screenWidth/2-len(errorText)*3, statusY)
		ls.drawButton(screen, ls.backButton)
		return
	}
	
	statusText := fmt.Sprintf("Room: %s", roomInfo.Name)
	ebitenutil.DebugPrintAt(screen, statusText, screenWidth/2-len(statusText)*3, statusY)

	playerText := fmt.Sprintf("Players: %d/%d", roomInfo.Players, roomInfo.MaxPlayers)
	ebitenutil.DebugPrintAt(screen, playerText, screenWidth/2-len(playerText)*3, statusY+30)

	// Check if we can start the game
	canStart := false
	if roomInfo.GameType == "yahtzee" || roomInfo.GameType == "memory" {
		// Multi-player games can start with 1+ players
		canStart = roomInfo.Players >= 1
	} else {
		// 2-player games need exactly 2 players
		canStart = roomInfo.Players == roomInfo.MaxPlayers
	}
	
	if !canStart {
		// Only show waiting message for 2-player games
		waitText := "Waiting for another player..."
		ebitenutil.DebugPrintAt(screen, waitText, screenWidth/2-len(waitText)*3, statusY+80)
	} else if !ls.waitingForGame {
		var readyText string
		if roomInfo.GameType == "yahtzee" || roomInfo.GameType == "memory" {
			if roomInfo.Players == 1 {
				readyText = "Ready to start solo or wait for more players!"
			} else {
				readyText = fmt.Sprintf("Ready with %d players! Start or wait for more.", roomInfo.Players)
			}
		} else {
			readyText = "Ready to start!"
		}
		ebitenutil.DebugPrintAt(screen, readyText, screenWidth/2-len(readyText)*3, statusY+80)
		ls.drawButton(screen, ls.startButton)
	} else {
		startingText := "Starting game..."
		ebitenutil.DebugPrintAt(screen, startingText, screenWidth/2-len(startingText)*3, statusY+80)
	}

	// Player ID
	playerID := ls.networkClient.GetPlayerID()
	if playerID != "" {
		idText := fmt.Sprintf("Your ID: %s", playerID[:8])
		ebitenutil.DebugPrintAt(screen, idText, 20, screenHeight-20)
	}

	// Back button
	ls.drawButton(screen, ls.backButton)
}

func (ls *LobbyScreen) drawButton(screen *ebiten.Image, btn *Button) {
	bgColor := color.RGBA{30, 50, 80, 255}
	borderColor := color.RGBA{100, 150, 220, 255}

	if btn.hovered {
		bgColor = color.RGBA{50, 80, 120, 255}
		borderColor = color.RGBA{150, 200, 255, 255}
	}

	if !btn.enabled {
		bgColor = color.RGBA{40, 40, 50, 255}
		borderColor = color.RGBA{80, 80, 90, 255}
	}

	x := float32(btn.x)
	y := float32(btn.y)
	w := float32(btn.width)
	h := float32(btn.height)

	vector.DrawFilledRect(screen, x, y, w, h, bgColor, false)
	vector.StrokeRect(screen, x, y, w, h, 2, borderColor, false)

	textX := int(x + (w-float32(len(btn.text)*6))/2)
	textY := int(y + h/2 - 4)

	ebitenutil.DebugPrintAt(screen, btn.text, textX, textY)
	if btn.enabled {
		ebitenutil.DebugPrintAt(screen, btn.text, textX+1, textY)
	}
}

func (ls *LobbyScreen) drawAvatarSelection(screen *ebiten.Image) {
	// Title
	titleWidth := float32(400)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "SELECT YOUR AVATAR"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	// Info text
	infoText := "Choose your avatar:"
	infoX := screenWidth/2 - len(infoText)*3
	ebitenutil.DebugPrintAt(screen, infoText, infoX, 80)

	// Draw avatar options
	for i := 0; i < int(AvatarNumTypes); i++ {
		btn := ls.avatarButtons[i]
		
		// Draw background
		bgColor := color.RGBA{40, 40, 50, 255}
		borderColor := color.RGBA{100, 150, 220, 255}
		
		if AvatarType(i) == ls.selectedAvatar {
			bgColor = color.RGBA{60, 80, 120, 255}
			borderColor = color.RGBA{255, 220, 100, 255}
		} else if btn.hovered {
			bgColor = color.RGBA{50, 60, 80, 255}
			borderColor = color.RGBA{150, 200, 255, 255}
		}
		
		x := float32(btn.x)
		y := float32(btn.y)
		w := float32(btn.width)
		h := float32(btn.height)
		
		vector.DrawFilledRect(screen, x, y, w, h, bgColor, false)
		vector.StrokeRect(screen, x, y, w, h, 2, borderColor, false)
		
		// Draw avatar
		DrawAvatar(screen, AvatarType(i), x+15, y+15, 1)
		
		// Draw avatar name
		name := GetAvatarName(AvatarType(i))
		// For longer names, we need to adjust the centering
		nameWidth := len(name) * 6  // Each character is approximately 6 pixels wide
		nameX := int(x + w/2 - float32(nameWidth)/2)
		nameY := int(y + h + 5)
		ebitenutil.DebugPrintAt(screen, name, nameX, nameY)
	}

	// Draw random avatar button
	randomBtn := ls.randomAvatarButton
	randomBgColor := color.RGBA{50, 40, 60, 255}
	randomBorderColor := color.RGBA{150, 100, 200, 255}
	if randomBtn.hovered {
		randomBgColor = color.RGBA{70, 50, 90, 255}
		randomBorderColor = color.RGBA{200, 150, 255, 255}
	}

	randomX := float32(randomBtn.x)
	randomY := float32(randomBtn.y)
	randomW := float32(randomBtn.width)
	randomH := float32(randomBtn.height)

	vector.DrawFilledRect(screen, randomX, randomY, randomW, randomH, randomBgColor, false)
	vector.StrokeRect(screen, randomX, randomY, randomW, randomH, 2, randomBorderColor, false)

	// Draw simple human silhouette (similar to avatar style but simpler)
	scale := float32(1)
	avatarX := randomX + 15
	avatarY := randomY + 15

	humanColor := color.RGBA{200, 200, 200, 255}

	// Head (circle)
	vector.DrawFilledCircle(screen, avatarX+25*scale, avatarY+12*scale, 8*scale, humanColor, false)

	// Body
	vector.DrawFilledRect(screen, avatarX+20*scale, avatarY+20*scale, 10*scale, 15*scale, humanColor, false)

	// Arms
	vector.DrawFilledRect(screen, avatarX+12*scale, avatarY+22*scale, 8*scale, 4*scale, humanColor, false)
	vector.DrawFilledRect(screen, avatarX+30*scale, avatarY+22*scale, 8*scale, 4*scale, humanColor, false)

	// Legs
	vector.DrawFilledRect(screen, avatarX+20*scale, avatarY+35*scale, 4*scale, 10*scale, humanColor, false)
	vector.DrawFilledRect(screen, avatarX+26*scale, avatarY+35*scale, 4*scale, 10*scale, humanColor, false)

	// Double frame for random button
	vector.StrokeRect(screen, randomX, randomY, randomW, randomH, 3, randomBorderColor, false)
	vector.StrokeRect(screen, randomX+4, randomY+4, randomW-8, randomH-8, 2, randomBorderColor, false)

	// Draw "Random" label below
	randomLabel := "Random"
	labelWidth := len(randomLabel) * 6
	labelX := int(randomX + randomW/2 - float32(labelWidth)/2)
	labelY := int(randomY + randomH + 5)
	ebitenutil.DebugPrintAt(screen, randomLabel, labelX, labelY)

	// Back button
	ls.drawButton(screen, ls.backButton)
}

func (ls *LobbyScreen) drawUpdateMessage(screen *ebiten.Image, gr *GameRoom) {
	msgX := float32(20)
	msgY := float32(screenHeight - 40)
	msgWidth := float32(300)
	msgHeight := float32(30)

	// Background color - subtle green with hover effect
	bgColor := color.RGBA{40, 80, 40, 200}
	borderColor := color.RGBA{60, 120, 60, 255}
	if ls.updateMessageHovered {
		bgColor = color.RGBA{60, 100, 60, 220}
		borderColor = color.RGBA{80, 160, 80, 255}
	}

	// Draw background
	vector.DrawFilledRect(screen, msgX, msgY, msgWidth, msgHeight, bgColor, false)
	vector.StrokeRect(screen, msgX, msgY, msgWidth, msgHeight, 2, borderColor, false)

	// Draw text
	message := "Update " + gr.updateVersion + " available - click here"
	textX := int(msgX + 10)
	textY := int(msgY + msgHeight/2 - 4)
	ebitenutil.DebugPrintAt(screen, message, textX, textY)
}

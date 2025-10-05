package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type LobbyScreen struct {
	networkClient    *NetworkClient
	createButtons    []*Button // One button per game type
	roomButtons      []*Button // Buttons for joinable rooms
	createRoomButton *Button   // Button to create new room
	backButton       *Button
	startButton      *Button
	selectedGame     string
	showingRooms     bool
	inRoom           bool
	waitingForGame   bool
}

func NewLobbyScreen(nc *NetworkClient) *LobbyScreen {
	ls := &LobbyScreen{
		networkClient: nc,
		createButtons: make([]*Button, 5),
		roomButtons:   make([]*Button, 0),
		showingRooms:  false,
		inRoom:        false,
	}

	// Create buttons for each game type (5 games)
	games := []string{"YAHTZEE", "SANTORINI", "CONNECT FOUR", "MANCALA", "MEMORY MATCH"}
	buttonWidth := 280.0
	buttonHeight := 90.0
	startX := float64(screenWidth/2) - buttonWidth - 20
	startY := 180.0
	spacingX := buttonWidth + 40.0
	spacingY := buttonHeight + 20.0

	for i, game := range games {
		var x, y float64
		if i < 4 {
			row := i / 2
			col := i % 2
			x = startX + float64(col)*spacingX
			y = startY + float64(row)*spacingY
		} else {
			// Center the 5th button
			x = float64(screenWidth/2) - buttonWidth/2
			y = startY + spacingY*2
		}

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
		y:       float64(screenHeight - 100),
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

	// Register network handlers
	nc.RegisterHandler(MsgStartGame, func(msg Message) {
		ls.waitingForGame = false
		// Game will be started by the handler in main.go
	})

	nc.RegisterHandler("player_joined", func(msg Message) {
		ls.inRoom = true
		ls.waitingForGame = false
	})

	nc.RegisterHandler("room_created", func(msg Message) {
		ls.inRoom = true
		ls.waitingForGame = true
	})

	nc.RegisterHandler("player_left", func(msg Message) {
		// Reset lobby state when player leaves
		ls.inRoom = false
		ls.showingRooms = false
		ls.waitingForGame = false
	})

	nc.RegisterHandler("game_ended", func(msg Message) {
		// Reset lobby state when game ends
		log.Println("CLIENT: Received game_ended, resetting lobby state")
		ls.inRoom = false
		ls.showingRooms = false
		ls.waitingForGame = false
	})

	return ls
}

func (ls *LobbyScreen) Reset() {
	ls.inRoom = false
	ls.showingRooms = false
	ls.waitingForGame = false
	ls.selectedGame = ""
}

func (ls *LobbyScreen) Update(gr *GameRoom) error {
	mx, my := ebiten.CursorPosition()

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
						if room.GameType == ls.selectedGame && !room.Started && room.Players < room.MaxPlayers {
							availableRooms = append(availableRooms, room)
						}
					}
					if i < len(availableRooms) {
						ls.networkClient.JoinRoom(availableRooms[i].ID)
						ls.inRoom = true
						ls.showingRooms = false
					}
				}
			}
			if ls.createRoomButton.hovered {
				roomName := fmt.Sprintf("%s Room", ls.selectedGame)
				log.Printf("CLIENT: Creating new room for %s\n", ls.selectedGame)
				ls.networkClient.CreateRoom(ls.selectedGame, roomName)
				ls.inRoom = true
				ls.showingRooms = false
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

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			for i, btn := range ls.createButtons {
				if btn.hovered {
					games := []string{"yahtzee", "santorini", "connect_four", "mancala", "memory"}
					ls.selectedGame = games[i]
					log.Printf("CLIENT: Player clicked on %s game\n", games[i])
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

func (ls *LobbyScreen) updateRoomButtons() {
	rooms := ls.networkClient.GetRooms()
	availableRooms := make([]RoomInfo, 0)

	for _, room := range rooms {
		if room.GameType == ls.selectedGame && !room.Started && room.Players < room.MaxPlayers {
			availableRooms = append(availableRooms, room)
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
			text:    fmt.Sprintf("%s (%d/%d)", availableRooms[i].Name, availableRooms[i].Players, availableRooms[i].MaxPlayers),
			enabled: true,
		}
	}
}

func (ls *LobbyScreen) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)

	if ls.inRoom {
		ls.drawRoomWaiting(screen)
	} else if ls.showingRooms {
		ls.drawRoomList(screen)
	} else {
		ls.drawGameSelection(screen)
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
	if roomInfo != nil {
		statusText := fmt.Sprintf("Room: %s", roomInfo.Name)
		ebitenutil.DebugPrintAt(screen, statusText, screenWidth/2-len(statusText)*3, statusY)

		playerText := fmt.Sprintf("Players: %d/%d", roomInfo.Players, roomInfo.MaxPlayers)
		ebitenutil.DebugPrintAt(screen, playerText, screenWidth/2-len(playerText)*3, statusY+30)

		if roomInfo.Players < roomInfo.MaxPlayers {
			waitText := "Waiting for another player..."
			ebitenutil.DebugPrintAt(screen, waitText, screenWidth/2-len(waitText)*3, statusY+80)
		} else if !ls.waitingForGame {
			readyText := "Ready to start!"
			ebitenutil.DebugPrintAt(screen, readyText, screenWidth/2-len(readyText)*3, statusY+80)
			ls.drawButton(screen, ls.startButton)
		} else {
			startingText := "Starting game..."
			ebitenutil.DebugPrintAt(screen, startingText, screenWidth/2-len(startingText)*3, statusY+80)
		}
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	serverURL    = "wss://o-and-m-online.onrender.com/ws"
	// serverURL    = "ws://127.0.0.1:8080/ws" // Local testing
)

type ConnectionState int

const (
	StateConnecting ConnectionState = iota
	StateConnected
	StateFailed
)

type GameRoom struct {
	currentGame       GameInterface
	homeScreen        *HomeScreen
	lobbyScreen       *LobbyScreen
	networkClient     *NetworkClient
	isOnlineMode      bool
	updateAvailable   bool
	updateVersion     string
	updateURL         string
	connectionState   ConnectionState
	connectionError   string
}

func (gr *GameRoom) Update() error {
	if gr.isOnlineMode && gr.lobbyScreen != nil {
		return gr.lobbyScreen.Update(gr)
	}
	if gr.currentGame != nil {
		return gr.currentGame.Update(gr)
	}
	return gr.homeScreen.Update(gr)
}

func (gr *GameRoom) Draw(screen *ebiten.Image) {
	if gr.isOnlineMode && gr.lobbyScreen != nil {
		gr.lobbyScreen.Draw(screen, gr)
	} else if gr.currentGame != nil {
		gr.currentGame.Draw(screen, gr)
	} else {
		gr.homeScreen.Draw(screen, gr)
	}
}

func (gr *GameRoom) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (gr *GameRoom) SwitchToGame(game GameInterface) {
	gr.currentGame = game
}

func (gr *GameRoom) ReturnHome() {
	gr.currentGame = nil
	// Return to lobby if we have a network client
	if gr.networkClient != nil && gr.networkClient.IsConnected() {
		gr.isOnlineMode = true
		// Reset lobby state
		if gr.lobbyScreen != nil {
			gr.lobbyScreen.Reset()
		}
		// Leave current room if in one
		if gr.networkClient.GetCurrentRoom() != "" {
			gr.networkClient.LeaveRoom()
		}
	} else {
		gr.isOnlineMode = false
	}
}

func (gr *GameRoom) SwitchToOnline() {
	gr.isOnlineMode = true
}

func (gr *GameRoom) TryGoOnline() {
	log.Println("Attempting to connect to server...")

	// Set connecting state
	gr.connectionState = StateConnecting
	gr.connectionError = ""

	// Connect in goroutine to avoid blocking
	go func() {
		// Try to connect
		maxRetries := 3
		retryDelay := 2

		networkClient, err := NewNetworkClientWithRetry(serverURL, maxRetries, retryDelay)
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			gr.connectionState = StateFailed
			gr.connectionError = fmt.Sprintf("Connection failed: %v", err)
			return
		}

		log.Println("Connected successfully!")
		gr.connectionState = StateConnected
		gr.networkClient = networkClient
		gr.lobbyScreen = NewLobbyScreen(networkClient)
		gr.isOnlineMode = true
		// Send initial avatar selection
		networkClient.SetAvatar(0) // Default Human avatar

		// Register handlers
		networkClient.RegisterHandler(MsgStartGame, func(msg Message) {
			log.Printf("Starting game: %s\n", msg.GameType)

			// Get player number and game info from server
			var data struct {
				PlayerNumber int                      `json:"player_number"`
				TotalPlayers int                      `json:"total_players"`
				Players      []map[string]interface{} `json:"players"`
			}
			playerNum := 0
			totalPlayers := 2
			if err := json.Unmarshal(msg.Data, &data); err == nil {
				playerNum = data.PlayerNumber
				totalPlayers = data.TotalPlayers
			}
			log.Printf("I am player number: %d (total players: %d)\n", playerNum, totalPlayers)

			// Switch to the appropriate game with network support
			switch msg.GameType {
			case "yahtzee":
				gr.SwitchToGame(NewYahtzeeGameWithPlayers(networkClient, playerNum, data.Players))
			case "santorini":
				gr.SwitchToGame(NewSantoriniGameWithPlayers(networkClient, playerNum, data.Players))
			case "connect_four":
				gr.SwitchToGame(NewConnectFourGameWithPlayers(networkClient, playerNum, data.Players))
			case "memory":
				gr.SwitchToGame(NewMemoryGameWithPlayers(networkClient, playerNum, data.Players))
			}
			gr.isOnlineMode = false
		})

		networkClient.RegisterHandler("game_ended", func(msg Message) {
			log.Println("Game ended - player left")
			gr.ReturnHome()
		})
	}() // End of goroutine
}


func main() {
	log.Println("Starting Olive & Millie's Game Room")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Olive & Millie's Game Room - ONLINE")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	gameRoom := &GameRoom{
		homeScreen:      NewHomeScreen(),
		connectionState: StateConnecting,
	}

	// Check for updates (needs gameRoom to exist first)
	checkForUpdates(gameRoom)

	// Auto-connect on startup
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure UI is ready
		log.Println("Auto-connecting to server...")
		gameRoom.TryGoOnline()
	}()


	if err := ebiten.RunGame(gameRoom); err != nil {
		log.Fatal(err)
	}

	// Cleanup
	if gameRoom.networkClient != nil {
		gameRoom.networkClient.Close()
	}
}

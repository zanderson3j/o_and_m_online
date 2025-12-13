package main

import (
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	serverURL    = "wss://o-and-m-online.onrender.com/ws"
)

type GameRoom struct {
	currentGame   GameInterface
	homeScreen    *HomeScreen
	lobbyScreen   *LobbyScreen
	networkClient *NetworkClient
	isOnlineMode  bool
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

	// Try to connect with retries
	networkClient, err := NewNetworkClientWithRetry(serverURL, 10, 5)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}

	log.Println("Connected successfully!")
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
}

func main() {
	log.Println("Starting Olive & Millie's Game Room")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Olive & Millie's Game Room - ONLINE")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	gameRoom := &GameRoom{
		homeScreen: NewHomeScreen(),
	}

	// Try to connect to server (single attempt, no retries)
	log.Println("Attempting quick connection to server...")
	networkClient, err := NewNetworkClient(serverURL)
	if err != nil {
		log.Printf("Server not available: %v", err)
		log.Println("Starting in offline mode. Use 'Go Online' button to connect.")
		gameRoom.networkClient = nil
		gameRoom.isOnlineMode = false
	} else {
		log.Println("Connected to server successfully!")
		gameRoom.networkClient = networkClient
		gameRoom.lobbyScreen = NewLobbyScreen(networkClient)
		gameRoom.isOnlineMode = true
		// Send initial avatar selection
		networkClient.SetAvatar(0) // Default Human avatar
	}

	if networkClient != nil {
		// Register handler for when game starts
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
				gameRoom.SwitchToGame(NewYahtzeeGameWithPlayers(networkClient, playerNum, data.Players))
			case "santorini":
				gameRoom.SwitchToGame(NewSantoriniGameWithPlayers(networkClient, playerNum, data.Players))
			case "connect_four":
				gameRoom.SwitchToGame(NewConnectFourGameWithPlayers(networkClient, playerNum, data.Players))
			case "memory":
				gameRoom.SwitchToGame(NewMemoryGameWithPlayers(networkClient, playerNum, data.Players))
			}
			gameRoom.isOnlineMode = false
		})

		// Handle game ended (player left during game)
		networkClient.RegisterHandler("game_ended", func(msg Message) {
			log.Println("Game ended - player left")
			gameRoom.ReturnHome()
		})
	}

	if err := ebiten.RunGame(gameRoom); err != nil {
		log.Fatal(err)
	}

	// Cleanup
	if networkClient != nil {
		networkClient.Close()
	}
}

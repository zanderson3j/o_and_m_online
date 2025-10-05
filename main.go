package main

import (
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	serverURL    = "ws://localhost:8080/ws"
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

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Olive & Millie's Game Room - ONLINE")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Connect to server
	networkClient, err := NewNetworkClient(serverURL)
	if err != nil {
		log.Printf("Warning: Could not connect to server: %v", err)
		log.Println("Continuing in offline mode...")
	}

	gameRoom := &GameRoom{
		homeScreen:    NewHomeScreen(),
		networkClient: networkClient,
	}

	if networkClient != nil {
		gameRoom.lobbyScreen = NewLobbyScreen(networkClient)
		gameRoom.isOnlineMode = true

		// Register handler for when game starts
		networkClient.RegisterHandler(MsgStartGame, func(msg Message) {
			log.Printf("Starting game: %s\n", msg.GameType)

			// Get player number from server
			var data struct {
				PlayerNumber int `json:"player_number"`
			}
			playerNum := 0
			if err := json.Unmarshal(msg.Data, &data); err == nil {
				playerNum = data.PlayerNumber
			}
			log.Printf("I am player number: %d\n", playerNum)

			// Switch to the appropriate game with network support
			switch msg.GameType {
			case "yahtzee":
				gameRoom.SwitchToGame(NewYahtzeeGameWithNetwork(networkClient, playerNum))
			case "santorini":
				gameRoom.SwitchToGame(NewSantoriniGameWithNetwork(networkClient, playerNum))
			case "connect_four":
				gameRoom.SwitchToGame(NewConnectFourGameWithNetwork(networkClient, playerNum+1)) // Connect Four uses 1/2
			case "mancala":
				gameRoom.SwitchToGame(NewMancalaGameWithNetwork(networkClient, playerNum))
			case "memory":
				gameRoom.SwitchToGame(NewMemoryGameWithNetwork(networkClient, playerNum))
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

package main

import (
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// NetworkedGame wraps any game with network synchronization
type NetworkedGame struct {
	game          GameInterface
	networkClient *NetworkClient
	playerNumber  int // 0 or 1
	isMyTurn      bool
}

func NewNetworkedGame(game GameInterface, nc *NetworkClient) *NetworkedGame {
	ng := &NetworkedGame{
		game:          game,
		networkClient: nc,
		playerNumber:  0, // Will be determined by join order
		isMyTurn:      false,
	}

	// Register handler for receiving opponent moves
	nc.RegisterHandler(MsgGameMove, func(msg Message) {
		ng.handleOpponentMove(msg.Data)
	})

	return ng
}

func (ng *NetworkedGame) handleOpponentMove(data json.RawMessage) {
	// This will be overridden by specific game implementations
	log.Printf("Received opponent move: %s\n", string(data))
}

func (ng *NetworkedGame) sendMove(moveData interface{}) {
	if ng.networkClient != nil {
		ng.networkClient.SendGameMove(moveData)
	}
}

func (ng *NetworkedGame) Update(gr *GameRoom) error {
	return ng.game.Update(gr)
}

func (ng *NetworkedGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	ng.game.Draw(screen, gr)
}

func (ng *NetworkedGame) Reset() {
	ng.game.Reset()
}

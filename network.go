package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MsgJoinLobby    MessageType = "join_lobby"
	MsgLeaveLobby   MessageType = "leave_lobby"
	MsgCreateRoom   MessageType = "create_room"
	MsgJoinRoom     MessageType = "join_room"
	MsgLeaveRoom    MessageType = "leave_room"
	MsgStartGame    MessageType = "start_game"
	MsgGameMove     MessageType = "game_move"
	MsgGameState    MessageType = "game_state"
	MsgPlayerList   MessageType = "player_list"
	MsgRoomList     MessageType = "room_list"
	MsgError        MessageType = "error"
	MsgChat         MessageType = "chat"
)

type Message struct {
	Type      MessageType     `json:"type"`
	PlayerID  string          `json:"player_id,omitempty"`
	RoomID    string          `json:"room_id,omitempty"`
	GameType  string          `json:"game_type,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type RoomInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GameType   string `json:"game_type"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"max_players"`
	Started    bool   `json:"started"`
}

type NetworkClient struct {
	conn         *websocket.Conn
	playerID     string
	currentRoom  string
	rooms        []RoomInfo
	mu           sync.RWMutex
	msgHandlers  map[MessageType]func(Message)
	connected    bool
}

func NewNetworkClient(serverURL string) (*NetworkClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return nil, err
	}

	nc := &NetworkClient{
		conn:        conn,
		msgHandlers: make(map[MessageType]func(Message)),
		connected:   true,
	}

	// Start listening for messages
	go nc.listen()

	return nc, nil
}

func (nc *NetworkClient) listen() {
	defer func() {
		nc.mu.Lock()
		nc.connected = false
		nc.mu.Unlock()
		nc.conn.Close()
	}()

	for {
		var msg Message
		err := nc.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		nc.handleMessage(msg)
	}
}

func (nc *NetworkClient) handleMessage(msg Message) {
	// Handle special messages
	switch msg.Type {
	case "connected":
		nc.mu.Lock()
		nc.playerID = msg.PlayerID
		nc.mu.Unlock()
		log.Printf("Connected as player %s\n", msg.PlayerID)

	case MsgRoomList:
		var data struct {
			Rooms []RoomInfo `json:"rooms"`
		}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			nc.mu.Lock()
			nc.rooms = data.Rooms
			nc.mu.Unlock()
		}

	case "room_created":
		nc.mu.Lock()
		nc.currentRoom = msg.RoomID
		nc.mu.Unlock()
		log.Printf("Created room %s\n", msg.RoomID)

	case "player_joined":
		log.Printf("Player %s joined room\n", msg.PlayerID)
		// If we joined, update our current room
		if msg.PlayerID == nc.playerID {
			nc.mu.Lock()
			nc.currentRoom = msg.RoomID
			nc.mu.Unlock()
		}

	case "player_left":
		log.Printf("Player %s left room\n", msg.PlayerID)
		// If we're the one who left, clear our current room
		if msg.PlayerID == nc.playerID {
			nc.mu.Lock()
			nc.currentRoom = ""
			nc.mu.Unlock()
		}

	case "game_ended":
		log.Printf("Game ended in room %s\n", msg.RoomID)
		// Clear our current room when game ends
		nc.mu.Lock()
		nc.currentRoom = ""
		nc.mu.Unlock()

	case MsgError:
		var errData struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(msg.Data, &errData); err == nil {
			log.Printf("Server error: %s\n", errData.Error)
		}
	}

	// Call registered handler if exists
	nc.mu.RLock()
	handler, exists := nc.msgHandlers[msg.Type]
	nc.mu.RUnlock()

	if exists {
		handler(msg)
	}
}

func (nc *NetworkClient) RegisterHandler(msgType MessageType, handler func(Message)) {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	nc.msgHandlers[msgType] = handler
}

func (nc *NetworkClient) SendMessage(msg Message) error {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	if !nc.connected {
		log.Printf("CLIENT: Cannot send message %s - not connected!\n", msg.Type)
		return fmt.Errorf("not connected")
	}

	log.Printf("CLIENT: Sending message type=%s\n", msg.Type)
	return nc.conn.WriteJSON(msg)
}

func (nc *NetworkClient) CreateRoom(gameType, roomName string) error {
	data, _ := json.Marshal(map[string]string{
		"game_type": gameType,
		"room_name": roomName,
	})

	return nc.SendMessage(Message{
		Type:      MsgCreateRoom,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (nc *NetworkClient) JoinRoom(roomID string) error {
	data, _ := json.Marshal(map[string]string{
		"room_id": roomID,
	})

	return nc.SendMessage(Message{
		Type:      MsgJoinRoom,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (nc *NetworkClient) LeaveRoom() error {
	return nc.SendMessage(Message{
		Type:      MsgLeaveRoom,
		Timestamp: time.Now(),
	})
}

func (nc *NetworkClient) StartGame() error {
	return nc.SendMessage(Message{
		Type:      MsgStartGame,
		Timestamp: time.Now(),
	})
}

func (nc *NetworkClient) SendGameMove(moveData interface{}) error {
	data, err := json.Marshal(moveData)
	if err != nil {
		return err
	}

	return nc.SendMessage(Message{
		Type:      MsgGameMove,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (nc *NetworkClient) GetRooms() []RoomInfo {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.rooms
}

func (nc *NetworkClient) GetPlayerID() string {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.playerID
}

func (nc *NetworkClient) GetCurrentRoom() string {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.currentRoom
}

func (nc *NetworkClient) IsConnected() bool {
	nc.mu.RLock()
	defer nc.mu.RUnlock()
	return nc.connected
}

func (nc *NetworkClient) Close() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	if nc.connected {
		nc.conn.Close()
		nc.connected = false
	}
}

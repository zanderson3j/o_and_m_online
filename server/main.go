package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Log origin for debugging
		origin := r.Header.Get("Origin")
		log.Printf("WebSocket upgrade request from origin: %s", origin)
		return true // Allow all origins
	},
}

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
	MsgSetAvatar    MessageType = "set_avatar"
	MsgPlayerUpdate MessageType = "player_update"
)

type Message struct {
	Type      MessageType     `json:"type"`
	PlayerID  string          `json:"player_id,omitempty"`
	RoomID    string          `json:"room_id,omitempty"`
	GameType  string          `json:"game_type,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type Player struct {
	ID     string
	Name   string
	Avatar int
	Conn   *websocket.Conn
	RoomID string
	mu     sync.Mutex
}

type Room struct {
	ID         string
	Name       string
	GameType   string
	Players    []*Player
	MaxPlayers int
	Started    bool
	mu         sync.RWMutex
}

type Server struct {
	players map[string]*Player
	rooms   map[string]*Room
	mu      sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		players: make(map[string]*Player),
		rooms:   make(map[string]*Room),
	}
}

// Get max players for a game type
func getMaxPlayers(gameType string) int {
	switch gameType {
	case "yahtzee", "memory":
		return 20
	default:
		return 2
	}
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	playerID := generateID()
	player := &Player{
		ID:     playerID,
		Name:   avatarNames[0], // Default to "Human"
		Avatar: 0,              // Default to human avatar
		Conn:   conn,
	}

	s.mu.Lock()
	s.players[playerID] = player
	s.mu.Unlock()

	log.Printf("Player %s connected\n", playerID)

	// Send initial player ID
	s.sendMessage(player, Message{
		Type:      "connected",
		PlayerID:  playerID,
		Timestamp: time.Now(),
	})

	// Send current room list
	s.sendRoomList(player)

	// Handle messages from this player
	go s.handlePlayer(player)
}

func (s *Server) handlePlayer(player *Player) {
	defer func() {
		// Clean up on disconnect
		s.mu.Lock()
		if player.RoomID != "" {
			s.removePlayerFromRoom(player)
		}
		delete(s.players, player.ID)
		s.mu.Unlock()
		player.Conn.Close()
		log.Printf("Player %s disconnected\n", player.ID)
	}()

	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading from player %s: %v\n", player.ID, err)
			}
			break
		}

		msg.PlayerID = player.ID
		msg.Timestamp = time.Now()

		s.handleMessage(player, msg)
	}
}

func (s *Server) handleMessage(player *Player, msg Message) {
	log.Printf("SERVER: Received message type=%s from player %s\n", msg.Type, player.ID)

	switch msg.Type {
	case MsgCreateRoom:
		s.handleCreateRoom(player, msg)
	case MsgJoinRoom:
		s.handleJoinRoom(player, msg)
	case MsgLeaveRoom:
		s.handleLeaveRoom(player, msg)
	case MsgStartGame:
		s.handleStartGame(player, msg)
	case MsgGameMove:
		s.handleGameMove(player, msg)
	case MsgChat:
		s.handleChat(player, msg)
	case MsgSetAvatar:
		s.handleSetAvatar(player, msg)
	default:
		s.sendError(player, "Unknown message type")
	}
}

func (s *Server) handleCreateRoom(player *Player, msg Message) {
	var data struct {
		GameType string `json:"game_type"`
		RoomName string `json:"room_name"`
	}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		s.sendError(player, "Invalid create room data")
		return
	}

	s.mu.Lock()
	// Remove player from any existing room first
	if player.RoomID != "" {
		log.Printf("Player %s leaving old room %s to create new room\n", player.ID, player.RoomID)
		s.removePlayerFromRoom(player)
		s.broadcastRoomList()
	}

	roomID := generateID()
	room := &Room{
		ID:         roomID,
		Name:       data.RoomName,
		GameType:   data.GameType,
		Players:    []*Player{player},
		MaxPlayers: getMaxPlayers(data.GameType),
		Started:    false,
	}

	s.rooms[roomID] = room
	player.RoomID = roomID
	s.mu.Unlock()

	log.Printf("Player %s created room %s for game %s (Players: %d/%d)\n", player.ID, roomID, data.GameType, len(room.Players), room.MaxPlayers)

	// Send confirmation to creator
	s.sendMessage(player, Message{
		Type:      "room_created",
		RoomID:    roomID,
		GameType:  data.GameType,
		Timestamp: time.Now(),
	})

	// Broadcast updated room list to all players
	s.broadcastRoomList()
}

func (s *Server) handleJoinRoom(player *Player, msg Message) {
	var data struct {
		RoomID string `json:"room_id"`
	}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		s.sendError(player, "Invalid join room data")
		return
	}

	s.mu.Lock()
	// Remove player from any existing room first
	if player.RoomID != "" {
		log.Printf("Player %s leaving old room %s to join new room\n", player.ID, player.RoomID)
		s.removePlayerFromRoom(player)
		s.broadcastRoomList()
	}

	room, exists := s.rooms[data.RoomID]
	if !exists {
		s.mu.Unlock()
		s.sendError(player, "Room not found")
		return
	}

	room.mu.Lock()
	if len(room.Players) >= room.MaxPlayers {
		room.mu.Unlock()
		s.mu.Unlock()
		s.sendError(player, "Room is full")
		return
	}

	if room.Started {
		room.mu.Unlock()
		s.mu.Unlock()
		s.sendError(player, "Game already started")
		return
	}

	log.Printf("Before join: Room %s has %d players\n", data.RoomID, len(room.Players))
	room.Players = append(room.Players, player)
	player.RoomID = data.RoomID
	log.Printf("After join: Room %s has %d players\n", data.RoomID, len(room.Players))
	room.mu.Unlock()
	s.mu.Unlock()

	log.Printf("Player %s joined room %s (Players: %d/%d)\n", player.ID, data.RoomID, len(room.Players), room.MaxPlayers)

	// Notify all players in room
	s.broadcastToRoom(room, Message{
		Type:      "player_joined",
		PlayerID:  player.ID,
		RoomID:    data.RoomID,
		Timestamp: time.Now(),
	})

	// Update room list for everyone
	s.broadcastRoomList()
}

func (s *Server) handleLeaveRoom(player *Player, msg Message) {
	s.mu.Lock()

	if player.RoomID == "" {
		s.mu.Unlock()
		s.sendError(player, "Not in a room")
		return
	}

	log.Printf("Player %s leaving room %s\n", player.ID, player.RoomID)
	s.removePlayerFromRoom(player)
	s.mu.Unlock()

	s.broadcastRoomList()
}

func (s *Server) handleStartGame(player *Player, msg Message) {
	s.mu.RLock()
	room, exists := s.rooms[player.RoomID]
	s.mu.RUnlock()

	if !exists {
		s.sendError(player, "Not in a room")
		return
	}

	room.mu.Lock()
	// For games that support many players (Yahtzee, Memory), allow 1+ players
	// For 2-player only games, need exactly 2
	if room.MaxPlayers == 2 && len(room.Players) != 2 {
		room.mu.Unlock()
		s.sendError(player, "Need exactly 2 players to start")
		return
	}

	// For multi-player games, need at least 1 player (which we always have)
	if len(room.Players) < 1 {
		room.mu.Unlock()
		s.sendError(player, "Need at least 1 player to start")
		return
	}

	room.Started = true
	room.mu.Unlock()

	log.Printf("Game starting in room %s\n", room.ID)

	// Notify each player with their player number and all player info
	room.mu.RLock()

	// Build player info array
	playerInfos := make([]map[string]interface{}, len(room.Players))
	for i, p := range room.Players {
		playerInfos[i] = map[string]interface{}{
			"id":     p.ID,
			"name":   p.Name,
			"avatar": p.Avatar,
		}
	}

	for i, p := range room.Players {
		playerData, _ := json.Marshal(map[string]interface{}{
			"player_number": i, // 0 for first player, 1 for second
			"total_players": len(room.Players),
			"players":       playerInfos,
		})
		s.sendMessage(p, Message{
			Type:      MsgStartGame,
			RoomID:    room.ID,
			GameType:  room.GameType,
			Data:      playerData,
			Timestamp: time.Now(),
		})
	}
	room.mu.RUnlock()

	s.broadcastRoomList()
}

func (s *Server) handleGameMove(player *Player, msg Message) {
	s.mu.RLock()
	room, exists := s.rooms[player.RoomID]
	s.mu.RUnlock()

	if !exists {
		s.sendError(player, "Not in a room")
		return
	}

	// Broadcast move to all other players in room
	room.mu.RLock()
	for _, p := range room.Players {
		if p.ID != player.ID {
			s.sendMessage(p, msg)
		}
	}
	room.mu.RUnlock()
}

func (s *Server) handleChat(player *Player, msg Message) {
	s.mu.RLock()
	room, exists := s.rooms[player.RoomID]
	s.mu.RUnlock()

	if !exists {
		s.sendError(player, "Not in a room")
		return
	}

	// Broadcast chat to all players in room
	s.broadcastToRoom(room, msg)
}

// Avatar names matching client side
var avatarNames = []string{
	"Human", "Teddy", "Kaycat", "Zach Rabbit", "Kiraffe", "Owlive", "Milliepede", "Sweet Puppy Paw", "Tygler", "Chimpancici", "Papapus", "Kaitlynx", "Reagator", "Ocelivia", "Hen-ry", "Tomouse", "Karabou", "Valkyrie", "Eleanor", "Stella", "Huckleberry", "Winston", "Baxter",
}

func (s *Server) handleSetAvatar(player *Player, msg Message) {
	var data struct {
		Avatar int `json:"avatar"`
	}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		s.sendError(player, "Invalid avatar data")
		return
	}

	// Update player's avatar and name
	s.mu.Lock()
	player.Avatar = data.Avatar
	// Set player name to avatar name
	if data.Avatar >= 0 && data.Avatar < len(avatarNames) {
		player.Name = avatarNames[data.Avatar]
	}
	s.mu.Unlock()

	// If player is in a room, notify other players
	if player.RoomID != "" {
		s.mu.RLock()
		room, exists := s.rooms[player.RoomID]
		s.mu.RUnlock()

		if exists {
			// Broadcast player update to room with both avatar and name
			updateData, _ := json.Marshal(map[string]interface{}{
				"player_id": player.ID,
				"avatar":    data.Avatar,
				"name":      player.Name,
			})
			s.broadcastToRoom(room, Message{
				Type:      MsgPlayerUpdate,
				PlayerID:  player.ID,
				Data:      updateData,
				Timestamp: time.Now(),
			})
		}
	}
}

func (s *Server) removePlayerFromRoom(player *Player) {
	if player.RoomID == "" {
		return
	}

	room, exists := s.rooms[player.RoomID]
	if !exists {
		player.RoomID = ""
		return
	}

	room.mu.Lock()
	log.Printf("Before removal: Room %s has %d players\n", room.ID, len(room.Players))

	// Remove player from room
	newPlayers := make([]*Player, 0)
	for _, p := range room.Players {
		if p.ID != player.ID {
			newPlayers = append(newPlayers, p)
		} else {
			log.Printf("Removing player %s from room %s\n", p.ID, room.ID)
		}
	}
	room.Players = newPlayers

	log.Printf("After removal: Room %s has %d players\n", room.ID, len(room.Players))

	// Check if game was in progress
	wasStarted := room.Started

	// Reset room state if someone left during game
	if room.Started && len(room.Players) < room.MaxPlayers {
		room.Started = false
		log.Printf("Room %s reset to not-started (player left during game)\n", room.ID)
	}

	// If room is empty, delete it
	isEmpty := len(room.Players) == 0
	roomID := room.ID
	room.mu.Unlock()

	if isEmpty {
		delete(s.rooms, roomID)
		log.Printf("Room %s deleted (empty)\n", roomID)
	} else {
		// If game was started, end it and kick everyone out
		if wasStarted {
			log.Printf("Game was in progress, ending room %s and removing all remaining players\n", roomID)
			s.broadcastToRoom(room, Message{
				Type:      "game_ended",
				PlayerID:  player.ID,
				RoomID:    roomID,
				Timestamp: time.Now(),
			})

			// Remove all remaining players from the room
			room.mu.Lock()
			for _, p := range room.Players {
				p.RoomID = ""
				log.Printf("Removed player %s from ended room %s\n", p.ID, roomID)
			}
			room.Players = make([]*Player, 0)
			room.mu.Unlock()

			// Delete the room since the game ended
			delete(s.rooms, roomID)
			log.Printf("Room %s deleted (game ended)\n", roomID)
		} else {
			// Game hadn't started yet, just notify remaining players
			s.broadcastToRoom(room, Message{
				Type:      "player_left",
				PlayerID:  player.ID,
				RoomID:    roomID,
				Timestamp: time.Now(),
			})
		}
	}

	player.RoomID = ""
}

func (s *Server) sendMessage(player *Player, msg Message) {
	player.mu.Lock()
	defer player.mu.Unlock()

	err := player.Conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Error sending to player %s: %v\n", player.ID, err)
	}
}

func (s *Server) sendError(player *Player, errMsg string) {
	s.sendMessage(player, Message{
		Type:      MsgError,
		Data:      json.RawMessage(`{"error":"` + errMsg + `"}`),
		Timestamp: time.Now(),
	})
}

func (s *Server) broadcastToRoom(room *Room, msg Message) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	for _, player := range room.Players {
		s.sendMessage(player, msg)
	}
}

func (s *Server) broadcastRoomList() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		s.sendRoomList(player)
	}
}

func (s *Server) sendRoomList(player *Player) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type RoomInfo struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		GameType   string `json:"game_type"`
		Players    int    `json:"players"`
		MaxPlayers int    `json:"max_players"`
		Started    bool   `json:"started"`
	}

	rooms := make([]RoomInfo, 0)
	for _, room := range s.rooms {
		room.mu.RLock()
		rooms = append(rooms, RoomInfo{
			ID:         room.ID,
			Name:       room.Name,
			GameType:   room.GameType,
			Players:    len(room.Players),
			MaxPlayers: room.MaxPlayers,
			Started:    room.Started,
		})
		room.mu.RUnlock()
	}

	data, _ := json.Marshal(map[string]interface{}{
		"rooms": rooms,
	})

	s.sendMessage(player, Message{
		Type:      MsgRoomList,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

func main() {
	server := NewServer()

	http.HandleFunc("/ws", server.handleConnection)

	// Simple root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Game server is running. Use a game client to connect."))
	})

	// Use PORT env variable if available (for Render), otherwise 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/givensuman/teletyperacer/server/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Message represents a WebSocket message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// Room represents a game room
type Room struct {
	clients   map[string]*websocket.Conn
	indices   map[string]int // clientID -> playerIndex
	nextIndex int
}

// RoomManager manages WebSocket connections and rooms
type RoomManager struct {
	rooms map[string]*Room // roomCode -> room
	mu    sync.RWMutex
}

// NewRoomManager creates a new room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// AddClient adds a client to a room
func (rm *RoomManager) AddClient(roomCode, clientID string, conn *websocket.Conn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.rooms[roomCode] == nil {
		rm.rooms[roomCode] = &Room{
			clients:   make(map[string]*websocket.Conn),
			indices:   make(map[string]int),
			nextIndex: 0,
		}
	}
	room := rm.rooms[roomCode]
	if _, exists := room.indices[clientID]; !exists {
		room.indices[clientID] = room.nextIndex
		room.nextIndex++
	}
	room.clients[clientID] = conn
}

// RemoveClient removes a client from a room
func (rm *RoomManager) RemoveClient(roomCode, clientID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[roomCode]; exists {
		delete(room.clients, clientID)
		delete(room.indices, clientID)
		if len(room.clients) == 0 {
			delete(rm.rooms, roomCode)
		}
	}
}

// BroadcastToRoom broadcasts a message to all clients in a room except the sender
func (rm *RoomManager) BroadcastToRoom(roomCode, senderID string, msg interface{}) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomCode]
	if !exists {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for clientID, conn := range room.clients {
		if clientID != senderID {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error broadcasting to client %s: %v", clientID, err)
			}
		}
	}
}

// GetRoomClients returns the number of clients in a room
func (rm *RoomManager) GetRoomClients(roomCode string) int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if room, exists := rm.rooms[roomCode]; exists {
		return len(room.clients)
	}
	return 0
}

// GetPlayerIndex returns the player index for a client in a room
func (rm *RoomManager) GetPlayerIndex(roomCode, clientID string) int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if room, exists := rm.rooms[roomCode]; exists {
		if index, exists := room.indices[clientID]; exists {
			return index
		}
	}
	return -1
}

var roomManager = NewRoomManager()

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	clientID := uuid.New().String()
	log.Printf("🔌 New WebSocket connection established - Client ID: %s", clientID)

	// Handle messages from this client
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for client %s: %v", clientID, err)
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Error unmarshaling message from client %s: %v", clientID, err)
			continue
		}

		switch msg.Type {
		case "createRoom":
			var req types.CreateRoomRequest
			if dataBytes, err := json.Marshal(msg.Data); err == nil {
				json.Unmarshal(dataBytes, &req)
			}
			handleCreateRoom(conn, clientID, req.Code)

		case "joinRoom":
			var req types.JoinRoomRequest
			if dataBytes, err := json.Marshal(msg.Data); err == nil {
				json.Unmarshal(dataBytes, &req)
			}
			handleJoinRoom(conn, clientID, req.Code)

		case "getRoomState":
			var req types.JoinRoomRequest // reuse for code
			if dataBytes, err := json.Marshal(msg.Data); err == nil {
				json.Unmarshal(dataBytes, &req)
			}
			handleGetRoomState(conn, clientID, req.Code)

		default:
			log.Printf("Unknown message type from client %s: %s", clientID, msg.Type)
		}
	}

	// Clean up when client disconnects
	log.Printf("🔌 WebSocket connection closed - Client ID: %s disconnected", clientID)
	// Note: In a real implementation, you'd need to track which rooms the client was in
	// For simplicity, we'll assume clients are only in one room at a time
}

func handleCreateRoom(conn *websocket.Conn, clientID, code string) {
	log.Printf("🏠 Client %s attempting to create room with code %s", clientID, code)

	roomManager.AddClient(code, clientID, conn)
	log.Printf("✅ Room %s created successfully by client %s", code, clientID)

	// Send room created confirmation
	response := Message{Type: "roomCreated", Data: types.RoomCreatedResponse{Code: code}}
	sendMessage(conn, response)
	log.Printf("📤 Sent roomCreated confirmation to client %s for room %s", clientID, code)

	// Send initial room state (just the host for now)
	playerCount := roomManager.GetRoomClients(code)
	yourIndex := roomManager.GetPlayerIndex(code, clientID)
	roomState := types.RoomStateResponse{
		Code:        code,
		PlayerCount: playerCount,
		YourIndex:   yourIndex,
	}
	stateMsg := Message{Type: "roomState", Data: roomState}
	sendMessage(conn, stateMsg)
	log.Printf("📤 Sent initial roomState to client %s for room %s: %d players, yourIndex %d", clientID, code, roomState.PlayerCount, roomState.YourIndex)
}

func handleJoinRoom(conn *websocket.Conn, clientID, code string) {
	log.Printf("🚪 Client %s attempting to join room %s", clientID, code)

	roomManager.AddClient(code, clientID, conn)
	log.Printf("✅ Client %s successfully joined room %s", clientID, code)

	// Send join confirmation
	response := Message{Type: "roomJoined", Data: types.RoomJoinedResponse{Code: code}}
	sendMessage(conn, response)
	log.Printf("📤 Sent roomJoined confirmation to client %s for room %s", clientID, code)

	// Send room state to the new player
	playerCount := roomManager.GetRoomClients(code)
	yourIndex := roomManager.GetPlayerIndex(code, clientID)
	roomState := types.RoomStateResponse{
		Code:        code,
		PlayerCount: playerCount,
		YourIndex:   yourIndex,
	}
	stateMsg := Message{Type: "roomState", Data: roomState}
	sendMessage(conn, stateMsg)
	log.Printf("📤 Sent roomState to client %s for room %s: %d players, yourIndex %d", clientID, code, roomState.PlayerCount, roomState.YourIndex)

	// Send updated room state to all existing players
	roomManager.mu.RLock()
	room := roomManager.rooms[code]
	for existingClientID, existingConn := range room.clients {
		if existingClientID != clientID {
			existingIndex := room.indices[existingClientID]
			existingRoomState := types.RoomStateResponse{
				Code:        code,
				PlayerCount: playerCount,
				YourIndex:   existingIndex,
			}
			existingStateMsg := Message{Type: "roomState", Data: existingRoomState}
			sendMessage(existingConn, existingStateMsg)
			log.Printf("📤 Sent updated roomState to existing client %s for room %s: %d players, yourIndex %d", existingClientID, code, existingRoomState.PlayerCount, existingRoomState.YourIndex)
		}
	}
	roomManager.mu.RUnlock()

	// Broadcast player joined event
	broadcastMsg := Message{Type: "playerJoined", Data: types.PlayerJoinedResponse{PlayerIndex: yourIndex}}
	roomManager.BroadcastToRoom(code, clientID, broadcastMsg)
	log.Printf("📢 Broadcasted playerJoined event to room %s (excluding sender %s)", code, clientID)
}

func handleGetRoomState(conn *websocket.Conn, clientID, code string) {
	log.Printf("📥 Client %s requesting room state for room %s", clientID, code)

	playerCount := roomManager.GetRoomClients(code)
	yourIndex := roomManager.GetPlayerIndex(code, clientID)
	if yourIndex == -1 {
		log.Printf("Client %s not in room %s", clientID, code)
		return
	}

	roomState := types.RoomStateResponse{
		Code:        code,
		PlayerCount: playerCount,
		YourIndex:   yourIndex,
	}
	stateMsg := Message{Type: "roomState", Data: roomState}
	sendMessage(conn, stateMsg)
	log.Printf("📤 Sent roomState to client %s for room %s: %d players, yourIndex %d", clientID, code, roomState.PlayerCount, roomState.YourIndex)
}

func sendMessage(conn *websocket.Conn, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

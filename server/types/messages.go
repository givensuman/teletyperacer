package types

// Message types for WebSocket communication
type CreateRoomRequest struct {
	Code string `json:"code"`
}

type JoinRoomRequest struct {
	Code string `json:"code"`
}

type RoomCreatedResponse struct {
	Code string `json:"code"`
}

type RoomJoinedResponse struct {
	Code string `json:"code"`
}

type RoomStateResponse struct {
	Code    string   `json:"code"`
	Players []string `json:"players"`
}

type PlayerJoinedResponse struct {
	PlayerName string `json:"playerName"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

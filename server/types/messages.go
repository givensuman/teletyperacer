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
	Code        string `json:"code"`
	PlayerCount int    `json:"playerCount"`
	YourIndex   int    `json:"yourIndex"`
}

type PlayerJoinedResponse struct {
	PlayerIndex int `json:"playerIndex"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

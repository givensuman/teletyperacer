// Package types contains shared types
// for the TUI
package types

type Screen int

const (
	HomeScreen Screen = iota
	HostScreen
	PracticeScreen
	RoomInputScreen
)

type ScreenChangeMsg struct {
	Screen Screen
}

type ConnectionStatus int

const (
	Connecting ConnectionStatus = iota
	Connected
	ServerUnreachable
	ClientError
	Failed // Keep for backward compatibility
)

type ConnectionStatusMsg struct {
	Status ConnectionStatus
}

// Room-related messages
type CreateRoomMsg struct{}

type JoinRoomMsg struct {
	Code string
}

type RoomCreatedMsg struct {
	Code string
}

type RoomJoinedMsg struct {
	Code string
}

type PlayerJoinedMsg struct {
	PlayerName string
}

type RoomStateMsg struct {
	Code    string
	Players []string
}

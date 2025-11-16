package types

// View defines the application views which may
// be rendered at any given time.
type View int64

const (
	None View = iota
	Home
	Practice
	Host
	Join
	Lobby
	Play
)

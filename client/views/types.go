package views

// View defines the application views which may
// be rendered at any given time.
type View int64

const (
	None View = iota // Identity case
	Home
	Practice
	Host
	Join
	Lobby
	Play
)

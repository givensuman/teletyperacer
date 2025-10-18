module github.com/givensuman/teletyperacer/server

go 1.25.2

require (
	github.com/givensuman/teletyperacer/sockets v0.0.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
)

replace github.com/givensuman/teletyperacer/sockets => ../sockets

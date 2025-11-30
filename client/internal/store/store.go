// Package store provides a global singleton store for application state
package store

import (
	"github.com/givensuman/teletyperacer/client/internal/tui"
	"sync"
)

type Store struct {
	ConnectionStatus tui.ConnectionStatus
	Width            int
	Height           int
}

var instance *Store
var once sync.Once

// GetStore returns the singleton store instance
func GetStore() *Store {
	once.Do(func() {
		instance = &Store{
			ConnectionStatus: tui.Connecting,
			Width:            80,
			Height:           24,
		}
	})

	return instance
}

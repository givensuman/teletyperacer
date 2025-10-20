package sockets

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHub(t *testing.T) {
	hub := CreateHub()

	t.Run("AddRoom", func(t *testing.T) {
		room := hub.CreateRoom("test room")
		assert.NotNil(t, room)
		assert.Equal(t, "test room", room.Name)
		assert.False(t, room.IsPrivate)

		// Check room is in hub
		hubRoom, err := hub.GetRoom(room.ID.String())
		require.NoError(t, err)
		assert.Equal(t, room, hubRoom)
	})

	t.Run("AddClient", func(t *testing.T) {
		client := CreateClient(nil)
		hub.AddClient(client)

		hubClient, err := hub.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Equal(t, client, hubClient)
	})

	t.Run("MoveClientToRoom", func(t *testing.T) {
		client := CreateClient(nil)
		hub.AddClient(client)

		room := hub.CreateRoom("test room")

		err := hub.MoveClientToRoom(client.ID.String(), room.ID.String())
		require.NoError(t, err)

		// Check client is in room
		roomClient, err := room.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Equal(t, client, roomClient)

		// Check client not in hub
		hubClient, err := hub.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Nil(t, hubClient)
	})

	t.Run("MoveClientOutOfRoom", func(t *testing.T) {
		room := hub.CreateRoom("test room")

		client := CreateClient(nil)
		room.AddClient(client)

		err := hub.MoveClientOutOfRoom(client.ID.String(), room.ID.String())
		require.NoError(t, err)

		// Check client is back in hub
		hubClient, err := hub.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Equal(t, client, hubClient)

		// Check client not in room
		roomClient, err := room.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Nil(t, roomClient)
	})
}

func TestRoom(t *testing.T) {
	hub := CreateHub()
	room := hub.CreateRoom("test room")

	t.Run("AddClient", func(t *testing.T) {
		client := CreateClient(nil)
		room.AddClient(client)

		roomClient, err := room.GetClient(client.ID.String())
		require.NoError(t, err)
		assert.Equal(t, client, roomClient)
		assert.Equal(t, room, client.room)
	})

	t.Run("SendToAll", func(t *testing.T) {
		msg := &Message{
			Event: "test",
			Data:  []byte("test data"),
		}
		err := room.SendToAll(msg)
		// Should not error even with no clients
		assert.NoError(t, err)
	})

	t.Run("EventHandlers", func(t *testing.T) {
		called := false
		var receivedData any
		var receivedClient *Client

		room.On("test_event", func(c *Client, data any) {
			called = true
			receivedClient = c
			receivedData = data
		})

		client := CreateClient(nil)
		room.AddClient(client)

		// Simulate receiving a message
		room.eventHandlers["test event"][0](client, []byte("test data"))

		assert.True(t, called)
		assert.Equal(t, client, receivedClient)
		assert.Equal(t, []byte("test data"), receivedData)
	})
}

func TestClient(t *testing.T) {
	t.Run("CreateClient", func(t *testing.T) {
		client := CreateClient(nil)
		assert.NotNil(t, client)
		assert.NotEqual(t, uuid.Nil, client.ID)
		assert.Nil(t, client.room)
		assert.NotNil(t, client.Callbacks)
	})

	t.Run("SetRoom", func(t *testing.T) {
		client := CreateClient(nil)
		room := &Room{}
		client.SetRoom(room)
		assert.Equal(t, room, client.room)
	})

	t.Run("AddCallback", func(t *testing.T) {
		client := CreateClient(nil)
		called := false
		client.AddCallback("test", func() {
			called = true
		})

		client.ReceiveMessage(&Message{Event: "test"})
		assert.True(t, called)
	})
}

func TestMessage(t *testing.T) {
	t.Run("MessageStruct", func(t *testing.T) {
		id := uuid.New()
		msg := &Message{
			SenderID:   id,
			Event:      "test",
			Data:       []byte("data"),
			CallbackID: &id,
		}
		assert.Equal(t, id, msg.SenderID)
		assert.Equal(t, "test", msg.Event)
		assert.Equal(t, []byte("data"), msg.Data)
		assert.Equal(t, &id, msg.CallbackID)
	})
}

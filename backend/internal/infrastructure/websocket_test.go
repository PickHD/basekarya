package infrastructure

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// MockClient creates a client with a buffered Send channel for testing
func NewMockClient(hub *Hub, userID uint) *Client {
	return &Client{
		Hub:    hub,
		UserID: userID,
		Send:   make(chan []byte, 1), // Buffer of 1 for testing
		Conn:   &websocket.Conn{},    // Dummy conn, won't be used for write in this test
	}
}

func setupRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return rdb, mr
}

func TestHub_Broadcast_PanicOnSlowClient(t *testing.T) {
	rdb, mr := setupRedis(t)
	defer mr.Close()

	hub := NewHub(rdb)
	go hub.Run()

	userID := uint(1)
	client := NewMockClient(hub, userID)

	// Register client
	hub.Register <- client
	time.Sleep(10 * time.Millisecond) // Wait for registration

	// Fill the client's Send channel
	client.Send <- []byte("fill buffer")

	// Now broadcast to this user.
	// The client's Send channel is full.
	// OLD CODE: would close(client.Send).
	// NEW CODE: should trigger Unregister in a goroutine.
	hub.SendToUser(userID, []byte("message 1"))

	// Give time for the unregister goroutine to run
	time.Sleep(50 * time.Millisecond)

	// Verify client is unregistered
	hub.mu.RLock()
	clients := hub.Clients[userID]
	hub.mu.RUnlock()

	// Logic check:
	// If unregister happened, the client should be removed from the slice (or slice empty)
	found := false
	for _, c := range clients {
		if c == client {
			found = true
			break
		}
	}

	if found {
		t.Errorf("Client should have been unregistered due to full buffer")
	}

	// Double check: Try to close the channel again manually to see if it was already closed?
	// Actually, the Hub's Unregister loop closes the channel.
	// If the Panic bug existed (double close), it might have panicked the server during the Sleep.
	// If the test is still running here without panic, that's a good sign.
}

func TestHub_Broadcast_Success(t *testing.T) {
	rdb, mr := setupRedis(t)
	defer mr.Close()

	hub := NewHub(rdb)
	go hub.Run()

	userID := uint(2)
	client := NewMockClient(hub, userID)

	// Register
	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	// Subscribe to the channel to verify message publication (optional but good)
	// We can trust SendToUser publishes to Redis, and Run receives from Redis.

	// Broadcast
	msg := []byte("hello")
	hub.SendToUser(userID, msg)

	// Need to give a bit more time for Redis Pub/Sub roundtrip
	// 1. SendToUser -> Redis
	// 2. Redis -> Hub.Run (via subscribe)
	// 3. Hub.Run -> Client.Send

	// Check if message received
	select {
	case received := <-client.Send:
		var receivedMsg string
		// The received message is byte array.
		receivedMsg = string(received)

		if receivedMsg != string(msg) {
			t.Errorf("Expected %s, got %s", msg, receivedMsg)
		}
	case <-time.After(500 * time.Millisecond): // Increased timeout for roundtrip
		t.Errorf("Timeout waiting for message")
	}
}

func TestHub_RedisBroadcast_MultiInstance(t *testing.T) {
	// Simulate two Hubs sharing the same Redis
	rdb, mr := setupRedis(t)
	defer mr.Close()

	// Hub 1
	hub1 := NewHub(rdb)
	go hub1.Run()

	// Hub 2
	hub2 := NewHub(rdb)
	go hub2.Run()

	userID := uint(3)
	// Client connected to Hub 1
	client1 := NewMockClient(hub1, userID)
	hub1.Register <- client1

	time.Sleep(10 * time.Millisecond)

	// Send message from Hub 2
	msg := []byte("msg from hub 2")
	hub2.SendToUser(userID, msg)

	// Client on Hub 1 should receive it via Redis
	select {
	case received := <-client1.Send:
		if string(received) != string(msg) {
			t.Errorf("Expected %s, got %s", msg, received)
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Timeout waiting for message in multi-instance test")
	}
}

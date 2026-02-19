package infrastructure

import (
	"context"
	"encoding/json"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Hub    *Hub
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients     map[uint][]*Client
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Message
	mu          sync.RWMutex
	RedisClient *redis.Client
}

type Message struct {
	TargetUserID uint
	Data         []byte
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		Clients:     make(map[uint][]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan Message, 256),
		RedisClient: redisClient,
	}
}

func (h *Hub) Run() {
	logger.Info("Websocket Started...")

	ctx := context.Background()
	pubsub := h.RedisClient.Subscribe(ctx, constants.RedisBroadcastChannel)
	defer pubsub.Close()

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err == nil {
				h.Broadcast <- message
			} else {
				logger.Errorw("[WS] Failed to unmarshal redis message", err)
			}
		}
	}()

	for {
		select {
		case client := <-h.Register:
			logger.Infof("[WS] Registering client for user %d", client.UserID)
			h.mu.Lock()
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
			h.mu.Unlock()
		case client := <-h.Unregister:
			logger.Infof("[WS] Unregistering client for user %d", client.UserID)
			h.mu.Lock()

			if clients, ok := h.Clients[client.UserID]; ok {
				for i, c := range clients {
					if c == client {
						h.Clients[client.UserID] = append(clients[:i], clients[i+1:]...)
						close(c.Send)
						break
					}
				}

				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}

			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.RLock()
			clients, ok := h.Clients[message.TargetUserID]
			h.mu.RUnlock()

			if ok {
				for _, client := range clients {
					select {
					case client.Send <- message.Data:
					default:
						go func(c *Client) {
							h.Unregister <- c
						}(client)
					}
				}
			}
		}
	}
}

func (h *Hub) SendToUser(userID uint, payload []byte) {
	msg := Message{TargetUserID: userID, Data: payload}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Errorw("[WS] Failed to marshal message", err)
		return
	}

	h.RedisClient.Publish(context.Background(), constants.RedisBroadcastChannel, data)
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				logger.Errorw("[WS] Websocket error: ", err)
			}
			break
		}

	}
}

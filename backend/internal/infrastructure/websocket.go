package infrastructure

import (
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients    map[uint][]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.RWMutex
}

type Message struct {
	TargetUserID uint
	Data         []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	logger.Info("Websocket Started...")

	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			// add new client to slices
			h.Clients[client.UserID] = append(h.Clients[client.UserID], client)
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()

			if clients, ok := h.Clients[client.UserID]; ok {
				for i, c := range clients {
					if c == client {
						// remove client from slices safely
						h.Clients[client.UserID] = append(clients[:i], clients[i+1:]...)
						close(c.Send)
						break
					}
				}

				// if user already have no connections left, remove the map key
				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}

			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.RLock()
			// since only read maps, using RLock
			clients, ok := h.Clients[message.TargetUserID]
			h.mu.RUnlock()

			if ok {
				for _, client := range clients {
					select {
					case client.Send <- message.Data:
					default:
						close(client.Send)
					}
				}
			}
		}
	}
}

func (h *Hub) SendToUser(userID uint, payload []byte) {
	h.Broadcast <- Message{TargetUserID: userID, Data: payload}
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

			// Send batching messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// send ping with defined interval
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

	// read configuration
	c.Conn.SetReadLimit(512) // limit message incoming due security purposes
	c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorw("[WS] Websocket error: ", err)
			}
			break
		}

	}
}

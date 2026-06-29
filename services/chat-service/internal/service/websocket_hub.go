package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	RoomID   uuid.UUID
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *WebSocketHub
}

type WebSocketHub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*Client]bool
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		rooms: make(map[string]map[*Client]bool),
	}
}

func (h *WebSocketHub) JoinRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.RoomID.String()] == nil {
		h.rooms[client.RoomID.String()] = make(map[*Client]bool)
	}
	h.rooms[client.RoomID.String()][client] = true

	log.Info().Str("user_id", client.UserID.String()).Str("room_id", client.RoomID.String()).Msg("client joined room")
}

func (h *WebSocketHub) LeaveRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.RoomID.String()]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.rooms, client.RoomID.String())
		}
	}
	close(client.Send)
	log.Info().Str("user_id", client.UserID.String()).Str("room_id", client.RoomID.String()).Msg("client left room")
}

func (h *WebSocketHub) Broadcast(roomID string, message []byte, sender *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			if client == sender {
				continue
			}
			select {
			case client.Send <- message:
			default:
				h.mu.RUnlock()
				h.LeaveRoom(client)
				h.mu.RLock()
			}
		}
	}
}

func (h *WebSocketHub) GetRoomClients(roomID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	for client := range h.rooms[roomID] {
		clients = append(clients, client)
	}
	return clients
}

func (h *WebSocketHub) GetRoomCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.LeaveRoom(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error().Err(err).Msg("websocket read error")
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		c.Hub.Broadcast(c.RoomID.String(), message, c)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
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
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

package service

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type WebSocketClient struct {
	ID       uuid.UUID
	StreamID uuid.UUID
	Conn     *websocket.Conn
	Send     chan []byte
	mu       sync.Mutex
}

type WebSocketHub struct {
	mu      sync.RWMutex
	rooms   map[uuid.UUID]map[uuid.UUID]*WebSocketClient
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		rooms: make(map[uuid.UUID]map[uuid.UUID]*WebSocketClient),
	}
}

func (h *WebSocketHub) Register(client *WebSocketClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[client.StreamID]; !ok {
		h.rooms[client.StreamID] = make(map[uuid.UUID]*WebSocketClient)
	}
	h.rooms[client.StreamID][client.ID] = client

	log.Debug().
		Str("client_id", client.ID.String()).
		Str("stream_id", client.StreamID.String()).
		Msg("WebSocket client registered")
}

func (h *WebSocketHub) Unregister(client *WebSocketClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.StreamID]; ok {
		if _, ok := clients[client.ID]; ok {
			delete(clients, client.ID)
			close(client.Send)
			if len(clients) == 0 {
				delete(h.rooms, client.StreamID)
			}
		}
	}

	log.Debug().
		Str("client_id", client.ID.String()).
		Str("stream_id", client.StreamID.String()).
		Msg("WebSocket client unregistered")
}

func (h *WebSocketHub) BroadcastToStream(streamID uuid.UUID, msg WebSocketMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal WebSocket message")
		return
	}

	h.mu.RLock()
	clients, ok := h.rooms[streamID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			log.Warn().
				Str("client_id", client.ID.String()).
				Str("stream_id", streamID.String()).
				Msg("WebSocket client send buffer full, dropping message")
		}
	}
}

func (h *WebSocketHub) SendToClient(clientID uuid.UUID, streamID uuid.UUID, msg WebSocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	clients, ok := h.rooms[streamID]
	h.mu.RUnlock()

	if !ok {
		return nil
	}

	client, ok := clients[clientID]
	if !ok {
		return nil
	}

	select {
	case client.Send <- data:
	default:
		return nil
	}

	return nil
}

func (h *WebSocketHub) GetStreamClientCount(streamID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[streamID]
	if !ok {
		return 0
	}
	return len(clients)
}

func (h *WebSocketHub) GetStreamClients(streamID uuid.UUID) []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[streamID]
	if !ok {
		return nil
	}

	ids := make([]uuid.UUID, 0, len(clients))
	for id := range clients {
		ids = append(ids, id)
	}
	return ids
}

func (h *WebSocketHub) RemoveStream(streamID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.rooms[streamID]
	if !ok {
		return
	}

	for _, client := range clients {
		close(client.Send)
	}
	delete(h.rooms, streamID)
}

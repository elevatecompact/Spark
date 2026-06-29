package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/chat-service/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub *service.WebSocketHub
}

func NewWebSocketHandler(hub *service.WebSocketHub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	username, _ := r.Context().Value("username").(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("websocket upgrade failed")
		return
	}

	client := &service.Client{
		ID:       uuid.New(),
		UserID:   userID,
		RoomID:   roomID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      h.hub,
	}

	h.hub.JoinRoom(client)

	go client.WritePump()
	go client.ReadPump()
}

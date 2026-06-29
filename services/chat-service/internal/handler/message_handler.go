package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/domain"
	"github.com/elevatecompact/spark/services/chat-service/internal/service"
)

type MessageHandler struct {
	svc service.MessageService
}

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	username, _ := r.Context().Value("username").(string)

	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	msg, err := h.svc.SendMessage(r.Context(), roomID, userID, username, req)
	if err != nil {
		switch {
		case err == domain.ErrUserBanned:
			respondError(w, http.StatusForbidden, "user is banned from this room")
		case err == domain.ErrUserMuted:
			respondError(w, http.StatusForbidden, "user is muted in this room")
		case err == domain.ErrMessageTooLong:
			respondError(w, http.StatusRequestEntityTooLarge, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, msg)
}

func (h *MessageHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	var cursor time.Time
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		cursor, err = time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid cursor format, use RFC3339")
			return
		}
	} else {
		cursor = time.Now().UTC().Add(time.Hour)
	}

	messages, err := h.svc.GetHistory(r.Context(), roomID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, messages)
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	msg, err := h.svc.EditMessage(r.Context(), msgID, req.Content)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, msg)
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	if err := h.svc.DeleteMessage(r.Context(), msgID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

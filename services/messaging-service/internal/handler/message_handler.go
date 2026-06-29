package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/messaging-service/internal/domain"
	"github.com/elevatecompact/spark/services/messaging-service/internal/service"
)

type MessageHandler struct {
	svc service.MessageService
}

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	convID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid conversation id")
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	msg, err := h.svc.Send(r.Context(), convID, userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, msg)
}

func (h *MessageHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	convID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid conversation id")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var cursor time.Time
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		cursor, _ = time.Parse(time.RFC3339, cursorStr)
	}

	msgs, err := h.svc.GetHistory(r.Context(), convID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, msgs)
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "msgId"))
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

	if err := h.svc.Edit(r.Context(), msgID, req.Content); err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "msgId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	if err := h.svc.Delete(r.Context(), msgID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MessageHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "msgId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.AddReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.AddReaction(r.Context(), msgID, userID, req.Emoji); err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MessageHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	msgID, err := uuid.Parse(chi.URLParam(r, "msgId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid message id")
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	emoji := chi.URLParam(r, "emoji")

	if err := h.svc.RemoveReaction(r.Context(), msgID, userID, emoji); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

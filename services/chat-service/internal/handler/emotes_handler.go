package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/chat-service/internal/service"
)

type EmotesHandler struct {
	svc service.EmoteService
}

func NewEmotesHandler(svc service.EmoteService) *EmotesHandler {
	return &EmotesHandler{svc: svc}
}

func (h *EmotesHandler) GetGlobal(w http.ResponseWriter, r *http.Request) {
	emotes, err := h.svc.GetGlobal(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emotes)
}

func (h *EmotesHandler) GetByRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	emotes, err := h.svc.GetByRoom(r.Context(), roomID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emotes)
}

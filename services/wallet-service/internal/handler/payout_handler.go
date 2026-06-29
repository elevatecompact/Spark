package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/wallet-service/internal/domain"
	"github.com/elevatecompact/spark/services/wallet-service/internal/service"
)

type PayoutHandler struct {
	svc service.PayoutService
}

func NewPayoutHandler(svc service.PayoutService) *PayoutHandler {
	return &PayoutHandler{svc: svc}
}

func (h *PayoutHandler) Request(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payout, err := h.svc.Request(r.Context(), userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, payout)
}

func (h *PayoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid payout id")
		return
	}

	payout, err := h.svc.Get(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "payout not found")
		return
	}

	respondJSON(w, http.StatusOK, payout)
}

func (h *PayoutHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
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

	payouts, err := h.svc.ListByUser(r.Context(), userID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, payouts)
}

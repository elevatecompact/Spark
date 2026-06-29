package handler

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/wallet-service/internal/service"
)

type WalletHandler struct {
	svc service.WalletService
}

func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) GetMyWallet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wallet, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "wallet not found")
		return
	}

	respondJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) GetMyBalances(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wallet, err := h.svc.GetOrCreate(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	wallet, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "wallet not found")
		return
	}

	respondJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) Freeze(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wallet, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "wallet not found")
		return
	}

	if err := h.svc.Freeze(r.Context(), wallet.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WalletHandler) Close(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wallet, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "wallet not found")
		return
	}

	if err := h.svc.Close(r.Context(), wallet.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

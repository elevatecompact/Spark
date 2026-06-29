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

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Type = domain.TxnDeposit

	txn, err := h.svc.Deposit(r.Context(), userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Type = domain.TxnWithdraw

	txn, err := h.svc.Withdraw(r.Context(), userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Type = domain.TxnTransfer

	txn, err := h.svc.Transfer(r.Context(), userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) Tip(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Type = domain.TxnTip

	txn, err := h.svc.Tip(r.Context(), userID, req)
	if err != nil {
		respondError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	txn, err := h.svc.Get(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "transaction not found")
		return
	}

	respondJSON(w, http.StatusOK, txn)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
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

	txns, err := h.svc.ListByUser(r.Context(), userID, cursor, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, txns)
}

package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
	"github.com/elevatecompact/spark/services/payment-service/internal/service"
)

type PaymentHandler struct {
	svc service.PaymentService
}

func NewPaymentHandler(svc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// Intents
func (h *PaymentHandler) CreateIntent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.CreateIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if key := r.Header.Get("Idempotency-Key"); key != "" {
		req.IdempotencyKey = key
	}
	intent, err := h.svc.CreateIntent(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, intent)
}

func (h *PaymentHandler) GetIntent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid intent id")
		return
	}
	intent, err := h.svc.GetIntent(r.Context(), id)
	if err != nil {
		if err == domain.ErrIntentNotFound {
			writeError(w, http.StatusNotFound, "intent not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, intent)
}

func (h *PaymentHandler) ListIntents(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	cursor := parseCursor(r.URL.Query().Get("cursor"))
	limit := parseInt(r.URL.Query().Get("limit"), 50)
	intents, err := h.svc.ListIntents(r.Context(), userID, cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, intents)
}

func (h *PaymentHandler) ConfirmIntent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid intent id")
		return
	}
	var req domain.ConfirmIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	intent, err := h.svc.ConfirmIntent(r.Context(), id, userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, intent)
}

func (h *PaymentHandler) CancelIntent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid intent id")
		return
	}
	if err := h.svc.CancelIntent(r.Context(), id, userID); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "canceled"})
}

// Methods
func (h *PaymentHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.CreatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	method, err := h.svc.CreatePaymentMethod(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, method)
}

func (h *PaymentHandler) GetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid method id")
		return
	}
	method, err := h.svc.GetPaymentMethod(r.Context(), id)
	if err != nil {
		if err == domain.ErrMethodNotFound {
			writeError(w, http.StatusNotFound, "method not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, method)
}

func (h *PaymentHandler) ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	methods, err := h.svc.ListPaymentMethods(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, methods)
}

func (h *PaymentHandler) SetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid method id")
		return
	}
	if err := h.svc.SetDefaultPaymentMethod(r.Context(), id, userID); err != nil {
		if err == domain.ErrMethodNotFound {
			writeError(w, http.StatusNotFound, "method not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "default set"})
}

func (h *PaymentHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid method id")
		return
	}
	if err := h.svc.DeletePaymentMethod(r.Context(), id, userID); err != nil {
		if err == domain.ErrMethodNotFound {
			writeError(w, http.StatusNotFound, "method not found")
			return
		}
		if err == domain.ErrForbidden {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Refunds
func (h *PaymentHandler) RefundIntent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid intent id")
		return
	}
	var req domain.RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.RefundIntent(r.Context(), id, userID, req); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "refunded"})
}

func (h *PaymentHandler) ListRefunds(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, []interface{}{})
}

// Payouts
func (h *PaymentHandler) CreatePayout(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	payout, err := h.svc.CreatePayout(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, payout)
}

func (h *PaymentHandler) GetPayout(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid payout id")
		return
	}
	payout, err := h.svc.GetPayout(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			writeError(w, http.StatusNotFound, "payout not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, payout)
}

// Webhooks
func (h *PaymentHandler) ProcessStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}
	var payload struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		payload.ID = "evt_noop_" + uuid.New().String()
		payload.Type = "unknown"
	}
	if err := h.svc.ProcessWebhook(r.Context(), domain.ProcessorStripe, payload.ID, payload.Type, body); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "received"})
}

func (h *PaymentHandler) ProcessPayPalWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}
	var payload struct {
		ID   string `json:"id"`
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		payload.ID = "evt_noop_" + uuid.New().String()
		payload.EventType = "UNKNOWN"
	}
	if err := h.svc.ProcessWebhook(r.Context(), domain.ProcessorPayPal, payload.ID, payload.EventType, body); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "received"})
}

// Admin
func (h *PaymentHandler) GetProcessorStatus(w http.ResponseWriter, r *http.Request) {
	status := h.svc.GetProcessorStatus(r.Context())
	writeJSON(w, http.StatusOK, status)
}

func (h *PaymentHandler) RetryWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid webhook id")
		return
	}
	if err := h.svc.RetryWebhook(r.Context(), id); err != nil {
		if err == domain.ErrWebhookNotFound {
			writeError(w, http.StatusNotFound, "webhook not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "retried"})
}

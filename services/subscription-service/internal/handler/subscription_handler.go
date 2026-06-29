package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
	"github.com/elevatecompact/spark/services/subscription-service/internal/service"
)

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscriptionHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req domain.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sub, err := h.svc.Subscribe(r.Context(), userID, req)
	if err != nil {
		if err == domain.ErrPlanNotFound || err == domain.ErrPlanInactive {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == domain.ErrAlreadySubscribed || err == domain.ErrMaxSubscriptions {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, sub)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "subId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}
	sub, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err == domain.ErrSubscriptionNotFound {
			writeError(w, http.StatusNotFound, "subscription not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	subs, err := h.svc.GetMy(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "subId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}
	userID := getUserID(r)
	if err := h.svc.Cancel(r.Context(), id, userID); err != nil {
		if err == domain.ErrSubscriptionNotFound {
			writeError(w, http.StatusNotFound, "subscription not found")
			return
		}
		if err == domain.ErrNotOwner {
			writeError(w, http.StatusForbidden, "not owner")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Reactivate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "subId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}
	userID := getUserID(r)
	if err := h.svc.Reactivate(r.Context(), id, userID); err != nil {
		if err == domain.ErrSubscriptionNotFound {
			writeError(w, http.StatusNotFound, "subscription not found")
			return
		}
		if err == domain.ErrNotOwner {
			writeError(w, http.StatusForbidden, "not owner")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reactivated"})
}

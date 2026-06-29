package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
	"github.com/elevatecompact/spark/services/gift-service/internal/service"
)

type GiftHandler struct {
	svc service.GiftService
}

func NewGiftHandler(svc service.GiftService) *GiftHandler {
	return &GiftHandler{svc: svc}
}

// Catalog
func (h *GiftHandler) CreateGiftItem(w http.ResponseWriter, r *http.Request) {
	var item domain.GiftItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.CreateGiftItem(r.Context(), &item); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *GiftHandler) GetGiftItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetGiftItem(r.Context(), id)
	if err != nil {
		if err == domain.ErrGiftItemNotFound {
			writeError(w, http.StatusNotFound, "gift item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *GiftHandler) ListGiftItems(w http.ResponseWriter, r *http.Request) {
	admin := r.URL.Query().Get("admin") == "true"
	items, err := h.svc.ListGiftItems(r.Context(), admin)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *GiftHandler) UpdateGiftItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var item domain.GiftItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item.ID = id
	if err := h.svc.UpdateGiftItem(r.Context(), &item); err != nil {
		if err == domain.ErrGiftItemNotFound {
			writeError(w, http.StatusNotFound, "gift item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *GiftHandler) DeleteGiftItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteGiftItem(r.Context(), id); err != nil {
		if err == domain.ErrGiftItemNotFound {
			writeError(w, http.StatusNotFound, "gift item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Sending
func (h *GiftHandler) SendGift(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.SendGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	gift, err := h.svc.SendGift(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, gift)
}

func (h *GiftHandler) SendBatchGift(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.SendBatchGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	gifts, err := h.svc.SendBatchGift(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, gifts)
}

func (h *GiftHandler) SendSubscriptionGift(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.SendSubscriptionGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	gift, err := h.svc.SendSubscriptionGift(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, gift)
}

// My Gifts
func (h *GiftHandler) GetGift(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	gift, err := h.svc.GetGift(r.Context(), id)
	if err != nil {
		if err == domain.ErrGiftNotFound {
			writeError(w, http.StatusNotFound, "gift not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gift)
}

func (h *GiftHandler) ListSent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	cursor := parseCursor(r.URL.Query().Get("cursor"))
	limit := parseInt(r.URL.Query().Get("limit"), 50)
	gifts, err := h.svc.ListSent(r.Context(), userID, cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gifts)
}

func (h *GiftHandler) ListReceived(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	cursor := parseCursor(r.URL.Query().Get("cursor"))
	limit := parseInt(r.URL.Query().Get("limit"), 50)
	gifts, err := h.svc.ListReceived(r.Context(), userID, cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gifts)
}

// Gift Cards
func (h *GiftHandler) PurchaseGiftCard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.PurchaseGiftCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	card, err := h.svc.PurchaseGiftCard(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, card)
}

func (h *GiftHandler) RedeemGiftCard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.RedeemGiftCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	card, err := h.svc.RedeemGiftCard(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *GiftHandler) GetGiftCardByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	card, err := h.svc.GetGiftCardByCode(r.Context(), code)
	if err != nil {
		if err == domain.ErrGiftCardNotFound {
			writeError(w, http.StatusNotFound, "gift card not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

// Campaigns
func (h *GiftHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	campaign, err := h.svc.CreateCampaign(r.Context(), userID, req)
	if err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, campaign)
}

func (h *GiftHandler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	campaigns, err := h.svc.ListCampaigns(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, campaigns)
}

func (h *GiftHandler) ApplyCampaignMatch(w http.ResponseWriter, r *http.Request) {
	campaignID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign id")
		return
	}
	var req struct {
		GiftID uuid.UUID `json:"gift_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.ApplyCampaignMatch(r.Context(), req.GiftID, campaignID); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "matched"})
}

// Analytics
func (h *GiftHandler) GetTopGifts(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	limit := parseInt(r.URL.Query().Get("limit"), 50)
	gifts, err := h.svc.GetTopGifts(r.Context(), period, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, gifts)
}

func (h *GiftHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	limit := parseInt(r.URL.Query().Get("limit"), 50)
	entries, err := h.svc.GetLeaderboard(r.Context(), period, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// Admin
func (h *GiftHandler) RefundGift(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.RefundGift(r.Context(), id); err != nil {
		if err == domain.ErrGiftNotFound {
			writeError(w, http.StatusNotFound, "gift not found")
			return
		}
		if err == domain.ErrGiftNotCompleted {
			writeError(w, http.StatusBadRequest, "gift is not completed")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "refunded"})
}

func parseCursor(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	var t time.Time
	if err := t.UnmarshalText([]byte(s)); err != nil {
		return time.Time{}
	}
	return t
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/advertising-service/internal/domain"
	"github.com/elevatecompact/spark/services/advertising-service/internal/service"
)

type Handler struct {
	svc service.AdvertisingService
}

func New(svc service.AdvertisingService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/campaigns", func(r chi.Router) {
		r.Post("/", h.createCampaign)
		r.Get("/", h.listCampaigns)
		r.Get("/{id}", h.getCampaign)
		r.Patch("/{id}", h.updateCampaign)
		r.Delete("/{id}", h.deleteCampaign)
	})
	r.Route("/v1/ad-units", func(r chi.Router) {
		r.Post("/", h.createAdUnit)
		r.Get("/{id}", h.getAdUnit)
		r.Patch("/{id}", h.updateAdUnit)
		r.Delete("/{id}", h.deleteAdUnit)
		r.Post("/{id}/approve", h.approveAdUnit)
	})
	r.Get("/v1/ads/request", h.requestAd)
	r.Post("/v1/ads/impression", h.recordImpression)
	r.Post("/v1/ads/click", h.recordClick)
	r.Get("/v1/analytics/campaigns/{id}/performance", h.getCampaignPerformance)
	r.Route("/v1/admin/campaigns", func(r chi.Router) {
		r.Post("/{id}/pause", h.pauseCampaign)
		r.Post("/{id}/resume", h.resumeCampaign)
	})
	r.Get("/v1/admin/revenue", h.getRevenueStats)
}

func (h *Handler) createCampaign(w http.ResponseWriter, r *http.Request) {
	var c domain.Campaign
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateCampaign(r.Context(), &c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getCampaign(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	c, err := h.svc.GetCampaign(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "campaign not found")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) updateCampaign(w http.ResponseWriter, r *http.Request) {
	var c domain.Campaign
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	c.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateCampaign(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	advertiserID := uuid.MustParse(r.URL.Query().Get("advertiserId"))
	camps, err := h.svc.ListCampaigns(r.Context(), advertiserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, camps)
}

func (h *Handler) deleteCampaign(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) createAdUnit(w http.ResponseWriter, r *http.Request) {
	var u domain.AdUnit
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateAdUnit(r.Context(), &u)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getAdUnit(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	u, err := h.svc.GetAdUnit(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "ad unit not found")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) updateAdUnit(w http.ResponseWriter, r *http.Request) {
	var u domain.AdUnit
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	u.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateAdUnit(r.Context(), &u); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) deleteAdUnit(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteAdUnit(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) approveAdUnit(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		Approved bool   `json:"approved"`
		Note     string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ApproveAdUnit(r.Context(), id, req.Approved, req.Note); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) requestAd(w http.ResponseWriter, r *http.Request) {
	placementID := r.URL.Query().Get("placementId")
	var userID *uuid.UUID
	if uid := r.URL.Query().Get("userId"); uid != "" {
		id := uuid.MustParse(uid)
		userID = &id
	}
	ad, err := h.svc.RequestAd(r.Context(), placementID, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "no ads available")
		return
	}
	writeJSON(w, http.StatusOK, ad)
}

func (h *Handler) recordImpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CampaignID  uuid.UUID  `json:"campaignId"`
		AdUnitID    uuid.UUID  `json:"adUnitId"`
		PlacementID string     `json:"placementId"`
		UserID      *uuid.UUID `json:"userId,omitempty"`
		CostMicro   int64      `json:"costMicroCents"`
		Device      string     `json:"deviceType"`
		Geo         string     `json:"geo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RecordImpression(r.Context(), req.CampaignID, req.AdUnitID, req.PlacementID, req.UserID, req.CostMicro, req.Device, req.Geo); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func (h *Handler) recordClick(w http.ResponseWriter, r *http.Request) {
	var req struct{ ImpressionID uuid.UUID `json:"impressionId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RecordClick(r.Context(), req.ImpressionID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func (h *Handler) getCampaignPerformance(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	p, err := h.svc.GetCampaignPerformance(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "campaign not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) pauseCampaign(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.PauseCampaign(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "paused"})
}

func (h *Handler) resumeCampaign(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.ResumeCampaign(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resumed"})
}

func (h *Handler) getRevenueStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetRevenueStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/recommendation-service/internal/service"
)

type Handler struct {
	svc service.RecService
}

func New(svc service.RecService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/feeds", func(r chi.Router) {
		r.Get("/home", h.getHomeFeed)
		r.Get("/trending", h.getTrendingFeed)
		r.Get("/up-next/{contentId}", h.getUpNext)
		r.Get("/similar/{contentId}", h.getSimilar)
		r.Get("/creator/{creatorId}", h.getCreatorFeed)
	})
	r.Route("/v1/models", func(r chi.Router) {
		r.Get("/active", h.getActiveModel)
		r.Post("/deploy", h.deployModel)
		r.Get("/metrics", h.getModelMetrics)
	})
	r.Route("/v1/feedback", func(r chi.Router) {
		r.Post("/click", h.recordClick)
		r.Post("/dismiss", h.recordDismiss)
		r.Get("/explain/{recId}", h.explain)
	})
	r.Route("/v1/admin", func(r chi.Router) {
		r.Post("/refresh-features", h.refreshFeatures)
		r.Get("/feature-importance", h.getFeatureImportance)
		r.Post("/invalidate-cache", h.invalidateCache)
	})
}

func (h *Handler) getHomeFeed(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(r.URL.Query().Get("userId"))
	if userID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "userId required")
		return
	}
	feed, err := h.svc.GetHomeFeed(r.Context(), userID, parseLimit(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *Handler) getTrendingFeed(w http.ResponseWriter, r *http.Request) {
	feed, err := h.svc.GetTrendingFeed(r.Context(), parseLimit(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *Handler) getUpNext(w http.ResponseWriter, r *http.Request) {
	contentID := uuid.MustParse(chi.URLParam(r, "contentId"))
	userID := uuid.MustParse(r.URL.Query().Get("userId"))
	feed, err := h.svc.GetUpNext(r.Context(), userID, contentID, parseLimit(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *Handler) getSimilar(w http.ResponseWriter, r *http.Request) {
	contentID := uuid.MustParse(chi.URLParam(r, "contentId"))
	feed, err := h.svc.GetSimilar(r.Context(), contentID, parseLimit(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *Handler) getCreatorFeed(w http.ResponseWriter, r *http.Request) {
	creatorID := uuid.MustParse(chi.URLParam(r, "creatorId"))
	feed, err := h.svc.GetCreatorFeed(r.Context(), creatorID, parseLimit(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *Handler) getActiveModel(w http.ResponseWriter, r *http.Request) {
	model, err := h.svc.GetActiveModel(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model)
}

func (h *Handler) deployModel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Version string `json:"version"`
		Metrics string `json:"metrics"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.DeployModel(r.Context(), req.Version, req.Metrics); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deployed"})
}

func (h *Handler) getModelMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.svc.GetModelMetrics(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) recordClick(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    uuid.UUID `json:"userId"`
		ContentID uuid.UUID `json:"contentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RecordClick(r.Context(), req.UserID, req.ContentID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func (h *Handler) recordDismiss(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    uuid.UUID `json:"userId"`
		ContentID uuid.UUID `json:"contentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RecordDismiss(r.Context(), req.UserID, req.ContentID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "dismissed"})
}

func (h *Handler) explain(w http.ResponseWriter, r *http.Request) {
	recID := uuid.MustParse(chi.URLParam(r, "recId"))
	exp, err := h.svc.Explain(r.Context(), recID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, exp)
}

func (h *Handler) refreshFeatures(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.RefreshFeatures(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
}

func (h *Handler) getFeatureImportance(w http.ResponseWriter, r *http.Request) {
	imp, err := h.svc.GetFeatureImportance(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, imp)
}

func (h *Handler) invalidateCache(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.InvalidateCache(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "invalidated"})
}

func parseLimit(r *http.Request) int {
	return 50
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

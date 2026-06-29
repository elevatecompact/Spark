package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/moderation-service/internal/domain"
	"github.com/elevatecompact/spark/services/moderation-service/internal/service"
)

type Handler struct {
	svc service.ModerationService
}

func New(svc service.ModerationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/scan", func(r chi.Router) {
		r.Post("/text", h.scanText)
		r.Post("/image", h.scanImage)
		r.Post("/batch", h.scanBatch)
	})
	r.Route("/v1/rules", func(r chi.Router) {
		r.Get("/", h.listRules)
		r.Post("/", h.createRule)
		r.Patch("/{id}", h.updateRule)
		r.Delete("/{id}", h.deleteRule)
	})
	r.Route("/v1/queue", func(r chi.Router) {
		r.Get("/", h.listReviewQueue)
		r.Get("/stats", h.getQueueStats)
		r.Post("/{id}/approve", h.approveReview)
		r.Post("/{id}/reject", h.rejectReview)
	})
	r.Route("/v1/actions", func(r chi.Router) {
		r.Post("/warn", h.warnUser)
		r.Post("/restrict", h.restrictUser)
		r.Post("/remove", h.removeContent)
		r.Post("/suspend", h.suspendUser)
		r.Post("/reverse", h.reverseAction)
	})
	r.Route("/v1/reports", func(r chi.Router) {
		r.Post("/", h.createReport)
		r.Get("/{id}", h.getReport)
		r.Get("/", h.listReports)
	})
	r.Route("/v1/admin", func(r chi.Router) {
		r.Get("/stats", h.getAdminStats)
		r.Get("/accuracy", h.getAccuracy)
	})
}

func (h *Handler) scanText(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContentID uuid.UUID `json:"contentId"`
		Text      string    `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.ScanText(r.Context(), req.ContentID, req.Text)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) scanImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContentID uuid.UUID `json:"contentId"`
		ImageURL  string    `json:"imageUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.ScanImage(r.Context(), req.ContentID, req.ImageURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) scanBatch(w http.ResponseWriter, r *http.Request) {
	var items []domain.ScanResult
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	results, err := h.svc.ScanBatch(r.Context(), items)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *Handler) listRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.svc.ListRules(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func (h *Handler) createRule(w http.ResponseWriter, r *http.Request) {
	var rule domain.ModerationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.CreateRule(r.Context(), &rule); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rule)
}

func (h *Handler) updateRule(w http.ResponseWriter, r *http.Request) {
	var rule domain.ModerationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rule.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateRule(r.Context(), &rule); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rule)
}

func (h *Handler) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteRule(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) listReviewQueue(w http.ResponseWriter, r *http.Request) {
	status := domain.ReviewStatus(r.URL.Query().Get("status"))
	if status == "" {
		status = domain.ReviewPending
	}
	items, err := h.svc.ListReviewQueue(r.Context(), status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) getQueueStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetQueueStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) approveReview(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ Resolution string `json:"resolution"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ApproveReview(r.Context(), id, req.Resolution); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

func (h *Handler) rejectReview(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ Resolution string `json:"resolution"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RejectReview(r.Context(), id, req.Resolution); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (h *Handler) warnUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID uuid.UUID `json:"userId"`
		Reason string    `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	action, err := h.svc.WarnUser(r.Context(), req.UserID, req.Reason)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, action)
}

func (h *Handler) restrictUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   uuid.UUID `json:"userId"`
		Reason   string    `json:"reason"`
		Duration int       `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	action, err := h.svc.RestrictUser(r.Context(), req.UserID, req.Reason, req.Duration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, action)
}

func (h *Handler) removeContent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContentID uuid.UUID `json:"contentId"`
		Reason    string    `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	action, err := h.svc.RemoveContent(r.Context(), req.ContentID, req.Reason)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, action)
}

func (h *Handler) suspendUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   uuid.UUID `json:"userId"`
		Reason   string    `json:"reason"`
		Duration int       `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	action, err := h.svc.SuspendUser(r.Context(), req.UserID, req.Reason, req.Duration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, action)
}

func (h *Handler) reverseAction(w http.ResponseWriter, r *http.Request) {
	var req struct{ ActionID uuid.UUID `json:"actionId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ReverseAction(r.Context(), req.ActionID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reversed"})
}

func (h *Handler) createReport(w http.ResponseWriter, r *http.Request) {
	var report domain.ContentReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.CreateReport(r.Context(), &report); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	report, err := h.svc.GetReport(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "report not found")
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) listReports(w http.ResponseWriter, r *http.Request) {
	status := domain.ReportStatus(r.URL.Query().Get("status"))
	if status == "" {
		status = domain.ReportOpen
	}
	reports, err := h.svc.ListReports(r.Context(), status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reports)
}

func (h *Handler) getAdminStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetAdminStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) getAccuracy(w http.ResponseWriter, r *http.Request) {
	acc, err := h.svc.GetAccuracy(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, acc)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

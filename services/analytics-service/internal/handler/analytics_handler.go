package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/analytics-service/internal/domain"
	"github.com/elevatecompact/spark/services/analytics-service/internal/service"
)

type AnalyticsHandler struct {
	svc service.AnalyticsService
}

func NewAnalyticsHandler(svc service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.TrackEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	event, err := h.svc.TrackEvent(r.Context(), userID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

func (h *AnalyticsHandler) TrackBatch(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var reqs []domain.TrackEventRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	events, err := h.svc.TrackBatch(r.Context(), userID, reqs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, events)
}

func (h *AnalyticsHandler) GetCreatorDashboard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	dash, err := h.svc.GetDashboard(r.Context(), userID, domain.DashCreator)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"user_id": userID,
				"type":    "creator",
				"message": "no dashboard configured",
			})
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dash)
}

func (h *AnalyticsHandler) GetViewerDashboard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	dash, err := h.svc.GetDashboard(r.Context(), userID, domain.DashViewer)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"user_id": userID,
				"type":    "viewer",
				"message": "no dashboard configured",
			})
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dash)
}

func (h *AnalyticsHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	dash, err := h.svc.GetDashboard(r.Context(), uuid.Nil, domain.DashAdmin)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"type":    "admin",
				"message": "no dashboard configured",
			})
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dash)
}

func (h *AnalyticsHandler) GetRealtimeMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.svc.GetRealtimeMetrics(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

func (h *AnalyticsHandler) QueryMetrics(w http.ResponseWriter, r *http.Request) {
	var query domain.MetricQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	results, err := h.svc.QueryMetrics(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *AnalyticsHandler) GetHistoricalMetrics(w http.ResponseWriter, r *http.Request) {
	metricName := r.URL.Query().Get("metric")
	if metricName == "" {
		writeError(w, http.StatusBadRequest, "metric parameter required")
		return
	}
	results, err := h.svc.QueryMetrics(r.Context(), domain.MetricQuery{
		MetricName: metricName,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *AnalyticsHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req struct {
		Name   string          `json:"name"`
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	report, err := h.svc.GenerateReport(r.Context(), userID, req.Name, req.Type, req.Config)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, report)
}

func (h *AnalyticsHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid report id")
		return
	}
	report, err := h.svc.GetReport(r.Context(), id)
	if err != nil {
		if err == domain.ErrReportNotFound {
			writeError(w, http.StatusNotFound, "report not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *AnalyticsHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	reports, err := h.svc.ListReports(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reports)
}

func (h *AnalyticsHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.svc.ListTemplates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func (h *AnalyticsHandler) DefineFunnel(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req struct {
		Name  string          `json:"name"`
		Steps json.RawMessage `json:"steps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	funnel, err := h.svc.DefineFunnel(r.Context(), userID, req.Name, req.Steps)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, funnel)
}

func (h *AnalyticsHandler) GetFunnel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid funnel id")
		return
	}
	funnel, err := h.svc.GetFunnel(r.Context(), id)
	if err != nil {
		if err == domain.ErrFunnelNotFound {
			writeError(w, http.StatusNotFound, "funnel not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, funnel)
}

func (h *AnalyticsHandler) AnalyzeFunnel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid funnel id")
		return
	}
	funnel, err := h.svc.AnalyzeFunnel(r.Context(), id)
	if err != nil {
		if err == domain.ErrFunnelNotFound {
			writeError(w, http.StatusNotFound, "funnel not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, funnel)
}

func (h *AnalyticsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	var query domain.MetricQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	data, err := h.svc.ExportCSV(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *AnalyticsHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	var query domain.MetricQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	data, err := h.svc.ExportJSON(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

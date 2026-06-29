package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/trust-service/internal/domain"
	"github.com/elevatecompact/spark/services/trust-service/internal/service"
)

type Handler struct {
	svc *service.TrustService
}

func New(svc *service.TrustService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/v1/reputation/{userId}", h.getReputation)
	r.Get("/v1/reputation/{userId}/history", h.reputationHistory)
	r.Post("/v1/reputation/{userId}/recalculate", h.recalculateReputation)
	r.Get("/v1/trust/signals/{userId}", h.getSignals)
	r.Post("/v1/trust/signals", h.recordSignal)
	r.Get("/v1/trust/level/{userId}", h.getTrustLevel)
	r.Post("/v1/risk/assess", h.assessRisk)
	r.Get("/v1/risk/assessment/{id}", h.getRiskAssessment)
	r.Post("/v1/risk/rules", h.createRiskRule)
	r.Patch("/v1/risk/rules/{id}", h.updateRiskRule)
	r.Post("/v1/fraud/check-payment", h.checkPaymentFraud)
	r.Post("/v1/fraud/check-account", h.checkAccountFraud)
	r.Post("/v1/fraud/report", h.reportFraud)
	r.Get("/v1/fraud/cases", h.listFraudCases)
	r.Post("/v1/fraud/cases/{id}/resolve", h.resolveFraudCase)
	r.Get("/v1/admin/dashboard", h.dashboard)
	r.Get("/v1/admin/scores/distribution", h.scoreDistribution)
	r.Post("/v1/admin/thresholds/update", h.updateThresholds)
	r.Get("/v1/admin/flagged-users", h.flaggedUsers)
}

func (h *Handler) getReputation(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	rs, err := h.svc.GetReputation(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rs)
}

func (h *Handler) reputationHistory(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	history, err := h.svc.GetReputationHistory(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) recalculateReputation(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	rs, err := h.svc.RecalculateReputation(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rs)
}

func (h *Handler) getSignals(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	signals, err := h.svc.GetTrustSignals(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, signals)
}

func (h *Handler) recordSignal(w http.ResponseWriter, r *http.Request) {
	var s domain.TrustSignal
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.RecordSignal(r.Context(), &s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) getTrustLevel(w http.ResponseWriter, r *http.Request) {
	userID := uuid.MustParse(chi.URLParam(r, "userId"))
	rs, err := h.svc.GetTrustLevel(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rs)
}

func (h *Handler) assessRisk(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     uuid.UUID              `json:"userId"`
		ActionType string                 `json:"actionType"`
		Context    map[string]interface{} `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ra, err := h.svc.AssessRisk(r.Context(), req.UserID, req.ActionType, req.Context)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ra)
}

func (h *Handler) getRiskAssessment(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	ra, err := h.svc.GetRiskAssessment(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assessment not found")
		return
	}
	writeJSON(w, http.StatusOK, ra)
}

func (h *Handler) createRiskRule(w http.ResponseWriter, r *http.Request) {
	var rr domain.RiskRule
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateRiskRule(r.Context(), &rr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) updateRiskRule(w http.ResponseWriter, r *http.Request) {
	var rr domain.RiskRule
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rr.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateRiskRule(r.Context(), &rr); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rr)
}

func (h *Handler) checkPaymentFraud(w http.ResponseWriter, r *http.Request) {
	var ctx map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CheckPaymentFraud(r.Context(), ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) checkAccountFraud(w http.ResponseWriter, r *http.Request) {
	var ctx map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CheckAccountFraud(r.Context(), ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) reportFraud(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   uuid.UUID              `json:"userId"`
		Reason   string                 `json:"reason"`
		Evidence map[string]interface{} `json:"evidence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	fc, err := h.svc.ReportFraud(r.Context(), req.UserID, req.Reason, req.Evidence)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, fc)
}

func (h *Handler) listFraudCases(w http.ResponseWriter, r *http.Request) {
	var status *domain.FraudCaseStatus
	if s := r.URL.Query().Get("status"); s != "" {
		v := domain.FraudCaseStatus(s)
		status = &v
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	cases, err := h.svc.ListFraudCases(r.Context(), status, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cases)
}

func (h *Handler) resolveFraudCase(w http.ResponseWriter, r *http.Request) {
	caseID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		Status     string    `json:"status"`
		ReviewerID uuid.UUID `json:"reviewerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ResolveFraudCase(r.Context(), caseID, domain.FraudCaseStatus(req.Status), req.ReviewerID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	d, err := h.svc.GetDashboard(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) scoreDistribution(w http.ResponseWriter, r *http.Request) {
	sd, err := h.svc.GetScoreDistribution(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sd)
}

func (h *Handler) updateThresholds(w http.ResponseWriter, r *http.Request) {
	var thresholds map[string]int
	if err := json.NewDecoder(r.Body).Decode(&thresholds); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdateThresholds(r.Context(), thresholds); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) flaggedUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	users, err := h.svc.GetFlaggedUsers(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

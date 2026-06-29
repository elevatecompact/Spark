package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/licensing-service/internal/domain"
	"github.com/elevatecompact/spark/services/licensing-service/internal/service"
)

type Handler struct {
	svc *service.LicensingService
}

func New(svc *service.LicensingService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/licenses", func(r chi.Router) {
		r.Post("/", h.createLicense)
		r.Get("/", h.listLicenses)
		r.Get("/{id}", h.getLicense)
		r.Patch("/{id}", h.updateLicense)
		r.Delete("/{id}", h.deleteLicense)
	})
	r.Route("/v1/rights", func(r chi.Router) {
		r.Post("/content", h.registerContentRight)
		r.Get("/{contentId}", h.getContentRights)
		r.Post("/verify", h.verifyRights)
	})
	r.Route("/v1/royalties", func(r chi.Router) {
		r.Get("/calculate", h.calculateRoyalty)
		r.Post("/statement", h.generateStatement)
		r.Get("/statements", h.getStatements)
		r.Get("/pending", h.getPendingRoyalties)
	})
	r.Route("/v1/usage", func(r chi.Router) {
		r.Post("/record", h.recordUsage)
		r.Get("/report", h.usageReport)
		r.Get("/content/{id}", h.usageByContent)
	})
	r.Route("/v1/admin", func(r chi.Router) {
		r.Post("/licenses/{id}/approve", h.approveLicense)
		r.Post("/licenses/{id}/reject", h.rejectLicense)
		r.Post("/licenses/{id}/terminate", h.terminateLicense)
		r.Post("/royalties/process", h.processRoyalties)
		r.Get("/compliance/report", h.complianceReport)
	})
}

func (h *Handler) createLicense(w http.ResponseWriter, r *http.Request) {
	var l domain.License
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateLicense(r.Context(), &l)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) getLicense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	l, err := h.svc.GetLicense(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "license not found")
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *Handler) updateLicense(w http.ResponseWriter, r *http.Request) {
	var l domain.License
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	l.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateLicense(r.Context(), &l); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *Handler) deleteLicense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteLicense(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) listLicenses(w http.ResponseWriter, r *http.Request) {
	var rightsHolderID, licenseeID, contentID *uuid.UUID
	if id := r.URL.Query().Get("rightsHolderId"); id != "" {
		v := uuid.MustParse(id)
		rightsHolderID = &v
	}
	if id := r.URL.Query().Get("licenseeId"); id != "" {
		v := uuid.MustParse(id)
		licenseeID = &v
	}
	if id := r.URL.Query().Get("contentId"); id != "" {
		v := uuid.MustParse(id)
		contentID = &v
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	licenses, err := h.svc.ListLicenses(r.Context(), rightsHolderID, licenseeID, contentID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, licenses)
}

func (h *Handler) registerContentRight(w http.ResponseWriter, r *http.Request) {
	var cr domain.ContentRight
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.RegisterContentRight(r.Context(), &cr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) getContentRights(w http.ResponseWriter, r *http.Request) {
	contentID := uuid.MustParse(chi.URLParam(r, "contentId"))
	cr, err := h.svc.GetContentRights(r.Context(), contentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "rights not found")
		return
	}
	writeJSON(w, http.StatusOK, cr)
}

func (h *Handler) verifyRights(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContentID uuid.UUID `json:"contentId"`
		UsageType string    `json:"usageType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.VerifyContentRights(r.Context(), req.ContentID, req.UsageType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) calculateRoyalty(w http.ResponseWriter, r *http.Request) {
	licenseID := uuid.MustParse(r.URL.Query().Get("licenseId"))
	rs, err := h.svc.CalculateProjectedRoyalty(r.Context(), licenseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rs)
}

func (h *Handler) generateStatement(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LicenseID   uuid.UUID `json:"licenseId"`
		PeriodStart time.Time `json:"periodStart"`
		PeriodEnd   time.Time `json:"periodEnd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rs, err := h.svc.GenerateRoyaltyStatement(r.Context(), req.LicenseID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rs)
}

func (h *Handler) getStatements(w http.ResponseWriter, r *http.Request) {
	var rightsHolderID, licenseID *uuid.UUID
	if id := r.URL.Query().Get("rightsHolderId"); id != "" {
		v := uuid.MustParse(id)
		rightsHolderID = &v
	}
	if id := r.URL.Query().Get("licenseId"); id != "" {
		v := uuid.MustParse(id)
		licenseID = &v
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	statements, err := h.svc.GetRoyaltyStatements(r.Context(), rightsHolderID, licenseID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, statements)
}

func (h *Handler) getPendingRoyalties(w http.ResponseWriter, r *http.Request) {
	royalties, err := h.svc.GetPendingRoyalties(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, royalties)
}

func (h *Handler) recordUsage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LicenseID uuid.UUID              `json:"licenseId"`
		ContentID uuid.UUID              `json:"contentId"`
		UsageType string                 `json:"usageType"`
		Context   map[string]interface{} `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	u, err := h.svc.RecordUsage(r.Context(), req.LicenseID, req.ContentID, domain.UsageType(req.UsageType), req.Context)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *Handler) usageReport(w http.ResponseWriter, r *http.Request) {
	rightsHolderID := uuid.MustParse(r.URL.Query().Get("rightsHolderId"))
	start, _ := time.Parse(time.RFC3339, r.URL.Query().Get("periodStart"))
	end, _ := time.Parse(time.RFC3339, r.URL.Query().Get("periodEnd"))
	report, err := h.svc.GetUsageReport(r.Context(), rightsHolderID, start, end)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) usageByContent(w http.ResponseWriter, r *http.Request) {
	contentID := uuid.MustParse(chi.URLParam(r, "id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	logs, err := h.svc.GetUsageByContent(r.Context(), contentID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (h *Handler) approveLicense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.ApproveLicense(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

func (h *Handler) rejectLicense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.RejectLicense(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (h *Handler) terminateLicense(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.TerminateLicense(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "terminated"})
}

func (h *Handler) processRoyalties(w http.ResponseWriter, r *http.Request) {
	var req struct{ StatementID uuid.UUID `json:"statementId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ProcessRoyaltyPayout(r.Context(), req.StatementID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "processed"})
}

func (h *Handler) complianceReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.svc.GetComplianceReport(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
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

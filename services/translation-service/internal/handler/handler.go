package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/translation-service/internal/domain"
	"github.com/elevatecompact/spark/services/translation-service/internal/service"
)

type Handler struct {
	svc service.TranslationService
}

func New(svc service.TranslationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/v1/translate", h.translate)
	r.Post("/v1/translate/batch", h.translateBatch)
	r.Get("/v1/translate/languages", h.getLanguages)

	r.Post("/v1/detect", h.detect)
	r.Post("/v1/detect/batch", h.detectBatch)

	r.Get("/v1/tm/lookup", h.lookupMemory)
	r.Post("/v1/tm/store", h.storeMemory)
	r.Delete("/v1/tm/entries/{id}", h.deleteMemory)

	r.Get("/v1/review/queue", h.listReviewQueue)
	r.Post("/v1/review/{id}/approve", h.approveReview)
	r.Post("/v1/review/{id}/reject", h.rejectReview)

	r.Get("/v1/admin/usage", h.getUsage)
}

func (h *Handler) translate(w http.ResponseWriter, r *http.Request) {
	var req domain.TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.Translate(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) translateBatch(w http.ResponseWriter, r *http.Request) {
	var reqs []domain.TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	results, err := h.svc.TranslateBatch(r.Context(), reqs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *Handler) getLanguages(w http.ResponseWriter, r *http.Request) {
	langs, err := h.svc.GetLanguages(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, langs)
}

func (h *Handler) detect(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.Detect(r.Context(), req.Text)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) detectBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Texts []string `json:"texts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	results, err := h.svc.DetectBatch(r.Context(), req.Texts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *Handler) lookupMemory(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	sourceLang := r.URL.Query().Get("sourceLang")
	targetLang := r.URL.Query().Get("targetLang")
	entry, err := h.svc.LookupMemory(r.Context(), text, sourceLang, targetLang)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (h *Handler) storeMemory(w http.ResponseWriter, r *http.Request) {
	var entry domain.TranslationMemoryEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.StoreMemory(r.Context(), &entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stored"})
}

func (h *Handler) deleteMemory(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteMemory(r.Context(), id); err != nil {
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
	entries, err := h.svc.ListReviewQueue(r.Context(), status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *Handler) approveReview(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		ReviewerID    uuid.UUID `json:"reviewerId"`
		CorrectedText string    `json:"correctedText"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ApproveReview(r.Context(), id, req.ReviewerID, req.CorrectedText); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

func (h *Handler) rejectReview(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		ReviewerID uuid.UUID `json:"reviewerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RejectReview(r.Context(), id, req.ReviewerID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (h *Handler) getUsage(w http.ResponseWriter, r *http.Request) {
	usage, err := h.svc.GetUsage(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, usage)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

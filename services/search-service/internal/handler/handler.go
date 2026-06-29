package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/search-service/internal/domain"
	"github.com/elevatecompact/spark/services/search-service/internal/service"
)

type Handler struct {
	svc service.SearchService
}

func New(svc service.SearchService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/search", func(r chi.Router) {
		r.Get("/", h.search)
	})
	r.Get("/v1/autocomplete", h.autocomplete)
	r.Route("/v1/index", func(r chi.Router) {
		r.Post("/{contentType}", h.indexDoc)
		r.Put("/{contentType}/{id}", h.updateDoc)
		r.Delete("/{contentType}/{id}", h.removeDoc)
		r.Post("/reindex", h.reindex)
	})
	r.Post("/v1/suggestions/click", h.suggestionClick)
	r.Route("/v1/admin", func(r chi.Router) {
		r.Get("/stats", h.getStats)
		r.Post("/synonyms", h.putSynonyms)
		r.Put("/analyzers", h.putAnalyzers)
		r.Get("/health", h.health)
	})
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	q := &domain.SearchQuery{
		Query:   r.URL.Query().Get("q"),
		Type:    domain.ContentType(r.URL.Query().Get("type")),
		Sort:    domain.SortBy(r.URL.Query().Get("sort")),
		Page:    queryInt(r, "page", 1),
		Size:    queryInt(r, "size", 20),
	}
	if uid := r.URL.Query().Get("userId"); uid != "" {
		q.UserID = uuid.MustParse(uid)
	}
	result, err := h.svc.Search(r.Context(), q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) autocomplete(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("q")
	size := queryInt(r, "size", 10)
	suggestions, err := h.svc.Autocomplete(r.Context(), prefix, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, suggestions)
}

func (h *Handler) indexDoc(w http.ResponseWriter, r *http.Request) {
	ct := domain.ContentType(chi.URLParam(r, "contentType"))
	var doc domain.SearchDocument
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.IndexDocument(r.Context(), ct, &doc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "indexed"})
}

func (h *Handler) updateDoc(w http.ResponseWriter, r *http.Request) {
	ct := domain.ContentType(chi.URLParam(r, "contentType"))
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var doc map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdateDocument(r.Context(), ct, id, doc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) removeDoc(w http.ResponseWriter, r *http.Request) {
	ct := domain.ContentType(chi.URLParam(r, "contentType"))
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.RemoveDocument(r.Context(), ct, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (h *Handler) reindex(w http.ResponseWriter, r *http.Request) {
	ct := domain.ContentType(r.URL.Query().Get("type"))
	if err := h.svc.Reindex(r.Context(), ct); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reindex started"})
}

func (h *Handler) suggestionClick(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     uuid.UUID `json:"userId"`
		Suggestion string    `json:"suggestion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RecordSuggestionClick(r.Context(), req.UserID, req.Suggestion); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) putSynonyms(w http.ResponseWriter, r *http.Request) {
	var set domain.SynonymSet
	if err := json.NewDecoder(r.Body).Decode(&set); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.PutSynonyms(r.Context(), &set); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "synonyms updated"})
}

func (h *Handler) putAnalyzers(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.PutAnalyzers(r.Context(), config); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "analyzers updated"})
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	health, err := h.svc.Health(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, health)
}

func queryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

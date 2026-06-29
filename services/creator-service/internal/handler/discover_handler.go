package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/service"
)

type DiscoverHandler struct {
	discoverService *service.DiscoverService
}

func NewDiscoverHandler(discoverService *service.DiscoverService) *DiscoverHandler {
	return &DiscoverHandler{
		discoverService: discoverService,
	}
}

func (h *DiscoverHandler) Trending(w http.ResponseWriter, r *http.Request) {
	limit, offset := getLimitOffset(r)
	creators, err := h.discoverService.GetTrending(r.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   creators,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *DiscoverHandler) Recommended(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserID(r.Context())
	limit, _ := getLimitOffset(r)

	var userID uuid.UUID
	var err error
	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			userID = uuid.Nil
		}
	}

	creators, err := h.discoverService.GetRecommended(r.Context(), userID, limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  creators,
		"limit": limit,
	})
}

func (h *DiscoverHandler) Nearby(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	limit, offset := getLimitOffset(r)

	creators, err := h.discoverService.GetNearby(r.Context(), country, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   creators,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *DiscoverHandler) ByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := uuid.Parse(chi.URLParam(r, "categoryID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	limit, offset := getLimitOffset(r)
	creators, total, err := h.discoverService.GetByCategory(r.Context(), categoryID, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   creators,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *DiscoverHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := service.SearchFilters{
		Query:    q.Get("q"),
		Category: q.Get("category"),
		Language: q.Get("language"),
		Country:  q.Get("country"),
	}
	filters.Limit, filters.Offset = getLimitOffset(r)

	creators, total, err := h.discoverService.Search(r.Context(), filters)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   creators,
		"total":  total,
		"limit":  filters.Limit,
		"offset": filters.Offset,
	})
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/discovery-service/internal/domain"
	"github.com/elevatecompact/spark/services/discovery-service/internal/service"
)

type Handler struct {
	svc *service.DiscoveryService
}

func New(svc *service.DiscoveryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/feeds", func(r chi.Router) {
		r.Get("/home", h.homeFeed)
		r.Get("/trending", h.trendingFeed)
		r.Get("/category/{slug}", h.categoryFeed)
		r.Get("/new", h.newFeed)
		r.Get("/related/{contentId}", h.relatedFeed)
	})
	r.Route("/v1/categories", func(r chi.Router) {
		r.Get("/", h.listCategories)
		r.Get("/{slug}", h.getCategory)
		r.Get("/{slug}/contents", h.categoryContents)
	})
	r.Route("/v1/collections", func(r chi.Router) {
		r.Get("/", h.listCollections)
		r.Get("/{id}", h.getCollection)
		r.Post("/", h.createCollection)
		r.Patch("/{id}", h.updateCollection)
		r.Post("/{id}/items", h.addCollectionItem)
		r.Delete("/{id}/items/{contentId}", h.removeCollectionItem)
	})
	r.Route("/v1/trending", func(r chi.Router) {
		r.Get("/", h.getTrending)
		r.Get("/category/{slug}", h.trendingByCategory)
		r.Get("/creators", h.trendingCreators)
	})
	r.Route("/v1/editorial", func(r chi.Router) {
		r.Get("/picks", h.staffPicks)
		r.Get("/spotlight", h.spotlight)
		r.Get("/holiday/{campaign}", h.holidayPicks)
	})
	r.Route("/v1/admin", func(r chi.Router) {
		r.Post("/feeds/cache/warm", h.warmFeeds)
		r.Post("/trending/refresh", h.refreshTrending)
		r.Post("/categories/reorder", h.reorderCategories)
	})
}

func (h *Handler) homeFeed(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var userID *uuid.UUID
	if uid := r.URL.Query().Get("userId"); uid != "" {
		id := uuid.MustParse(uid)
		userID = &id
	}
	ids, err := h.svc.GetHomeFeed(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) trendingFeed(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	items, err := h.svc.GetTrendingFeed(r.Context(), r.URL.Query().Get("timeframe"), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) categoryFeed(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ids, err := h.svc.GetCategoryFeed(r.Context(), slug, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) newFeed(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ids, err := h.svc.GetNewFeed(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) relatedFeed(w http.ResponseWriter, r *http.Request) {
	contentID := uuid.MustParse(chi.URLParam(r, "contentId"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	ids, err := h.svc.GetRelatedFeed(r.Context(), contentID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.GetCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *Handler) getCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cat, err := h.svc.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	subs, _ := h.svc.GetSubcategories(r.Context(), cat.ID)
	result := map[string]interface{}{"category": cat, "subcategories": subs}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) categoryContents(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ids, err := h.svc.GetCategoryFeed(r.Context(), slug, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) listCollections(w http.ResponseWriter, r *http.Request) {
	featured := r.URL.Query().Get("featured") == "true"
	cols, err := h.svc.ListCollections(r.Context(), featured)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cols)
}

func (h *Handler) getCollection(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	col, err := h.svc.GetCollection(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "collection not found")
		return
	}
	writeJSON(w, http.StatusOK, col)
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var c domain.Collection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateCollection(r.Context(), &c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) updateCollection(w http.ResponseWriter, r *http.Request) {
	var c domain.Collection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	c.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateCollection(r.Context(), &c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) addCollectionItem(w http.ResponseWriter, r *http.Request) {
	collectionID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct {
		ContentID uuid.UUID `json:"contentId"`
		SortOrder int       `json:"sortOrder"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.AddCollectionItem(r.Context(), collectionID, req.ContentID, req.SortOrder); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "added"})
}

func (h *Handler) removeCollectionItem(w http.ResponseWriter, r *http.Request) {
	collectionID := uuid.MustParse(chi.URLParam(r, "id"))
	contentID := uuid.MustParse(chi.URLParam(r, "contentId"))
	if err := h.svc.RemoveCollectionItem(r.Context(), collectionID, contentID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) getTrending(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	items, err := h.svc.GetTrending(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) trendingByCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	items, err := h.svc.GetTrendingByCategory(r.Context(), slug, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) trendingCreators(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	ids, err := h.svc.GetTrendingCreators(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ids)
}

func (h *Handler) staffPicks(w http.ResponseWriter, r *http.Request) {
	picks, err := h.svc.GetStaffPicks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, picks)
}

func (h *Handler) spotlight(w http.ResponseWriter, r *http.Request) {
	picks, err := h.svc.GetSpotlight(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, picks)
}

func (h *Handler) holidayPicks(w http.ResponseWriter, r *http.Request) {
	campaign := chi.URLParam(r, "campaign")
	picks, err := h.svc.GetHolidayPicks(r.Context(), campaign)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, picks)
}

func (h *Handler) warmFeeds(w http.ResponseWriter, r *http.Request) {
	var req struct{ FeedTypes []string `json:"feedTypes"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.WarmFeedCache(r.Context(), req.FeedTypes); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "warmed"})
}

func (h *Handler) refreshTrending(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.RefreshTrending(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "refreshed"})
}

func (h *Handler) reorderCategories(w http.ResponseWriter, r *http.Request) {
	var req struct{ Order []uuid.UUID `json:"order"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ReorderCategories(r.Context(), req.Order); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reordered"})
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

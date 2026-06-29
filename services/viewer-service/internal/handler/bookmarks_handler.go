package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/service"
)

type BookmarkHandler struct {
	bookmarkSvc   service.BookmarkService
	watchLaterSvc service.WatchLaterService
}

func NewBookmarkHandler(bookmarkSvc service.BookmarkService, watchLaterSvc service.WatchLaterService) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkSvc:   bookmarkSvc,
		watchLaterSvc: watchLaterSvc,
	}
}

type CreateBookmarkRequest struct {
	ContentID uuid.UUID `json:"content_id"`
	Note      string    `json:"note,omitempty"`
	Folder    string    `json:"folder,omitempty"`
}

func (h *BookmarkHandler) Create(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bookmark, err := h.bookmarkSvc.Create(r.Context(), viewerID, req.ContentID, req.Note, req.Folder)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, bookmark)
}

func (h *BookmarkHandler) List(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	folder := r.URL.Query().Get("folder")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	bookmarks, err := h.bookmarkSvc.List(r.Context(), viewerID, folder, limit, offset)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"bookmarks": bookmarks,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *BookmarkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	bookmarkID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid bookmark id")
		return
	}

	if err := h.bookmarkSvc.Delete(r.Context(), viewerID, bookmarkID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AddWatchLaterRequest struct {
	ContentID uuid.UUID `json:"content_id"`
}

func (h *BookmarkHandler) AddWatchLater(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req AddWatchLaterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.watchLaterSvc.Add(r.Context(), viewerID, req.ContentID)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, item)
}

func (h *BookmarkHandler) ListWatchLater(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.watchLaterSvc.List(r.Context(), viewerID)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
	})
}

func (h *BookmarkHandler) RemoveWatchLater(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid watch later id")
		return
	}

	if err := h.watchLaterSvc.Remove(r.Context(), viewerID, itemID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

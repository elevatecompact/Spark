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

type WatchHistoryHandler struct {
	watchSvc service.WatchHistoryService
}

func NewWatchHistoryHandler(watchSvc service.WatchHistoryService) *WatchHistoryHandler {
	return &WatchHistoryHandler{watchSvc: watchSvc}
}

type RecordWatchRequest struct {
	ContentID            uuid.UUID              `json:"content_id"`
	ContentType          domain.ContentType     `json:"content_type"`
	Progress             float64                `json:"progress"`
	WatchDurationSeconds int                    `json:"watch_duration_seconds"`
	Completed            bool                   `json:"completed"`
}

func (h *WatchHistoryHandler) RecordWatch(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req RecordWatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	update := domain.WatchProgressUpdate{
		ContentID:            req.ContentID,
		ContentType:          req.ContentType,
		Progress:             req.Progress,
		WatchDurationSeconds: req.WatchDurationSeconds,
		Completed:            req.Completed,
	}

	entry, err := h.watchSvc.RecordWatch(r.Context(), viewerID, update)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, entry)
}

func (h *WatchHistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	contentType := r.URL.Query().Get("content_type")
	daysStr := r.URL.Query().Get("days")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	days := 90
	limit := 20
	offset := 0

	if daysStr != "" {
		if v, err := strconv.Atoi(daysStr); err == nil && v > 0 {
			days = v
		}
	}
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

	entries, err := h.watchSvc.GetHistory(r.Context(), viewerID, contentType, days, limit, offset)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *WatchHistoryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	entryID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	if err := h.watchSvc.DeleteEntry(r.Context(), viewerID, entryID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WatchHistoryHandler) ClearHistory(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.watchSvc.ClearHistory(r.Context(), viewerID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

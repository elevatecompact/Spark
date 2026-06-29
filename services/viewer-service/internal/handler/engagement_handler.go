package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/service"
)

type EngagementHandler struct {
	engSvc service.EngagementService
}

func NewEngagementHandler(engSvc service.EngagementService) *EngagementHandler {
	return &EngagementHandler{engSvc: engSvc}
}

type RateContentRequest struct {
	ContentID uuid.UUID `json:"content_id"`
	Score     int       `json:"score"`
}

func (h *EngagementHandler) RateContent(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req RateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rating, err := h.engSvc.RateContent(r.Context(), viewerID, req.ContentID, req.Score)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, rating)
}

type ToggleReactionRequest struct {
	ContentID uuid.UUID          `json:"content_id"`
	Type      domain.ReactionType `json:"type"`
}

func (h *EngagementHandler) ToggleReaction(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ToggleReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reaction, err := h.engSvc.ToggleReaction(r.Context(), viewerID, req.ContentID, req.Type)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	if reaction == nil {
		WriteJSON(w, http.StatusOK, map[string]string{"message": "reaction removed"})
		return
	}

	WriteJSON(w, http.StatusOK, reaction)
}

type ReportContentRequest struct {
	ContentID   uuid.UUID        `json:"content_id"`
	Type        domain.ReportType `json:"type"`
	Description string           `json:"description,omitempty"`
}

func (h *EngagementHandler) ReportContent(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ReportContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.engSvc.ReportContent(r.Context(), viewerID, req.ContentID, req.Type, req.Description)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, report)
}

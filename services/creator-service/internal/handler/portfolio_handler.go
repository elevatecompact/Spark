package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
)

type PortfolioHandler struct {
	portfolioRepo repository.PortfolioRepository
}

func NewPortfolioHandler(portfolioRepo repository.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioRepo: portfolioRepo,
	}
}

func (h *PortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	items, err := h.portfolioRepo.GetByCreatorID(r.Context(), creatorID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
	})
}

func (h *PortfolioHandler) Create(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	var req domain.CreatePortfolioItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.MediaURL == "" || req.MediaType == "" {
		writeError(w, http.StatusBadRequest, "title, media_url, and media_type are required")
		return
	}

	validTypes := map[string]bool{"video": true, "image": true, "audio": true}
	if !validTypes[req.MediaType] {
		writeError(w, http.StatusBadRequest, "media_type must be video, image, or audio")
		return
	}

	item := &domain.PortfolioItem{
		ID:           uuid.New(),
		CreatorID:    creatorID,
		Title:        req.Title,
		Description:  req.Description,
		MediaURL:     req.MediaURL,
		MediaType:    req.MediaType,
		ThumbnailURL: req.ThumbnailURL,
		Featured:     req.Featured,
		SortOrder:    req.SortOrder,
	}

	if err := h.portfolioRepo.Create(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *PortfolioHandler) Update(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item ID")
		return
	}

	existing, err := h.portfolioRepo.GetByID(r.Context(), itemID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	var req domain.CreatePortfolioItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.MediaURL != "" {
		existing.MediaURL = req.MediaURL
	}
	if req.MediaType != "" {
		existing.MediaType = req.MediaType
	}
	if req.ThumbnailURL != "" {
		existing.ThumbnailURL = req.ThumbnailURL
	}
	existing.Featured = req.Featured
	existing.SortOrder = req.SortOrder

	if err := h.portfolioRepo.Update(r.Context(), existing); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

func (h *PortfolioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item ID")
		return
	}

	if err := h.portfolioRepo.Delete(r.Context(), itemID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *PortfolioHandler) SetFeatured(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item ID")
		return
	}

	if err := h.portfolioRepo.SetFeatured(r.Context(), itemID, creatorID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "featured"})
}

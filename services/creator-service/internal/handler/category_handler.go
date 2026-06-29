package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
	"github.com/elevatecompact/spark/services/creator-service/internal/service"
)

type CategoryHandler struct {
	categoryRepo  repository.CategoryRepository
	creatorService *service.CreatorService
}

func NewCategoryHandler(categoryRepo repository.CategoryRepository, creatorService *service.CreatorService) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo:  categoryRepo,
		creatorService: creatorService,
	}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	categories, err := h.categoryRepo.List(r.Context(), activeOnly)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": categories,
	})
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Slug == "" {
		writeError(w, http.StatusBadRequest, "name and slug are required")
		return
	}

	now := &domain.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconURL:     req.IconURL,
		Color:       req.Color,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		Active:      true,
	}

	if err := h.categoryRepo.Create(r.Context(), now); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, now)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	category, err := h.categoryRepo.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) GetCreators(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	limit, offset := getLimitOffset(r)
	creators, total, err := h.creatorService.GetByCategory(r.Context(), id, limit, offset)
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

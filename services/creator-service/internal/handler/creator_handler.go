package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/service"
)

type CreatorHandler struct {
	creatorService *service.CreatorService
	verifyService  *service.VerificationService
}

func NewCreatorHandler(creatorService *service.CreatorService, verifyService *service.VerificationService) *CreatorHandler {
	return &CreatorHandler{
		creatorService: creatorService,
		verifyService:  verifyService,
	}
}

func (h *CreatorHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserID(r.Context())
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req domain.CreateCreatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "display name is required")
		return
	}
	if req.Language == "" {
		writeError(w, http.StatusBadRequest, "language is required")
		return
	}
	if req.Country == "" {
		writeError(w, http.StatusBadRequest, "country is required")
		return
	}

	creator, err := h.creatorService.CreateProfile(r.Context(), userID, req)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, creator)
}

func (h *CreatorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	creator, err := h.creatorService.GetProfile(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, creator)
}

func (h *CreatorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	userIDStr := GetUserID(r.Context())
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	creator, err := h.creatorService.GetProfile(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	userID, _ := uuid.Parse(userIDStr)
	if creator.UserID != userID && GetUserRole(r.Context()) != "admin" {
		writeError(w, http.StatusForbidden, "you can only update your own profile")
		return
	}

	var req domain.UpdateCreatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.creatorService.UpdateProfile(r.Context(), id, req); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *CreatorHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("q")
	categories := splitCSV(q.Get("categories"))
	tags := splitCSV(q.Get("tags"))
	language := q.Get("language")
	country := q.Get("country")
	limit, offset := getLimitOffset(r)

	creators, total, err := h.creatorService.SearchCreators(r.Context(), query, categories, tags, language, country, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  creators,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

func (h *CreatorHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	limit, offset := getLimitOffset(r)
	followers, total, err := h.creatorService.GetFollowers(r.Context(), id, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   followers,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *CreatorHandler) Follow(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	userIDStr := GetUserID(r.Context())
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	followerID, _ := uuid.Parse(userIDStr)
	if err := h.creatorService.FollowCreator(r.Context(), followerID, creatorID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "following"})
}

func (h *CreatorHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	userIDStr := GetUserID(r.Context())
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	followerID, _ := uuid.Parse(userIDStr)
	if err := h.creatorService.UnfollowCreator(r.Context(), followerID, creatorID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "unfollowed"})
}

func (h *CreatorHandler) Verify(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	adminIDStr := GetUserID(r.Context())
	if adminIDStr == "" || GetUserRole(r.Context()) != "admin" {
		writeError(w, http.StatusForbidden, "admin access required")
		return
	}

	adminID, _ := uuid.Parse(adminIDStr)
	if err := h.verifyService.ApproveVerification(r.Context(), creatorID, adminID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "verified"})
}

func (h *CreatorHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	limit, offset := getLimitOffset(r)
	following, total, err := h.creatorService.GetFollowing(r.Context(), id, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   following,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

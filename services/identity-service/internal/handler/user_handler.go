package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userSvc.GetProfile(r.Context(), user.ID)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, profile)
}

type UpdateUserRequest struct {
	Username    *string   `json:"username,omitempty"`
	DisplayName *string   `json:"display_name,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	BannerURL   *string   `json:"banner_url,omitempty"`
	Categories  *[]string `json:"categories,omitempty"`
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := service.UpdateUserProfile{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
		BannerURL:   req.BannerURL,
		Categories:  req.Categories,
	}

	updatedUser, err := h.userSvc.UpdateProfile(r.Context(), user.ID, updates)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, updatedUser)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.userSvc.ChangePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.userSvc.DeleteAccount(r.Context(), user.ID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	profile, err := h.userSvc.GetPublicProfile(r.Context(), userID)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	role := r.URL.Query().Get("role")
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

	users, err := h.userSvc.SearchUsers(r.Context(), query, role, limit, offset)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *UserHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.userSvc.VerifyUser(r.Context(), userID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "user verified successfully"})
}

func (h *UserHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	reason := r.URL.Query().Get("reason")

	if err := h.userSvc.SuspendUser(r.Context(), userID, reason); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "user suspended successfully"})
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	role := domain.UserRole(req.Role)
	switch role {
	case domain.RoleViewer, domain.RoleCreator, domain.RoleAdmin, domain.RoleMod:
	default:
		WriteError(w, http.StatusBadRequest, "invalid role")
		return
	}

	if err := h.userSvc.UpdateRole(r.Context(), userID, role); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "role updated successfully"})
}

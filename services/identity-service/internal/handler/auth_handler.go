package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/domain"
	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type RegisterResponse struct {
	User    *domain.User    `json:"user"`
	Session *domain.Session `json:"session"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authSvc.Register(r.Context(), req.Email, req.Username, req.Password, req.DisplayName)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, RegisterResponse{
		User:    user,
		Session: session,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User    *domain.User    `json:"user"`
	Session *domain.Session `json:"session"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := splitAndTrim(forwarded, ",")
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
	userAgent := r.Header.Get("User-Agent")

	session, err := h.authSvc.Login(r.Context(), req.Email, req.Password, ip, userAgent)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	user, err := h.authSvc.ValidateToken(r.Context(), session.Token)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load user")
		return
	}

	WriteJSON(w, http.StatusOK, LoginResponse{
		User:    user,
		Session: session,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Session *domain.Session `json:"session"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		WriteError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	session, err := h.authSvc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, RefreshResponse{Session: session})
}

type LogoutRequest struct {
	SessionID string `json:"session_id"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid session_id")
		return
	}

	if err := h.authSvc.Logout(r.Context(), sessionID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authSvc.LogoutAll(r.Context(), user.ID); err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpaceStr(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	if start <= len(s) {
		result = append(result, s[start:])
	}
	return result
}

func trimSpaceStr(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

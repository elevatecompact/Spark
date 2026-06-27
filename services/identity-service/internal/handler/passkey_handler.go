package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

type PasskeyHandler struct {
	passkeySvc service.PasskeyService
}

func NewPasskeyHandler(passkeySvc service.PasskeyService) *PasskeyHandler {
	return &PasskeyHandler{passkeySvc: passkeySvc}
}

func (h *PasskeyHandler) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	options, err := h.passkeySvc.BeginRegistration(r.Context(), user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, options)
}

func (h *PasskeyHandler) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var response service.AuthenticatorAttestationResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.passkeySvc.FinishRegistration(r.Context(), user, response); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *PasskeyHandler) BeginAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	options, err := h.passkeySvc.BeginAuthentication(r.Context(), user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, options)
}

func (h *PasskeyHandler) FinishAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var response service.AuthenticatorAssertionResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.passkeySvc.FinishAuthentication(r.Context(), user, response)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"session": session,
	})
}

func (h *PasskeyHandler) ListPasskeys(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	passkeys, err := h.passkeySvc.GetPasskeys(r.Context(), user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{"passkeys": passkeys})
}

func (h *PasskeyHandler) DeletePasskey(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromContext(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	passkeyID, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid passkey id")
		return
	}

	if err := h.passkeySvc.DeletePasskey(r.Context(), passkeyID, user.ID); err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

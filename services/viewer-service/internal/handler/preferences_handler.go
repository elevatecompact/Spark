package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/service"
)

type PreferencesHandler struct {
	prefsSvc service.PreferencesService
}

func NewPreferencesHandler(prefsSvc service.PreferencesService) *PreferencesHandler {
	return &PreferencesHandler{prefsSvc: prefsSvc}
}

func (h *PreferencesHandler) Get(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	prefs, err := h.prefsSvc.Get(r.Context(), viewerID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get preferences")
		return
	}

	WriteJSON(w, http.StatusOK, prefs)
}

func (h *PreferencesHandler) Replace(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var prefs domain.ViewerPreferences
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.prefsSvc.Replace(r.Context(), viewerID, &prefs)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func (h *PreferencesHandler) Patch(w http.ResponseWriter, r *http.Request) {
	viewerID, err := GetViewerID(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var updates domain.UpdatePreferences
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.prefsSvc.Patch(r.Context(), viewerID, updates)
	if err != nil {
		status := domain.HTTPStatusFromError(err)
		WriteError(w, status, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func (h *PreferencesHandler) GetDefault(w http.ResponseWriter, r *http.Request) {
	prefs := domain.ViewerPreferences{
		PreferredCategories: []uuid.UUID{},
		ContentLanguage:     "en",
		Autoplay:            true,
		MatureContentAllowed: false,
		NotificationPrefs:   make(map[string]interface{}),
	}
	WriteJSON(w, http.StatusOK, prefs)
}

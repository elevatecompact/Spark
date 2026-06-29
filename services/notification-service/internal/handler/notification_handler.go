package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/notification-service/internal/domain"
	"github.com/elevatecompact/spark/services/notification-service/internal/service"
)

type NotifHandler struct {
	svc service.NotificationService
}

func NewNotifHandler(svc service.NotificationService) *NotifHandler {
	return &NotifHandler{svc: svc}
}

func getUserID(r *http.Request) uuid.UUID {
	idStr := r.Header.Get("X-User-ID")
	if idStr == "" {
		return uuid.Nil
	}
	id, _ := uuid.Parse(idStr)
	return id
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *NotifHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	cursor := time.Time{}
	if c := r.URL.Query().Get("cursor"); c != "" {
		cursor.UnmarshalText([]byte(c))
	}
	notifs, err := h.svc.ListNotifications(r.Context(), uid, cursor, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, notifs)
}

func (h *NotifHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.MarkRead(r.Context(), id, uid); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "read"})
}

func (h *NotifHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	if err := h.svc.MarkAllRead(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "all read"})
}

func (h *NotifHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id, uid); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotifHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	prefs, err := h.svc.GetPreferences(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, prefs)
}

func (h *NotifHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req struct {
		Preferences string `json:"preferences"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdatePreferences(r.Context(), uid, req.Preferences); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *NotifHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req domain.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	n, err := h.svc.SendNotification(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (h *NotifHandler) SendBatch(w http.ResponseWriter, r *http.Request) {
	var reqs []domain.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ns, err := h.svc.SendBatch(r.Context(), reqs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ns)
}

func (h *NotifHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	var req domain.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	dev, err := h.svc.RegisterDevice(r.Context(), uid, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, dev)
}

func (h *NotifHandler) UnregisterDevice(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.UnregisterDevice(r.Context(), id, uid); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotifHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	if uid == uuid.Nil {
		writeError(w, http.StatusUnauthorized, "missing user id")
		return
	}
	devices, err := h.svc.ListDevices(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, devices)
}

func (h *NotifHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.svc.ListTemplates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func (h *NotifHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var t domain.Template
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.CreateTemplate(r.Context(), &t); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *NotifHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var t domain.Template
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	t.ID = id
	if err := h.svc.UpdateTemplate(r.Context(), &t); err != nil {
		writeError(w, domain.HTTPStatusFromError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *NotifHandler) TestPush(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.TestPush(r.Context(), req.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "push sent"})
}

func (h *NotifHandler) TestEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.TestEmail(r.Context(), req.Email); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "email sent"})
}

func (h *NotifHandler) DeliveryStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.DeliveryStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

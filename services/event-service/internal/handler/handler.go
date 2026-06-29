package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/event-service/internal/domain"
	"github.com/elevatecompact/spark/services/event-service/internal/service"
)

type Handler struct {
	svc service.EventService
}

func New(svc service.EventService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/events", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)

		r.Post("/{id}/tickets", h.createTicketTier)
		r.Get("/{id}/tickets", h.listTicketTiers)
		r.Post("/{id}/rsvp", h.rsvp)

		r.Get("/{id}/schedule", h.listSessions)
		r.Post("/{id}/schedule/sessions", h.createSession)

		r.Get("/{id}/attendees", h.listAttendees)
	})
	r.Post("/v1/tickets/{id}/purchase", h.purchaseTicket)

	r.Patch("/v1/sessions/{id}", h.updateSession)

	r.Route("/v1/series", func(r chi.Router) {
		r.Post("/", h.createSeries)
		r.Get("/{id}", h.getSeries)
		r.Patch("/{id}", h.updateSeries)
		r.Delete("/{id}", h.deleteSeries)
	})

	r.Route("/v1/admin/events", func(r chi.Router) {
		r.Post("/{id}/cancel", h.cancelEvent)
		r.Get("/stats", h.getAdminStats)
	})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.Create(r.Context(), &e)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	e, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "event not found")
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var e domain.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	e.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.Update(r.Context(), &e); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	status := domain.EventStatus(r.URL.Query().Get("status"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	events, err := h.svc.List(r.Context(), category, status, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) createTicketTier(w http.ResponseWriter, r *http.Request) {
	eventID := uuid.MustParse(chi.URLParam(r, "id"))
	var t domain.EventTicketTier
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	t.EventID = eventID
	if err := h.svc.CreateTicketTier(r.Context(), &t); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) listTicketTiers(w http.ResponseWriter, r *http.Request) {
	eventID := uuid.MustParse(chi.URLParam(r, "id"))
	tiers, err := h.svc.ListTicketTiers(r.Context(), eventID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tiers)
}

func (h *Handler) rsvp(w http.ResponseWriter, r *http.Request) {
	eventID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RSVP(r.Context(), eventID, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "registered"})
}

func (h *Handler) purchaseTicket(w http.ResponseWriter, r *http.Request) {
	tierID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ UserID uuid.UUID `json:"userId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.PurchaseTicket(r.Context(), tierID, req.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "purchased"})
}

func (h *Handler) listAttendees(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	attendees, err := h.svc.ListAttendees(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, attendees)
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	eventID := uuid.MustParse(chi.URLParam(r, "id"))
	var s domain.EventSession
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	s.EventID = eventID
	result, err := h.svc.CreateSession(r.Context(), &s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) updateSession(w http.ResponseWriter, r *http.Request) {
	var s domain.EventSession
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	s.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateSession(r.Context(), &s); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) listSessions(w http.ResponseWriter, r *http.Request) {
	eventID := uuid.MustParse(chi.URLParam(r, "id"))
	sessions, err := h.svc.ListSessions(r.Context(), eventID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

func (h *Handler) createSeries(w http.ResponseWriter, r *http.Request) {
	var s domain.EventSeries
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateSeries(r.Context(), &s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getSeries(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	s, err := h.svc.GetSeries(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "series not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) updateSeries(w http.ResponseWriter, r *http.Request) {
	var s domain.EventSeries
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	s.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateSeries(r.Context(), &s); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) deleteSeries(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteSeries(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) cancelEvent(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.CancelEvent(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handler) getAdminStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetAdminStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

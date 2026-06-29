package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/subscription-service/internal/domain"
	"github.com/elevatecompact/spark/services/subscription-service/internal/service"
)

type PlanHandler struct {
	svc service.PlanService
}

func NewPlanHandler(svc service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

func (h *PlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	plan, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

func (h *PlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "planId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan id")
		return
	}
	plan, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if err == domain.ErrPlanNotFound {
			writeError(w, http.StatusNotFound, "plan not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

func (h *PlanHandler) List(w http.ResponseWriter, r *http.Request) {
	creatorIDStr := r.URL.Query().Get("creator_id")
	var creatorID *uuid.UUID
	if creatorIDStr != "" {
		id, err := uuid.Parse(creatorIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid creator_id")
			return
		}
		creatorID = &id
	}

	cursorStr := r.URL.Query().Get("cursor")
	var cursor time.Time
	if cursorStr != "" {
		if err := cursor.UnmarshalText([]byte(cursorStr)); err != nil {
			writeError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if n, err := parseInt(limitStr); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	plans, err := h.svc.List(r.Context(), creatorID, cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, plans)
}

func (h *PlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "planId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan id")
		return
	}
	var req domain.UpdatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	plan, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		if err == domain.ErrPlanNotFound {
			writeError(w, http.StatusNotFound, "plan not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

func (h *PlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "planId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid plan id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if err == domain.ErrPlanNotFound {
			writeError(w, http.StatusNotFound, "plan not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

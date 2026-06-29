package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
	"github.com/elevatecompact/spark/services/creator-service/internal/service"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	slots, err := h.scheduleService.GetSchedule(r.Context(), creatorID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": slots,
	})
}

func (h *ScheduleHandler) AddSlot(w http.ResponseWriter, r *http.Request) {
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

	var req domain.CreateScheduleSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DayOfWeek < 0 || req.DayOfWeek > 6 {
		writeError(w, http.StatusBadRequest, "day_of_week must be 0-6")
		return
	}
	if req.StartTime == "" || req.EndTime == "" {
		writeError(w, http.StatusBadRequest, "start_time and end_time are required")
		return
	}

	slot, err := h.scheduleService.AddSlot(r.Context(), creatorID, req)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, slot)
}

func (h *ScheduleHandler) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	slotID, err := uuid.Parse(chi.URLParam(r, "slotID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid slot ID")
		return
	}

	var req domain.CreateScheduleSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.scheduleService.UpdateSlot(r.Context(), slotID, creatorID, req); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *ScheduleHandler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	creatorID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid creator ID")
		return
	}

	slotID, err := uuid.Parse(chi.URLParam(r, "slotID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid slot ID")
		return
	}

	if err := h.scheduleService.DeleteSlot(r.Context(), slotID, creatorID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

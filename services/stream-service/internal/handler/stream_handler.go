package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
	"github.com/elevatecompact/spark/services/stream-service/internal/service"
)

type StreamHandler struct {
	streamSvc *service.StreamService
	rtmpSvc   *service.RTMPService
}

func NewStreamHandler(streamSvc *service.StreamService, rtmpSvc *service.RTMPService) *StreamHandler {
	return &StreamHandler{
		streamSvc: streamSvc,
		rtmpSvc:   rtmpSvc,
	}
}

func (h *StreamHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, ok := GetUserID(r.Context())
	if ok {
		req.CreatorID = userID
	}

	stream, err := h.streamSvc.CreateStream(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create stream")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, stream)
}

func (h *StreamHandler) GetStream(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	stream, err := h.streamSvc.GetStream(r.Context(), id)
	if err != nil {
		if err == domain.ErrStreamNotFound {
			respondError(w, http.StatusNotFound, "stream not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stream)
}

func (h *StreamHandler) UpdateStream(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	var req domain.UpdateStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	stream, err := h.streamSvc.UpdateStream(r.Context(), id, req)
	if err != nil {
		if err == domain.ErrStreamNotFound {
			respondError(w, http.StatusNotFound, "stream not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stream)
}

func (h *StreamHandler) EndStream(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	if err := h.streamSvc.EndStream(r.Context(), id); err != nil {
		if err == domain.ErrStreamNotLive {
			respondError(w, http.StatusConflict, "stream is not live")
			return
		}
		if err == domain.ErrStreamNotFound {
			respondError(w, http.StatusNotFound, "stream not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

func (h *StreamHandler) ListStreams(w http.ResponseWriter, r *http.Request) {
	limit := queryParamInt(r, "limit", 20)
	offset := queryParamInt(r, "offset", 0)

	filter := domain.StreamFilter{
		Limit:  limit,
		Offset: offset,
	}

	if creatorIDStr := r.URL.Query().Get("creator_id"); creatorIDStr != "" {
		cid, err := uuid.Parse(creatorIDStr)
		if err == nil {
			filter.CreatorID = &cid
		}
	}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := domain.StreamStatus(statusStr)
		filter.Status = &status
	}

	filter.Category = r.URL.Query().Get("category")

	streams, err := h.streamSvc.ListStreams(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if streams == nil {
		streams = []domain.Stream{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"streams": streams,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *StreamHandler) ListLiveStreams(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	limit := queryParamInt(r, "limit", 50)
	offset := queryParamInt(r, "offset", 0)

	streams, err := h.streamSvc.ListLiveStreams(r.Context(), category, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if streams == nil {
		streams = []domain.Stream{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"streams": streams,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *StreamHandler) GetStreamHealth(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	stream, err := h.streamSvc.GetStream(r.Context(), id)
	if err != nil {
		if err == domain.ErrStreamNotFound {
			respondError(w, http.StatusNotFound, "stream not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":             stream.ID,
		"status":         stream.Status,
		"viewer_count":   stream.ViewerCount,
		"peak_viewers":   stream.PeakViewers,
		"total_views":    stream.TotalViews,
		"duration":       stream.Duration,
		"bitrate":        stream.Bitrate,
		"frame_rate":     stream.FrameRate,
		"is_live":        stream.IsLive(),
	})
}

func queryParamInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	intVal := 0
	for _, c := range val {
		if c >= '0' && c <= '9' {
			intVal = intVal*10 + int(c-'0')
		} else {
			return defaultVal
		}
	}
	return intVal
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

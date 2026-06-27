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

type WebRTCHandler struct {
	webrtcSvc *service.WebRTCService
	streamSvc *service.StreamService
}

func NewWebRTCHandler(webrtcSvc *service.WebRTCService, streamSvc *service.StreamService) *WebRTCHandler {
	return &WebRTCHandler{
		webrtcSvc: webrtcSvc,
		streamSvc: streamSvc,
	}
}

func (h *WebRTCHandler) HandleOffer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	var offer domain.WebRTCOffer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		respondError(w, http.StatusBadRequest, "invalid offer body")
		return
	}

	userID, _ := GetUserID(r.Context())

	answer, entry, err := h.webrtcSvc.HandleOffer(r.Context(), streamID, offer, userID)
	if err != nil {
		log.Error().Err(err).Str("stream_id", streamID.String()).Msg("WebRTC offer handling failed")
		if err == domain.ErrStreamNotLive {
			respondError(w, http.StatusConflict, "stream is not live")
			return
		}
		if err == domain.ErrViewerLimitReached {
			respondError(w, http.StatusTooManyRequests, "viewer limit reached")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to process offer")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"answer":    answer,
		"pc_id":     entry.ID,
	})
}

func (h *WebRTCHandler) HandleAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	pcID := chi.URLParam(r, "pcId")
	if pcID == "" {
		pcID = r.URL.Query().Get("pc_id")
	}

	var answer domain.WebRTCAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		respondError(w, http.StatusBadRequest, "invalid answer body")
		return
	}

	if err := h.webrtcSvc.HandleAnswer(r.Context(), streamID, pcID, answer); err != nil {
		log.Error().Err(err).Str("stream_id", streamID.String()).Str("pc_id", pcID).Msg("WebRTC answer handling failed")
		if err == domain.ErrPeerConnectionNotFound {
			respondError(w, http.StatusNotFound, "peer connection not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to process answer")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "accepted"})
}

func (h *WebRTCHandler) HandleICECandidate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	pcID := chi.URLParam(r, "pcId")
	if pcID == "" {
		pcID = r.URL.Query().Get("pc_id")
	}

	var candidate domain.ICECandidate
	if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
		respondError(w, http.StatusBadRequest, "invalid ICE candidate body")
		return
	}

	if err := h.webrtcSvc.HandleICECandidate(r.Context(), streamID, pcID, candidate); err != nil {
		log.Error().Err(err).Str("stream_id", streamID.String()).Str("pc_id", pcID).Msg("ICE candidate handling failed")
		if err == domain.ErrPeerConnectionNotFound {
			respondError(w, http.StatusNotFound, "peer connection not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to process ICE candidate")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "accepted"})
}

func (h *WebRTCHandler) GetStream(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	stream, err := h.streamSvc.GetStream(r.Context(), streamID)
	if err != nil {
		if err == domain.ErrStreamNotFound {
			respondError(w, http.StatusNotFound, "stream not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	viewerCount := h.webrtcSvc.GetViewerCount(r.Context(), streamID)

	respondJSON(w, http.StatusOK, domain.WebRTCStreamInfo{
		StreamID:    stream.ID,
		ViewerCount: viewerCount,
		IsLive:      stream.IsLive(),
		Bitrate:     stream.Bitrate,
		Width:       stream.Width,
		Height:      stream.Height,
		FrameRate:   stream.FrameRate,
	})
}

func (h *WebRTCHandler) JoinViewer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	userID, _ := GetUserID(r.Context())

	viewerCount := h.webrtcSvc.GetViewerCount(r.Context(), streamID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"stream_id":    streamID,
		"viewer_id":    userID,
		"viewer_count": viewerCount,
	})
}

func (h *WebRTCHandler) LeaveViewer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	pcID := chi.URLParam(r, "pcId")
	if pcID == "" {
		pcID = r.URL.Query().Get("pc_id")
	}

	if err := h.webrtcSvc.RemoveViewer(r.Context(), streamID, pcID); err != nil {
		log.Warn().Err(err).Str("stream_id", streamID.String()).Str("pc_id", pcID).Msg("Failed to remove viewer")
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "left"})
}

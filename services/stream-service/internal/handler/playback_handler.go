package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
	"github.com/elevatecompact/spark/services/stream-service/internal/service"
)

type PlaybackHandler struct {
	hlsSvc      *service.HLSService
	streamSvc   *service.StreamService
}

func NewPlaybackHandler(hlsSvc *service.HLSService, streamSvc *service.StreamService) *PlaybackHandler {
	return &PlaybackHandler{
		hlsSvc:    hlsSvc,
		streamSvc: streamSvc,
	}
}

func (h *PlaybackHandler) GetMasterPlaylist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	playlist, err := h.hlsSvc.GenerateMasterPlaylist(r.Context(), streamID)
	if err != nil {
		if err == domain.ErrPlaylistNotFound {
			respondError(w, http.StatusNotFound, "playlist not found")
			return
		}
		log.Error().Err(err).Str("stream_id", streamID.String()).Msg("Failed to generate master playlist")
		respondError(w, http.StatusInternalServerError, "failed to generate playlist")
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, playlist)
}

func (h *PlaybackHandler) GetQualityPlaylist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	quality := chi.URLParam(r, "quality")
	if quality == "" {
		quality = "source"
	}

	manifest, err := h.hlsSvc.GetManifest(r.Context(), streamID, quality)
	if err != nil {
		if err == domain.ErrPlaylistNotFound {
			respondError(w, http.StatusNotFound, "playlist not found")
			return
		}
		log.Error().Err(err).Str("stream_id", streamID.String()).Str("quality", quality).Msg("Failed to get manifest")
		respondError(w, http.StatusInternalServerError, "failed to get playlist")
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, manifest)
}

func (h *PlaybackHandler) GetSegment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	quality := chi.URLParam(r, "quality")
	if quality == "" {
		quality = "source"
	}

	segmentStr := chi.URLParam(r, "segment")
	segmentStr = strings.TrimPrefix(segmentStr, "segment-")
	segmentStr = strings.TrimSuffix(segmentStr, ".ts")

	segmentNumber, err := strconv.Atoi(segmentStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid segment number")
		return
	}

	data, err := h.hlsSvc.GetSegment(r.Context(), streamID, quality, segmentNumber)
	if err != nil {
		if err == domain.ErrSegmentNotFound {
			respondError(w, http.StatusNotFound, "segment not found")
			return
		}
		log.Error().Err(err).Str("stream_id", streamID.String()).Str("quality", quality).Int("segment", segmentNumber).Msg("Failed to get segment")
		respondError(w, http.StatusInternalServerError, "failed to get segment")
		return
	}

	w.Header().Set("Content-Type", "video/MP2T")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *PlaybackHandler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stream id")
		return
	}

	timestampStr := chi.URLParam(r, "timestamp")
	timestampStr = strings.TrimSuffix(timestampStr, ".jpg")

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte{})
}

func (h *PlaybackHandler) GetStreamInfo(w http.ResponseWriter, r *http.Request) {
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

	qualities, err := h.hlsSvc.GetAvailableQualities(r.Context(), streamID)
	if err != nil {
		qualities = []string{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":                  stream.ID,
		"title":               stream.Title,
		"status":              stream.Status,
		"is_live":             stream.IsLive(),
		"started_at":          stream.StartedAt,
		"ended_at":            stream.EndedAt,
		"duration":            stream.Duration,
		"available_qualities": qualities,
		"viewer_count":        stream.ViewerCount,
		"thumbnail_url":       stream.ThumbnailURL,
	})
}

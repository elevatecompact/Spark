package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/media-service/internal/domain"
	"github.com/elevatecompact/spark/services/media-service/internal/service"
)

type Handler struct {
	svc *service.MediaService
}

func New(svc *service.MediaService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/upload", func(r chi.Router) {
		r.Post("/init", h.initUpload)
		r.Post("/{id}/chunk", h.uploadChunk)
		r.Post("/{id}/complete", h.completeUpload)
		r.Get("/{id}/status", h.uploadStatus)
		r.Delete("/{id}", h.cancelUpload)
	})
	r.Post("/v1/media/transcode", h.startTranscoding)
	r.Post("/v1/media/thumbnail", h.generateThumbnail)
	r.Post("/v1/media/optimize", h.optimizeImage)
	r.Get("/v1/media/{id}/status", h.mediaStatus)
	r.Get("/v1/media/{id}/playback", h.getPlayback)
	r.Get("/v1/media/{id}/thumbnail/{time}", h.getThumbnail)
	r.Get("/v1/media/{id}/download", h.getDownload)
	r.Get("/v1/media/{id}/info", h.mediaInfo)
	r.Get("/v1/media/{id}/renditions", h.getRenditions)
	r.Post("/v1/drm/license", h.issueLicense)
	r.Post("/v1/drm/policy", h.createDRMPolicy)
	r.Get("/v1/drm/policies", h.getDRMPolicies)
	r.Route("/v1/admin", func(r chi.Router) {
		r.Get("/storage/usage", h.storageUsage)
		r.Post("/cache/purge", h.purgeCDN)
		r.Get("/processing/queue", h.processingQueue)
		r.Post("/media/{id}/retry", h.retryTranscoding)
	})
}

func (h *Handler) initUpload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UploaderID  uuid.UUID `json:"uploaderId"`
		Filename    string    `json:"filename"`
		ContentType string    `json:"contentType"`
		FileSize    int64     `json:"fileSizeBytes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	session, err := h.svc.InitUpload(r.Context(), req.UploaderID, req.Filename, req.ContentType, req.FileSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (h *Handler) uploadChunk(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ ChunkIndex int `json:"chunkIndex"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UploadChunk(r.Context(), id, req.ChunkIndex); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "chunk received"})
}

func (h *Handler) completeUpload(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	asset, err := h.svc.CompleteUpload(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, asset)
}

func (h *Handler) uploadStatus(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	session, err := h.svc.GetUploadStatus(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "upload not found")
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *Handler) cancelUpload(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.CancelUpload(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handler) startTranscoding(w http.ResponseWriter, r *http.Request) {
	var req struct{ MediaID uuid.UUID `json:"mediaId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.StartTranscoding(r.Context(), req.MediaID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "transcoding started"})
}

func (h *Handler) generateThumbnail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MediaID  uuid.UUID `json:"mediaId"`
		TimeSecs float64   `json:"timeSeconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rd, err := h.svc.GenerateThumbnail(r.Context(), req.MediaID, req.TimeSecs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rd)
}

func (h *Handler) optimizeImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MediaID uuid.UUID `json:"mediaId"`
		Width   int       `json:"width"`
		Height  int       `json:"height"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rd, err := h.svc.OptimizeImage(r.Context(), req.MediaID, req.Width, req.Height)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rd)
}

func (h *Handler) mediaStatus(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	media, err := h.svc.GetMediaStatus(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "media not found")
		return
	}
	writeJSON(w, http.StatusOK, media)
}

func (h *Handler) getPlayback(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	url, err := h.svc.GetPlaybackURL(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) getThumbnail(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	timeSecs, _ := strconv.ParseFloat(chi.URLParam(r, "time"), 64)
	url, err := h.svc.GetThumbnailURL(r.Context(), id, timeSecs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) getDownload(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	url, err := h.svc.GetDownloadURL(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) mediaInfo(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	media, err := h.svc.GetMediaInfo(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "media not found")
		return
	}
	writeJSON(w, http.StatusOK, media)
}

func (h *Handler) getRenditions(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	renditions, err := h.svc.GetRenditions(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, renditions)
}

func (h *Handler) issueLicense(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MediaID   uuid.UUID `json:"mediaId"`
		KeySystem string    `json:"keySystem"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	license, err := h.svc.IssueDRMLicense(r.Context(), req.MediaID, req.KeySystem)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"license": license})
}

func (h *Handler) createDRMPolicy(w http.ResponseWriter, r *http.Request) {
	var p domain.DRMPolicy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateDRMPolicy(r.Context(), &p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) getDRMPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.svc.GetDRMPolicies(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, policies)
}

func (h *Handler) storageUsage(w http.ResponseWriter, r *http.Request) {
	usage, err := h.svc.GetStorageUsage(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, usage)
}

func (h *Handler) purgeCDN(w http.ResponseWriter, r *http.Request) {
	var req struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.PurgeCDN(r.Context(), req.Path); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "purged"})
}

func (h *Handler) processingQueue(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	jobs, err := h.svc.GetProcessingQueue(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (h *Handler) retryTranscoding(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.RetryTranscoding(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "retrying"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

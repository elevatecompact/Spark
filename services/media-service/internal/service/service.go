package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/media-service/internal/domain"
	"github.com/elevatecompact/spark/services/media-service/internal/repository"
)

type MediaService struct {
	repo *repository.MediaRepository
	evt  domain.EventProducer
}

func NewMediaService(repo *repository.MediaRepository, evt domain.EventProducer) *MediaService {
	return &MediaService{repo: repo, evt: evt}
}

func (s *MediaService) InitUpload(ctx context.Context, uploaderID uuid.UUID, filename, contentType string, fileSize int64) (*domain.UploadSession, error) {
	chunkSize := int64(5 * 1024 * 1024)
	chunksTotal := int(fileSize / chunkSize)
	if fileSize%chunkSize != 0 {
		chunksTotal++
	}
	session := &domain.UploadSession{
		ID:            uuid.New(),
		UploaderID:    uploaderID,
		Filename:      filename,
		FileSizeBytes: fileSize,
		ContentType:   contentType,
		ChunksTotal:   chunksTotal,
		ChunksDone:    0,
		Status:        "initiated",
		StoragePath:   fmt.Sprintf("uploads/%s/%s", uploaderID.String(), uuid.New().String()),
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.SaveUploadSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *MediaService) UploadChunk(ctx context.Context, sessionID uuid.UUID, chunkIndex int) error {
	session, err := s.repo.GetUploadSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.Status == "completed" {
		return errors.New("upload already completed")
	}
	done := session.ChunksDone + 1
	status := "uploading"
	if done >= session.ChunksTotal {
		status = "completed"
	}
	return s.repo.UpdateUploadSession(ctx, sessionID, done, "", status)
}

func (s *MediaService) CompleteUpload(ctx context.Context, sessionID uuid.UUID) (*domain.MediaAsset, error) {
	session, err := s.repo.GetUploadSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.ChunksDone < session.ChunksTotal {
		return nil, errors.New("upload not complete")
	}
	asset := &domain.MediaAsset{
		ID:             uuid.New(),
		UploaderID:     session.UploaderID,
		ContentType:    classifyContentType(session.ContentType),
		SourceFilename: session.Filename,
		FileSizeBytes:  session.FileSizeBytes,
		MimeType:       session.ContentType,
		Status:         domain.MediaStatusProcessing,
		StoragePath:    session.StoragePath,
		CDNURL:         fmt.Sprintf("https://cdn.spark.dev/media/%s", uuid.New().String()),
		DurationSecs:   0,
		Width:          1920,
		Height:         1080,
		Checksum:       session.Checksum,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.CreateMediaAsset(ctx, asset); err != nil {
		return nil, err
	}
	s.repo.DeleteUploadSession(ctx, sessionID)
	if err := s.StartTranscoding(ctx, asset.ID); err != nil {
		log.Error().Err(err).Str("mediaId", asset.ID.String()).Msg("failed to start transcoding")
	}
	s.evt.Publish(ctx, "media.upload.completed", map[string]interface{}{
		"mediaId": asset.ID, "uploaderId": asset.UploaderID, "fileSize": asset.FileSizeBytes,
	})
	return asset, nil
}

func (s *MediaService) GetUploadStatus(ctx context.Context, sessionID uuid.UUID) (*domain.UploadSession, error) {
	return s.repo.GetUploadSession(ctx, sessionID)
}

func (s *MediaService) CancelUpload(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteUploadSession(ctx, sessionID)
}

func (s *MediaService) StartTranscoding(ctx context.Context, mediaID uuid.UUID) error {
	media, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return err
	}
	if media.ContentType != domain.ContentTypeVideo {
		return nil
	}
	profiles := []domain.RenditionProfile{domain.Rendition720p, domain.Rendition1080p, domain.RenditionSource}
	job := &domain.TranscodingJob{
		ID:        uuid.New(),
		MediaID:   mediaID,
		Profiles:  profiles,
		Status:    domain.TranscodingPending,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateTranscodingJob(ctx, job); err != nil {
		return err
	}
	s.evt.Publish(ctx, "media.transcoding.started", map[string]interface{}{
		"mediaId": mediaID, "jobId": job.ID, "profiles": profiles,
	})
	go s.processTranscodingJob(context.Background(), job.ID)
	return nil
}

func (s *MediaService) processTranscodingJob(ctx context.Context, jobID uuid.UUID) {
	s.repo.UpdateTranscodingJobStatus(ctx, jobID, domain.TranscodingProcessing, "")
	time.Sleep(500 * time.Millisecond)
	job, err := s.repo.GetTranscodingJob(ctx, jobID)
	if err != nil {
		log.Error().Err(err).Str("jobId", jobID.String()).Msg("failed to get transcoding job")
		return
	}
	for _, profile := range job.Profiles {
		format := s.selectFormat(profile)
		cdn := fmt.Sprintf("https://cdn.spark.dev/media/%s/%s", job.MediaID.String(), string(profile))
		rd := &domain.MediaRendition{
			ID:      uuid.New(),
			MediaID: job.MediaID,
			Profile: profile,
			Format:  format,
			FileSizeBytes: int64(rand.Intn(50000000) + 1000000),
			StoragePath: fmt.Sprintf("renditions/%s/%s", job.MediaID.String(), string(profile)),
			CDNURL:  cdn,
			Status:  domain.MediaStatusReady,
			CreatedAt: time.Now(),
		}
		s.repo.CreateRendition(ctx, rd)
	}
	s.repo.UpdateTranscodingJobStatus(ctx, jobID, domain.TranscodingCompleted, "")
	s.repo.UpdateMediaCDN(ctx, job.MediaID, fmt.Sprintf("https://cdn.spark.dev/media/%s/playlist.m3u8", job.MediaID.String()))
	s.evt.Publish(ctx, "media.transcoding.completed", map[string]interface{}{
		"mediaId": job.MediaID, "jobId": job.ID,
		"profiles": job.Profiles,
	})
}

func (s *MediaService) selectFormat(profile domain.RenditionProfile) domain.RenditionFormat {
	switch profile {
	case domain.RenditionThumbnail:
		return domain.RenditionJPG
	case domain.Rendition720p, domain.Rendition1080p:
		return domain.RenditionHLS
	default:
		return domain.RenditionMP4
	}
}

func (s *MediaService) GenerateThumbnail(ctx context.Context, mediaID uuid.UUID, timeSecs float64) (*domain.MediaRendition, error) {
	media, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return nil, err
	}
	if media.ContentType != domain.ContentTypeVideo {
		return nil, errors.New("thumbnails only supported for video")
	}
	rd := &domain.MediaRendition{
		ID:      uuid.New(),
		MediaID: mediaID,
		Profile: domain.RenditionThumbnail,
		Format:  domain.RenditionJPG,
		FileSizeBytes: int64(rand.Intn(50000) + 10000),
		StoragePath: fmt.Sprintf("thumbnails/%s/%f.jpg", mediaID.String(), timeSecs),
		CDNURL:  fmt.Sprintf("https://cdn.spark.dev/media/%s/thumb/%f.jpg", mediaID.String(), timeSecs),
		Status:  domain.MediaStatusReady,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateRendition(ctx, rd); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "media.thumbnail.generated", map[string]interface{}{
		"mediaId": mediaID, "time": timeSecs, "url": rd.CDNURL,
	})
	return rd, nil
}

func (s *MediaService) OptimizeImage(ctx context.Context, mediaID uuid.UUID, width, height int) (*domain.MediaRendition, error) {
	media, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return nil, err
	}
	if media.ContentType != domain.ContentTypeImage {
		return nil, errors.New("optimization only supported for images")
	}
	format := domain.RenditionWebP
	if width == 0 {
		width = media.Width
	}
	if height == 0 {
		height = media.Height
	}
	rd := &domain.MediaRendition{
		ID:      uuid.New(),
		MediaID: mediaID,
		Profile: domain.RenditionProfile(fmt.Sprintf("%dx%d", width, height)),
		Format:  format,
		FileSizeBytes: int64(rand.Intn(100000) + 10000),
		StoragePath: fmt.Sprintf("optimized/%s/%dx%d.webp", mediaID.String(), width, height),
		CDNURL:  fmt.Sprintf("https://cdn.spark.dev/media/%s/opt/%dx%d.webp", mediaID.String(), width, height),
		Status:  domain.MediaStatusReady,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateRendition(ctx, rd); err != nil {
		return nil, err
	}
	return rd, nil
}

func (s *MediaService) GetMediaInfo(ctx context.Context, mediaID uuid.UUID) (*domain.MediaAsset, error) {
	return s.repo.GetMediaAsset(ctx, mediaID)
}

func (s *MediaService) GetMediaStatus(ctx context.Context, mediaID uuid.UUID) (*domain.MediaAsset, error) {
	return s.repo.GetMediaAsset(ctx, mediaID)
}

func (s *MediaService) GetPlaybackURL(ctx context.Context, mediaID uuid.UUID) (string, error) {
	media, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return "", err
	}
	if media.Status != domain.MediaStatusReady {
		return "", errors.New("media not ready")
	}
	return fmt.Sprintf("https://cdn.spark.dev/media/%s/playlist.m3u8", mediaID.String()), nil
}

func (s *MediaService) GetDownloadURL(ctx context.Context, mediaID uuid.UUID) (string, error) {
	_, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://cdn.spark.dev/media/%s/download", mediaID.String()), nil
}

func (s *MediaService) GetThumbnailURL(ctx context.Context, mediaID uuid.UUID, timeSecs float64) (string, error) {
	return fmt.Sprintf("https://cdn.spark.dev/media/%s/thumb/%f.jpg", mediaID.String(), timeSecs), nil
}

func (s *MediaService) IssueDRMLicense(ctx context.Context, mediaID uuid.UUID, keySystem string) (string, error) {
	s.evt.Publish(ctx, "media.drm.license.issued", map[string]interface{}{
		"mediaId": mediaID, "keySystem": keySystem,
	})
	return fmt.Sprintf("drm-license-%s-%s", mediaID.String(), uuid.New().String()), nil
}

func (s *MediaService) CreateDRMPolicy(ctx context.Context, p *domain.DRMPolicy) (*domain.DRMPolicy, error) {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	if err := s.repo.CreateDRMPolicy(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *MediaService) GetDRMPolicies(ctx context.Context) ([]domain.DRMPolicy, error) {
	return s.repo.GetDRMPolicies(ctx)
}

func (s *MediaService) DeleteMedia(ctx context.Context, mediaID uuid.UUID) error {
	if err := s.repo.DeleteMediaAsset(ctx, mediaID); err != nil {
		return err
	}
	s.evt.Publish(ctx, "media.content.deleted", map[string]interface{}{
		"mediaId": mediaID,
	})
	return nil
}

func (s *MediaService) GetStorageUsage(ctx context.Context) (*domain.StorageUsage, error) {
	return s.repo.GetStorageUsage(ctx)
}

func (s *MediaService) PurgeCDN(ctx context.Context, path string) error {
	log.Info().Str("path", path).Msg("CDN cache purged (noop)")
	return nil
}

func (s *MediaService) GetProcessingQueue(ctx context.Context, limit, offset int) ([]domain.TranscodingJob, error) {
	status := domain.TranscodingPending
	return s.repo.ListTranscodingJobs(ctx, &status, limit, offset)
}

func (s *MediaService) RetryTranscoding(ctx context.Context, mediaID uuid.UUID) error {
	return s.StartTranscoding(ctx, mediaID)
}

func (s *MediaService) GetRenditions(ctx context.Context, mediaID uuid.UUID) ([]domain.MediaRendition, error) {
	return s.repo.GetRenditions(ctx, mediaID)
}

func classifyContentType(mime string) domain.ContentType {
	switch {
	case mime == "video/mp4" || mime == "video/webm" || mime == "video/quicktime":
		return domain.ContentTypeVideo
	case mime == "image/jpeg" || mime == "image/png" || mime == "image/webp" || mime == "image/gif":
		return domain.ContentTypeImage
	case mime == "audio/mpeg" || mime == "audio/ogg" || mime == "audio/wav":
		return domain.ContentTypeAudio
	}
	return domain.ContentTypeVideo
}

func (s *MediaService) GetMediaByUploader(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]domain.MediaAsset, error) {
	return s.repo.GetMediaByUploader(ctx, uploaderID, limit, offset)
}

func (s *MediaService) GetTranscodingJob(ctx context.Context, jobID uuid.UUID) (*domain.TranscodingJob, error) {
	return s.repo.GetTranscodingJob(ctx, jobID)
}

func (s *MediaService) IsMediaReady(ctx context.Context, mediaID uuid.UUID) (bool, error) {
	media, err := s.repo.GetMediaAsset(ctx, mediaID)
	if err != nil {
		return false, err
	}
	return media.Status == domain.MediaStatusReady, nil
}

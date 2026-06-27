package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
	"github.com/elevatecompact/spark/services/stream-service/internal/events"
	"github.com/elevatecompact/spark/services/stream-service/internal/repository"
)

type StreamService struct {
	repo      *repository.StreamRepository
	producer  *events.EventProducer
	cfg       *config.Config
	hub       *WebSocketHub

	mu          sync.RWMutex
	liveStreams map[uuid.UUID]*domain.Stream
}

func NewStreamService(repo *repository.StreamRepository, producer *events.EventProducer, cfg *config.Config, hub *WebSocketHub) *StreamService {
	return &StreamService{
		repo:        repo,
		producer:    producer,
		cfg:         cfg,
		hub:         hub,
		liveStreams: make(map[uuid.UUID]*domain.Stream),
	}
}

func generateStreamKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate stream key: %w", err)
	}
	return "spark_" + hex.EncodeToString(b), nil
}

func (s *StreamService) CreateStream(ctx context.Context, req domain.CreateStreamRequest) (*domain.Stream, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrValidation)
	}
	if req.CreatorID == uuid.Nil {
		return nil, fmt.Errorf("%w: creator_id is required", domain.ErrValidation)
	}

	streamKey, err := generateStreamKey()
	if err != nil {
		return nil, err
	}

	rtmpEndpoint := fmt.Sprintf("rtmp://%s/%s", s.cfg.RTMP.Domain, s.cfg.RTMP.AppName)

	stream := &domain.Stream{
		ID:             uuid.New(),
		CreatorID:      req.CreatorID,
		Title:          req.Title,
		Description:    req.Description,
		Category:       req.Category,
		Tags:           req.Tags,
		ThumbnailURL:   "",
		StreamKey:      streamKey,
		RTMPEndpoint:   rtmpEndpoint,
		IngestProtocol: "rtmp",
		Status:         domain.StreamPending,
		RecordEnabled:  req.RecordEnabled,
		ChatEnabled:    req.ChatEnabled,
		AgeRestricted:  req.AgeRestricted,
		DelaySeconds:   req.DelaySeconds,
	}

	if err := s.repo.Create(ctx, stream); err != nil {
		return nil, err
	}

	if err := s.producer.StreamCreated(ctx, stream); err != nil {
		log.Warn().Err(err).Str("stream_id", stream.ID.String()).Msg("Failed to emit StreamCreated event")
	}

	return stream, nil
}

func (s *StreamService) GetStream(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *StreamService) StartStream(ctx context.Context, streamKey string, ingestInfo domain.RTMPStreamInfo) error {
	stream, err := s.repo.GetByStreamKey(ctx, streamKey)
	if err != nil {
		return err
	}

	if !stream.CanBeStarted() {
		if stream.IsLive() {
			return domain.ErrStreamAlreadyLive
		}
		return domain.ErrStreamEnded
	}

	now := time.Now()
	stream.Status = domain.StreamLive
	stream.StartedAt = &now
	stream.Width = ingestInfo.Width
	stream.Height = ingestInfo.Height
	stream.FrameRate = ingestInfo.FrameRate
	stream.Bitrate = ingestInfo.Bitrate

	if err := s.repo.Update(ctx, stream); err != nil {
		return err
	}

	s.mu.Lock()
	s.liveStreams[stream.ID] = stream
	s.mu.Unlock()

	s.hub.BroadcastToStream(stream.ID, WebSocketMessage{
		Type: "stream_started",
		Data: stream,
	})

	if err := s.producer.StreamStarted(ctx, stream); err != nil {
		log.Warn().Err(err).Str("stream_id", stream.ID.String()).Msg("Failed to emit StreamStarted event")
	}

	return nil
}

func (s *StreamService) EndStream(ctx context.Context, id uuid.UUID) error {
	stream, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !stream.IsLive() {
		return domain.ErrStreamNotLive
	}

	if err := s.repo.EndStream(ctx, id); err != nil {
		return err
	}

	endedStream, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.liveStreams, stream.ID)
	s.mu.Unlock()

	s.hub.BroadcastToStream(stream.ID, WebSocketMessage{
		Type: "stream_ended",
		Data: endedStream,
	})

	if err := s.producer.StreamEnded(ctx, endedStream); err != nil {
		log.Warn().Err(err).Str("stream_id", id.String()).Msg("Failed to emit StreamEnded event")
	}

	return nil
}

func (s *StreamService) UpdateStream(ctx context.Context, id uuid.UUID, req domain.UpdateStreamRequest) (*domain.Stream, error) {
	stream, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	var changes []string

	if req.Title != nil {
		stream.Title = *req.Title
		updates["title"] = *req.Title
		changes = append(changes, "title")
	}
	if req.Description != nil {
		stream.Description = *req.Description
		updates["description"] = *req.Description
		changes = append(changes, "description")
	}
	if req.Category != nil {
		stream.Category = *req.Category
		updates["category"] = *req.Category
		changes = append(changes, "category")
	}
	if req.Tags != nil {
		stream.Tags = *req.Tags
		updates["tags"] = *req.Tags
		changes = append(changes, "tags")
	}
	if req.ThumbnailURL != nil {
		stream.ThumbnailURL = *req.ThumbnailURL
		updates["thumbnail_url"] = *req.ThumbnailURL
		changes = append(changes, "thumbnail_url")
	}
	if req.RecordEnabled != nil {
		stream.RecordEnabled = *req.RecordEnabled
		updates["record_enabled"] = *req.RecordEnabled
		changes = append(changes, "record_enabled")
	}
	if req.ChatEnabled != nil {
		stream.ChatEnabled = *req.ChatEnabled
		updates["chat_enabled"] = *req.ChatEnabled
		changes = append(changes, "chat_enabled")
	}
	if req.AgeRestricted != nil {
		stream.AgeRestricted = *req.AgeRestricted
		updates["age_restricted"] = *req.AgeRestricted
		changes = append(changes, "age_restricted")
	}
	if req.DelaySeconds != nil {
		stream.DelaySeconds = *req.DelaySeconds
		updates["delay_seconds"] = *req.DelaySeconds
		changes = append(changes, "delay_seconds")
	}

	if len(updates) == 0 {
		return stream, nil
	}

	if err := s.repo.UpdateSelective(ctx, id, updates); err != nil {
		return nil, err
	}

	if err := s.producer.StreamUpdated(ctx, stream, changes); err != nil {
		log.Warn().Err(err).Str("stream_id", id.String()).Msg("Failed to emit StreamUpdated event")
	}

	return stream, nil
}

func (s *StreamService) ListStreams(ctx context.Context, filter domain.StreamFilter) ([]domain.Stream, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.CreatorID != nil {
		return s.repo.ListByCreator(ctx, *filter.CreatorID, filter.Limit, filter.Offset)
	}
	if filter.Status != nil {
		return s.repo.ListByStatus(ctx, *filter.Status, filter.Limit, filter.Offset)
	}

	return s.repo.ListByStatus(ctx, domain.StreamLive, filter.Limit, filter.Offset)
}

func (s *StreamService) ListLiveStreams(ctx context.Context, category string, limit, offset int) ([]domain.Stream, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListLive(ctx, category, limit, offset)
}

func (s *StreamService) GetLiveCount(ctx context.Context) (int, error) {
	return s.repo.GetLiveCount(ctx)
}

func (s *StreamService) GetActiveStreamByCreator(ctx context.Context, creatorID uuid.UUID) (*domain.Stream, error) {
	streams, err := s.repo.ListByCreator(ctx, creatorID, 1, 0)
	if err != nil {
		return nil, err
	}
	for _, st := range streams {
		if st.IsLive() {
			return &st, nil
		}
	}
	return nil, domain.ErrStreamNotFound
}

func (s *StreamService) ReportStreamIngest(ctx context.Context, streamID uuid.UUID, info domain.RTMPStreamInfo) error {
	updates := make(map[string]interface{})
	if info.Width > 0 {
		updates["width"] = info.Width
	}
	if info.Height > 0 {
		updates["height"] = info.Height
	}
	if info.FrameRate > 0 {
		updates["frame_rate"] = info.FrameRate
	}
	if info.Bitrate > 0 {
		updates["bitrate"] = info.Bitrate
	}
	if info.VideoCodec != "" {
		updates["codec"] = info.VideoCodec
	}
	if info.BytesIn > 0 {
		// update stream stats
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateSelective(ctx, streamID, updates); err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamService) IsStreamLive(ctx context.Context, streamID uuid.UUID) bool {
	s.mu.RLock()
	_, ok := s.liveStreams[streamID]
	s.mu.RUnlock()
	if ok {
		return true
	}
	stream, err := s.repo.GetByID(ctx, streamID)
	if err != nil {
		return false
	}
	return stream.IsLive()
}

func (s *StreamService) UpdateStreamAfterTranscoding(ctx context.Context, streamID uuid.UUID, qualities []string) error {
	updates := map[string]interface{}{
		"available_qualities": qualities,
	}
	return s.repo.UpdateSelective(ctx, streamID, updates)
}

func (s *StreamService) RecordStream(ctx context.Context, streamID uuid.UUID) error {
	stream, err := s.repo.GetByID(ctx, streamID)
	if err != nil {
		return err
	}

	if !stream.RecordEnabled {
		return nil
	}

	return s.startRecording(ctx, stream.ID)
}

func (s *StreamService) startRecording(ctx context.Context, streamID uuid.UUID) error {
	log.Info().Str("stream_id", streamID.String()).Msg("Starting stream recording")
	return nil
}

func findStream(streams []domain.Stream, id uuid.UUID) *domain.Stream {
	for _, s := range streams {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

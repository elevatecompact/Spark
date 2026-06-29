package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
	"github.com/elevatecompact/spark/services/stream-service/internal/repository"
)

type RTMPService struct {
	repo         *repository.StreamRepository
	streamSvc    *StreamService
	publishCache map[string]*domain.RTMPStreamInfo
	mu           sync.RWMutex
}

func NewRTMPService(repo *repository.StreamRepository, streamSvc *StreamService) *RTMPService {
	return &RTMPService{
		repo:         repo,
		streamSvc:    streamSvc,
		publishCache: make(map[string]*domain.RTMPStreamInfo),
	}
}

func (s *RTMPService) ValidateStreamKey(ctx context.Context, streamKey string) (*domain.Stream, error) {
	if streamKey == "" {
		return nil, domain.ErrStreamKeyInvalid
	}

	return s.repo.GetByStreamKey(ctx, streamKey)
}

func (s *RTMPService) HandleRTMPPublish(ctx context.Context, req domain.RTMPPublishRequest) error {
	if req.StreamKey == "" {
		return fmt.Errorf("%w: stream_key is required", domain.ErrValidation)
	}

	stream, err := s.repo.GetByStreamKey(ctx, req.StreamKey)
	if err != nil {
		return fmt.Errorf("validate publish: %w", err)
	}

	ingestInfo := domain.RTMPStreamInfo{
		StreamID:    stream.ID,
		StreamKey:   req.StreamKey,
		CreatorID:   stream.CreatorID,
		Connected:   true,
		ConnectedAt: time.Now(),
		BytesIn:     0,
		BytesOut:    0,
	}

	s.mu.Lock()
	s.publishCache[req.StreamKey] = &ingestInfo
	s.mu.Unlock()

	if err := s.streamSvc.StartStream(ctx, req.StreamKey, ingestInfo); err != nil {
		s.mu.Lock()
		delete(s.publishCache, req.StreamKey)
		s.mu.Unlock()
		return fmt.Errorf("start stream: %w", err)
	}

	log.Info().
		Str("stream_id", stream.ID.String()).
		Str("stream_key", req.StreamKey).
		Str("ip", req.IP).
		Msg("RTMP publish started")

	return nil
}

func (s *RTMPService) HandleRTMPUnpublish(ctx context.Context, req domain.RTMPUnpublishRequest) error {
	s.mu.Lock()
	info, ok := s.publishCache[req.StreamKey]
	delete(s.publishCache, req.StreamKey)
	s.mu.Unlock()

	if !ok {
		log.Warn().Str("stream_key", req.StreamKey).Msg("RTMP unpublish for unknown stream")
		return domain.ErrStreamNotFound
	}

	if err := s.streamSvc.EndStream(ctx, info.StreamID); err != nil {
		return fmt.Errorf("end stream on unpublish: %w", err)
	}

	log.Info().
		Str("stream_id", info.StreamID.String()).
		Str("stream_key", req.StreamKey).
		Str("reason", req.Reason).
		Msg("RTMP publish ended")

	return nil
}

func (s *RTMPService) UpdateRTMPStreamInfo(ctx context.Context, streamKey string, info domain.RTMPStreamInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.publishCache[streamKey]
	if !ok {
		return domain.ErrStreamNotFound
	}

	if info.BytesIn > 0 {
		existing.BytesIn = info.BytesIn
	}
	if info.BytesOut > 0 {
		existing.BytesOut = info.BytesOut
	}
	if info.Bitrate > 0 {
		existing.Bitrate = info.Bitrate
	}
	if info.AudioCodec != "" {
		existing.AudioCodec = info.AudioCodec
	}
	if info.VideoCodec != "" {
		existing.VideoCodec = info.VideoCodec
	}
	if info.Width > 0 {
		existing.Width = info.Width
	}
	if info.Height > 0 {
		existing.Height = info.Height
	}
	if info.FrameRate > 0 {
		existing.FrameRate = info.FrameRate
	}

	if err := s.streamSvc.ReportStreamIngest(ctx, existing.StreamID, *existing); err != nil {
		log.Warn().Err(err).Str("stream_key", streamKey).Msg("Failed to update stream ingest info")
	}

	return nil
}

func (s *RTMPService) GetRTMPStats(ctx context.Context) ([]domain.RTMPStreamInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make([]domain.RTMPStreamInfo, 0, len(s.publishCache))
	for _, info := range s.publishCache {
		stats = append(stats, *info)
	}
	return stats, nil
}

func (s *RTMPService) GetActivePublishInfo(ctx context.Context, streamKey string) (*domain.RTMPStreamInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, ok := s.publishCache[streamKey]
	if !ok {
		return nil, domain.ErrStreamNotFound
	}
	return info, nil
}

func (s *RTMPService) ListActivePublishKeys(ctx context.Context) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.publishCache))
	for k := range s.publishCache {
		keys = append(keys, k)
	}
	return keys
}

func (s *RTMPService) IsKeyPublishing(ctx context.Context, streamKey string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.publishCache[streamKey]
	return ok
}

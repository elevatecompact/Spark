package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
	"github.com/elevatecompact/spark/services/stream-service/internal/events"
	"github.com/elevatecompact/spark/services/stream-service/internal/repository"
)

type RecordingService struct {
	cfg      *config.Config
	producer *events.EventProducer
	repo     *repository.StreamRepository

	mu             sync.RWMutex
	activeRecordings map[uuid.UUID]*ActiveRecording
}

type ActiveRecording struct {
	RecordingID uuid.UUID
	StreamID    uuid.UUID
	CreatorID   uuid.UUID
	StartedAt   time.Time
	FilePath    string
	SegmentCount int
	BytesWritten int64
}

func NewRecordingService(cfg *config.Config, producer *events.EventProducer, repo *repository.StreamRepository) *RecordingService {
	return &RecordingService{
		cfg:              cfg,
		producer:         producer,
		repo:             repo,
		activeRecordings: make(map[uuid.UUID]*ActiveRecording),
	}
}

func (s *RecordingService) StartRecording(ctx context.Context, streamID uuid.UUID) error {
	stream, err := s.repo.GetByID(ctx, streamID)
	if err != nil {
		return err
	}

	if !stream.RecordEnabled {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeRecordings[streamID]; exists {
		return nil
	}

	recordingID := uuid.New()
	s3Key := fmt.Sprintf("recordings/%s/%s/%d.mp4", stream.CreatorID, streamID, time.Now().Unix())

	activeRec := &ActiveRecording{
		RecordingID: recordingID,
		StreamID:    streamID,
		CreatorID:   stream.CreatorID,
		StartedAt:   time.Now(),
		FilePath:    s3Key,
	}

	s.activeRecordings[streamID] = activeRec

	recordingIDPtr := recordingID
	updates := map[string]interface{}{
		"recording_id": recordingIDPtr,
	}

	if err := s.repo.UpdateSelective(ctx, streamID, updates); err != nil {
		log.Warn().Err(err).Str("stream_id", streamID.String()).Msg("Failed to update stream recording_id")
	}

	if err := s.producer.RecordingStarted(ctx, streamID, stream.CreatorID, recordingID); err != nil {
		log.Warn().Err(err).Str("recording_id", recordingID.String()).Msg("Failed to emit RecordingStarted event")
	}

	log.Info().
		Str("recording_id", recordingID.String()).
		Str("stream_id", streamID.String()).
		Str("s3_key", s3Key).
		Msg("Recording started")

	return nil
}

func (s *RecordingService) StopRecording(ctx context.Context, streamID uuid.UUID) error {
	s.mu.Lock()
	activeRec, ok := s.activeRecordings[streamID]
	delete(s.activeRecordings, streamID)
	s.mu.Unlock()

	if !ok {
		return nil
	}

	duration := int(time.Since(activeRec.StartedAt).Seconds())

	if err := s.producer.RecordingCompleted(ctx, streamID, activeRec.CreatorID, activeRec.RecordingID, activeRec.FilePath); err != nil {
		log.Warn().Err(err).Str("recording_id", activeRec.RecordingID.String()).Msg("Failed to emit RecordingCompleted event")
	}

	log.Info().
		Str("recording_id", activeRec.RecordingID.String()).
		Str("stream_id", streamID.String()).
		Int("duration", duration).
		Int64("bytes", activeRec.BytesWritten).
		Msg("Recording completed")

	return nil
}

func (s *RecordingService) GetRecording(ctx context.Context, recordingID uuid.UUID) (*domain.Recording, error) {
	return nil, fmt.Errorf("%w: recording_id=%s", domain.ErrRecordingNotFound, recordingID.String())
}

func (s *RecordingService) ListRecordings(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]domain.Recording, error) {
	return []domain.Recording{}, nil
}

func (s *RecordingService) DeleteRecording(ctx context.Context, recordingID uuid.UUID) error {
	return nil
}

func (s *RecordingService) IsRecording(ctx context.Context, streamID uuid.UUID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.activeRecordings[streamID]
	return ok
}

func (s *RecordingService) GetActiveRecording(ctx context.Context, streamID uuid.UUID) (*ActiveRecording, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rec, ok := s.activeRecordings[streamID]
	if !ok {
		return nil, domain.ErrRecordingNotFound
	}
	return rec, nil
}

func (s *RecordingService) ListActiveRecordings(ctx context.Context) []uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]uuid.UUID, 0, len(s.activeRecordings))
	for sid := range s.activeRecordings {
		ids = append(ids, sid)
	}
	return ids
}

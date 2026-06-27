package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/domain"
)

type StreamRepository struct {
	pool *pgxpool.Pool
}

func NewStreamRepository(pool *pgxpool.Pool) *StreamRepository {
	return &StreamRepository{pool: pool}
}

func (r *StreamRepository) Create(ctx context.Context, s *domain.Stream) error {
	query := `
		INSERT INTO streams (
			id, creator_id, title, description, category, tags, thumbnail_url,
			stream_key, rtmp_endpoint, ingest_protocol, status,
			width, height, frame_rate, bitrate, codec, available_qualities,
			viewer_count, peak_viewers, total_views,
			record_enabled, chat_enabled, age_restricted, delay_seconds,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17,
			$18, $19, $20,
			$21, $22, $23, $24,
			$25, $26
		)`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query,
		s.ID, s.CreatorID, s.Title, s.Description, s.Category, s.Tags, s.ThumbnailURL,
		s.StreamKey, s.RTMPEndpoint, s.IngestProtocol, string(s.Status),
		s.Width, s.Height, s.FrameRate, s.Bitrate, s.Codec, s.AvailableQualities,
		s.ViewerCount, s.PeakViewers, s.TotalViews,
		s.RecordEnabled, s.ChatEnabled, s.AgeRestricted, s.DelaySeconds,
		now, now,
	)
	if err != nil {
		return fmt.Errorf("insert stream: %w", err)
	}
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

func (r *StreamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	query := `SELECT
		id, creator_id, title, description, category, tags, thumbnail_url,
		stream_key, rtmp_endpoint, ingest_protocol, status,
		started_at, ended_at, duration,
		width, height, frame_rate, bitrate, codec, available_qualities,
		viewer_count, peak_viewers, total_views,
		record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
		created_at, updated_at
		FROM streams WHERE id = $1`

	s := &domain.Stream{}
	var startedAt, endedAt *time.Time
	var recordingID *uuid.UUID
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.CreatorID, &s.Title, &s.Description, &s.Category, &s.Tags, &s.ThumbnailURL,
		&s.StreamKey, &s.RTMPEndpoint, &s.IngestProtocol, (*string)(&s.Status),
		&startedAt, &endedAt, &s.Duration,
		&s.Width, &s.Height, &s.FrameRate, &s.Bitrate, &s.Codec, &s.AvailableQualities,
		&s.ViewerCount, &s.PeakViewers, &s.TotalViews,
		&s.RecordEnabled, &recordingID, &s.ChatEnabled, &s.AgeRestricted, &s.DelaySeconds,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrStreamNotFound
		}
		return nil, fmt.Errorf("get stream by id: %w", err)
	}
	s.StartedAt = startedAt
	s.EndedAt = endedAt
	s.RecordingID = recordingID
	return s, nil
}

func (r *StreamRepository) GetByStreamKey(ctx context.Context, streamKey string) (*domain.Stream, error) {
	query := `SELECT
		id, creator_id, title, description, category, tags, thumbnail_url,
		stream_key, rtmp_endpoint, ingest_protocol, status,
		started_at, ended_at, duration,
		width, height, frame_rate, bitrate, codec, available_qualities,
		viewer_count, peak_viewers, total_views,
		record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
		created_at, updated_at
		FROM streams WHERE stream_key = $1`

	s := &domain.Stream{}
	var startedAt, endedAt *time.Time
	var recordingID *uuid.UUID
	err := r.pool.QueryRow(ctx, query, streamKey).Scan(
		&s.ID, &s.CreatorID, &s.Title, &s.Description, &s.Category, &s.Tags, &s.ThumbnailURL,
		&s.StreamKey, &s.RTMPEndpoint, &s.IngestProtocol, (*string)(&s.Status),
		&startedAt, &endedAt, &s.Duration,
		&s.Width, &s.Height, &s.FrameRate, &s.Bitrate, &s.Codec, &s.AvailableQualities,
		&s.ViewerCount, &s.PeakViewers, &s.TotalViews,
		&s.RecordEnabled, &recordingID, &s.ChatEnabled, &s.AgeRestricted, &s.DelaySeconds,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrStreamKeyInvalid
		}
		return nil, fmt.Errorf("get stream by key: %w", err)
	}
	s.StartedAt = startedAt
	s.EndedAt = endedAt
	s.RecordingID = recordingID
	return s, nil
}

func (r *StreamRepository) Update(ctx context.Context, s *domain.Stream) error {
	query := `
		UPDATE streams SET
			title = $2, description = $3, category = $4, tags = $5,
			thumbnail_url = $6, ingest_protocol = $7, status = $8,
			started_at = $9, ended_at = $10, duration = $11,
			width = $12, height = $13, frame_rate = $14, bitrate = $15,
			codec = $16, available_qualities = $17,
			viewer_count = $18, peak_viewers = $19, total_views = $20,
			record_enabled = $21, recording_id = $22,
			chat_enabled = $23, age_restricted = $24, delay_seconds = $25,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		s.ID, s.Title, s.Description, s.Category, s.Tags, s.ThumbnailURL,
		s.IngestProtocol, string(s.Status),
		s.StartedAt, s.EndedAt, s.Duration,
		s.Width, s.Height, s.FrameRate, s.Bitrate, s.Codec, s.AvailableQualities,
		s.ViewerCount, s.PeakViewers, s.TotalViews,
		s.RecordEnabled, s.RecordingID,
		s.ChatEnabled, s.AgeRestricted, s.DelaySeconds,
	)
	if err != nil {
		return fmt.Errorf("update stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) UpdateSelective(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates)+1)
	args := make([]interface{}, 0, len(updates)+1)
	idx := 1

	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
		args = append(args, val)
		idx++
	}
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", idx))
	args = append(args, time.Now())
	idx++
	args = append(args, id)

	query := fmt.Sprintf("UPDATE streams SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), idx)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("selective update stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) ListByStatus(ctx context.Context, status domain.StreamStatus, limit, offset int) ([]domain.Stream, error) {
	query := `SELECT
		id, creator_id, title, description, category, tags, thumbnail_url,
		stream_key, rtmp_endpoint, ingest_protocol, status,
		started_at, ended_at, duration,
		width, height, frame_rate, bitrate, codec, available_qualities,
		viewer_count, peak_viewers, total_views,
		record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
		created_at, updated_at
		FROM streams WHERE status = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	return r.scanStreams(ctx, query, string(status), limit, offset)
}

func (r *StreamRepository) ListByCreator(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]domain.Stream, error) {
	query := `SELECT
		id, creator_id, title, description, category, tags, thumbnail_url,
		stream_key, rtmp_endpoint, ingest_protocol, status,
		started_at, ended_at, duration,
		width, height, frame_rate, bitrate, codec, available_qualities,
		viewer_count, peak_viewers, total_views,
		record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
		created_at, updated_at
		FROM streams WHERE creator_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	return r.scanStreams(ctx, query, creatorID, limit, offset)
}

func (r *StreamRepository) ListLive(ctx context.Context, category string, limit, offset int) ([]domain.Stream, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT
			id, creator_id, title, description, category, tags, thumbnail_url,
			stream_key, rtmp_endpoint, ingest_protocol, status,
			started_at, ended_at, duration,
			width, height, frame_rate, bitrate, codec, available_qualities,
			viewer_count, peak_viewers, total_views,
			record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
			created_at, updated_at
			FROM streams WHERE status = 'live' AND category = $1
			ORDER BY viewer_count DESC LIMIT $2 OFFSET $3`
		args = []interface{}{category, limit, offset}
	} else {
		query = `SELECT
			id, creator_id, title, description, category, tags, thumbnail_url,
			stream_key, rtmp_endpoint, ingest_protocol, status,
			started_at, ended_at, duration,
			width, height, frame_rate, bitrate, codec, available_qualities,
			viewer_count, peak_viewers, total_views,
			record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
			created_at, updated_at
			FROM streams WHERE status = 'live'
			ORDER BY viewer_count DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}

	return r.scanStreams(ctx, query, args...)
}

func (r *StreamRepository) GetLiveCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM streams WHERE status = 'live'`
	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get live count: %w", err)
	}
	return count, nil
}

func (r *StreamRepository) UpdateViewerCount(ctx context.Context, id uuid.UUID, count int) error {
	query := `UPDATE streams SET viewer_count = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, count)
	if err != nil {
		return fmt.Errorf("update viewer count: %w", err)
	}
	return nil
}

func (r *StreamRepository) IncrementTotalViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE streams SET total_views = total_views + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("increment total views: %w", err)
	}
	return nil
}

func (r *StreamRepository) EndStream(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE streams SET status = 'ended', ended_at = $2, duration = EXTRACT(EPOCH FROM ($2 - started_at))::INT, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("end stream: %w", err)
	}
	return nil
}

func (r *StreamRepository) Search(ctx context.Context, query_str string, limit, offset int) ([]domain.Stream, error) {
	query := `SELECT
		id, creator_id, title, description, category, tags, thumbnail_url,
		stream_key, rtmp_endpoint, ingest_protocol, status,
		started_at, ended_at, duration,
		width, height, frame_rate, bitrate, codec, available_qualities,
		viewer_count, peak_viewers, total_views,
		record_enabled, recording_id, chat_enabled, age_restricted, delay_seconds,
		created_at, updated_at
		FROM streams WHERE title ILIKE $1 OR description ILIKE $1
		ORDER BY viewer_count DESC LIMIT $2 OFFSET $3`

	searchPattern := "%" + query_str + "%"
	return r.scanStreams(ctx, query, searchPattern, limit, offset)
}

func (r *StreamRepository) scanStreams(ctx context.Context, query string, args ...interface{}) ([]domain.Stream, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query streams: %w", err)
	}
	defer rows.Close()

	var streams []domain.Stream
	for rows.Next() {
		var s domain.Stream
		var startedAt, endedAt *time.Time
		var recordingID *uuid.UUID
		err := rows.Scan(
			&s.ID, &s.CreatorID, &s.Title, &s.Description, &s.Category, &s.Tags, &s.ThumbnailURL,
			&s.StreamKey, &s.RTMPEndpoint, &s.IngestProtocol, (*string)(&s.Status),
			&startedAt, &endedAt, &s.Duration,
			&s.Width, &s.Height, &s.FrameRate, &s.Bitrate, &s.Codec, &s.AvailableQualities,
			&s.ViewerCount, &s.PeakViewers, &s.TotalViews,
			&s.RecordEnabled, &recordingID, &s.ChatEnabled, &s.AgeRestricted, &s.DelaySeconds,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("scan stream row")
			continue
		}
		s.StartedAt = startedAt
		s.EndedAt = endedAt
		s.RecordingID = recordingID
		streams = append(streams, s)
	}

	if streams == nil {
		streams = []domain.Stream{}
	}
	return streams, nil
}



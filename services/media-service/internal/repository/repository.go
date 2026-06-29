package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/media-service/internal/domain"
)

type MediaRepository struct {
	pool *pgxpool.Pool
}

func NewMediaRepository(pool *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{pool: pool}
}

func (r *MediaRepository) CreateMediaAsset(ctx context.Context, m *domain.MediaAsset) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO media_assets (id, uploader_id, content_type, source_filename, file_size_bytes, mime_type, status, storage_path, cdn_url, duration_seconds, width, height, checksum, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`, m.ID, m.UploaderID, string(m.ContentType), m.SourceFilename, m.FileSizeBytes, m.MimeType, string(m.Status), m.StoragePath, m.CDNURL, m.DurationSecs, m.Width, m.Height, m.Checksum, m.CreatedAt)
	return err
}

func (r *MediaRepository) GetMediaAsset(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, uploader_id, content_type, source_filename, file_size_bytes, mime_type, status, storage_path, cdn_url, duration_seconds, width, height, checksum, created_at
		FROM media_assets WHERE id=$1
	`, id)
	m := &domain.MediaAsset{}
	err := row.Scan(&m.ID, &m.UploaderID, &m.ContentType, &m.SourceFilename, &m.FileSizeBytes, &m.MimeType, &m.Status, &m.StoragePath, &m.CDNURL, &m.DurationSecs, &m.Width, &m.Height, &m.Checksum, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MediaRepository) UpdateMediaStatus(ctx context.Context, id uuid.UUID, status domain.MediaStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE media_assets SET status=$2 WHERE id=$1`, id, string(status))
	return err
}

func (r *MediaRepository) UpdateMediaCDN(ctx context.Context, id uuid.UUID, cdnURL string) error {
	_, err := r.pool.Exec(ctx, `UPDATE media_assets SET cdn_url=$2, status='ready' WHERE id=$1`, id, cdnURL)
	return err
}

func (r *MediaRepository) DeleteMediaAsset(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE media_assets SET status='deleted' WHERE id=$1`, id)
	return err
}

func (r *MediaRepository) CreateRendition(ctx context.Context, rd *domain.MediaRendition) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO media_renditions (id, media_id, profile, format, file_size_bytes, storage_path, cdn_url, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, rd.ID, rd.MediaID, string(rd.Profile), string(rd.Format), rd.FileSizeBytes, rd.StoragePath, rd.CDNURL, string(rd.Status), rd.CreatedAt)
	return err
}

func (r *MediaRepository) GetRenditions(ctx context.Context, mediaID uuid.UUID) ([]domain.MediaRendition, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, media_id, profile, format, file_size_bytes, storage_path, cdn_url, status, created_at
		FROM media_renditions WHERE media_id=$1 ORDER BY created_at
	`, mediaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.MediaRendition
	for rows.Next() {
		var rd domain.MediaRendition
		if err := rows.Scan(&rd.ID, &rd.MediaID, &rd.Profile, &rd.Format, &rd.FileSizeBytes, &rd.StoragePath, &rd.CDNURL, &rd.Status, &rd.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rd)
	}
	return res, nil
}

func (r *MediaRepository) CreateTranscodingJob(ctx context.Context, j *domain.TranscodingJob) error {
	profiles, _ := json.Marshal(j.Profiles)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO transcoding_jobs (id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, j.ID, j.MediaID, profiles, string(j.Status), j.WorkerID, j.StartedAt, j.CompletedAt, j.ErrorMessage, j.CreatedAt)
	return err
}

func (r *MediaRepository) GetTranscodingJob(ctx context.Context, id uuid.UUID) (*domain.TranscodingJob, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at
		FROM transcoding_jobs WHERE id=$1
	`, id)
	j := &domain.TranscodingJob{}
	var profiles []byte
	err := row.Scan(&j.ID, &j.MediaID, &profiles, &j.Status, &j.WorkerID, &j.StartedAt, &j.CompletedAt, &j.ErrorMessage, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(profiles, &j.Profiles)
	return j, nil
}

func (r *MediaRepository) UpdateTranscodingJobStatus(ctx context.Context, id uuid.UUID, status domain.TranscodingJobStatus, errMsg string) error {
	_, err := r.pool.Exec(ctx, `UPDATE transcoding_jobs SET status=$2, error_message=$3, completed_at=CASE WHEN $2 IN ('completed','failed') THEN NOW() ELSE NULL END WHERE id=$1`, id, string(status), errMsg)
	return err
}

func (r *MediaRepository) ListTranscodingJobs(ctx context.Context, status *domain.TranscodingJobStatus, limit, offset int) ([]domain.TranscodingJob, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if status != nil {
		where += " AND status=$" + string(rune('0'+idx))
		args = append(args, string(*status))
		idx++
	}
	args = append(args, limit, offset)
	q := `SELECT id, media_id, profiles, status, worker_id, started_at, completed_at, error_message, created_at
	      FROM transcoding_jobs ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.TranscodingJob
	for rows.Next() {
		var j domain.TranscodingJob
		var profiles []byte
		if err := rows.Scan(&j.ID, &j.MediaID, &profiles, &j.Status, &j.WorkerID, &j.StartedAt, &j.CompletedAt, &j.ErrorMessage, &j.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(profiles, &j.Profiles)
		res = append(res, j)
	}
	return res, nil
}

func (r *MediaRepository) CreateDRMPolicy(ctx context.Context, p *domain.DRMPolicy) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO drm_policies (id, name, content_id, key_system, license_duration_seconds, security_level, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, p.ID, p.Name, p.ContentID, string(p.KeySystem), p.LicenseDurationSecs, p.SecurityLevel, p.IsActive, p.CreatedAt)
	return err
}

func (r *MediaRepository) GetDRMPolicies(ctx context.Context) ([]domain.DRMPolicy, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, content_id, key_system, license_duration_seconds, security_level, is_active, created_at FROM drm_policies WHERE is_active=true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.DRMPolicy
	for rows.Next() {
		var p domain.DRMPolicy
		if err := rows.Scan(&p.ID, &p.Name, &p.ContentID, &p.KeySystem, &p.LicenseDurationSecs, &p.SecurityLevel, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *MediaRepository) GetStorageUsage(ctx context.Context) (*domain.StorageUsage, error) {
	var u domain.StorageUsage
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(file_size_bytes),0) FROM media_assets WHERE status!='deleted'`).Scan(&u.TotalBytes)
	if err != nil {
		return nil, err
	}
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM media_assets WHERE status!='deleted'`).Scan(&u.AssetCount)
	return &u, nil
}

func (r *MediaRepository) CheckMediaExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM media_assets WHERE id=$1 AND status!='deleted')`, id).Scan(&exists)
	return exists, err
}

func (r *MediaRepository) GetMediaByUploader(ctx context.Context, uploaderID uuid.UUID, limit, offset int) ([]domain.MediaAsset, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, uploader_id, content_type, source_filename, file_size_bytes, mime_type, status, storage_path, cdn_url, duration_seconds, width, height, checksum, created_at
		FROM media_assets WHERE uploader_id=$1 AND status!='deleted' ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, uploaderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.MediaAsset
	for rows.Next() {
		var m domain.MediaAsset
		if err := rows.Scan(&m.ID, &m.UploaderID, &m.ContentType, &m.SourceFilename, &m.FileSizeBytes, &m.MimeType, &m.Status, &m.StoragePath, &m.CDNURL, &m.DurationSecs, &m.Width, &m.Height, &m.Checksum, &m.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func (r *MediaRepository) SaveUploadSession(ctx context.Context, s *domain.UploadSession) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO upload_sessions (id, uploader_id, filename, file_size_bytes, content_type, chunks_total, chunks_done, checksum, status, storage_path, created_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, s.ID, s.UploaderID, s.Filename, s.FileSizeBytes, s.ContentType, s.ChunksTotal, s.ChunksDone, s.Checksum, s.Status, s.StoragePath, s.CreatedAt, s.ExpiresAt)
	return err
}

func (r *MediaRepository) GetUploadSession(ctx context.Context, id uuid.UUID) (*domain.UploadSession, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, uploader_id, filename, file_size_bytes, content_type, chunks_total, chunks_done, checksum, status, storage_path, created_at, expires_at
		FROM upload_sessions WHERE id=$1
	`, id)
	s := &domain.UploadSession{}
	err := row.Scan(&s.ID, &s.UploaderID, &s.Filename, &s.FileSizeBytes, &s.ContentType, &s.ChunksTotal, &s.ChunksDone, &s.Checksum, &s.Status, &s.StoragePath, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *MediaRepository) UpdateUploadSession(ctx context.Context, id uuid.UUID, chunksDone int, checksum, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE upload_sessions SET chunks_done=$2, checksum=CASE WHEN $3='' THEN checksum ELSE $3 END, status=$4 WHERE id=$1`, id, chunksDone, checksum, status)
	return err
}

func (r *MediaRepository) DeleteUploadSession(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM upload_sessions WHERE id=$1`, id)
	return err
}

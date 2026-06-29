package repository

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/translation-service/internal/domain"
)

type TranslationRepository interface {
	SaveMemory(ctx context.Context, e *domain.TranslationMemoryEntry) error
	LookupMemory(ctx context.Context, sourceHash, sourceLang, targetLang string) (*domain.TranslationMemoryEntry, error)
	DeleteMemory(ctx context.Context, id uuid.UUID) error
	CreateJob(ctx context.Context, j *domain.TranslationJob) error
	UpdateJobStatus(ctx context.Context, id uuid.UUID, status domain.TranslationJobStatus) error
	ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewEntry, error)
	ApproveReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID, correctedText string) error
	RejectReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error
	LogUsage(ctx context.Context, provider string, charCount int) error
	GetUsage(ctx context.Context) ([]domain.ProviderUsage, error)
}

type transRepo struct{ pool *pgxpool.Pool }

func NewTranslationRepository(pool *pgxpool.Pool) TranslationRepository {
	return &transRepo{pool}
}

func (r *transRepo) SaveMemory(ctx context.Context, e *domain.TranslationMemoryEntry) error {
	e.SourceHash = fmt.Sprintf("%x", sha256.Sum256([]byte(e.SourceText)))
	_, err := r.pool.Exec(ctx,
		`INSERT INTO translation_memory (source_hash, source_text, translated_text, source_lang, target_lang, provider, quality_score, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		 ON CONFLICT (source_hash, source_lang, target_lang) DO UPDATE SET translated_text=$3, provider=$6, quality_score=$7, updated_at=NOW()`,
		e.SourceHash, e.SourceText, e.TranslatedText, e.SourceLang, e.TargetLang, e.Provider, e.QualityScore)
	return err
}

func (r *transRepo) LookupMemory(ctx context.Context, sourceHash, sourceLang, targetLang string) (*domain.TranslationMemoryEntry, error) {
	e := &domain.TranslationMemoryEntry{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, source_hash, source_text, translated_text, source_lang, target_lang, provider, quality_score, created_at, updated_at
		 FROM translation_memory WHERE source_hash=$1 AND source_lang=$2 AND target_lang=$3`,
		sourceHash, sourceLang, targetLang).Scan(&e.ID, &e.SourceHash, &e.SourceText, &e.TranslatedText,
		&e.SourceLang, &e.TargetLang, &e.Provider, &e.QualityScore, &e.CreatedAt, &e.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return e, err
}

func (r *transRepo) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM translation_memory WHERE id=$1`, id)
	return err
}

func (r *transRepo) CreateJob(ctx context.Context, j *domain.TranslationJob) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO translation_jobs (id, content_type, content_id, status, languages, created_at) VALUES ($1,$2,$3,$4,$5,NOW())`,
		j.ID, j.ContentType, j.ContentID, j.Status, j.Languages)
	return err
}

func (r *transRepo) UpdateJobStatus(ctx context.Context, id uuid.UUID, status domain.TranslationJobStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE translation_jobs SET status=$2 WHERE id=$1`, id, status)
	return err
}

func (r *transRepo) ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, translation_id, original_text, translated_text, source_lang, target_lang, reviewer_id, status, corrected_text, reviewed_at
		 FROM review_queue WHERE status=$1 ORDER BY reviewed_at NULLS FIRST`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []domain.ReviewEntry
	for rows.Next() {
		e := domain.ReviewEntry{}
		if err := rows.Scan(&e.ID, &e.TranslationID, &e.OriginalText, &e.TranslatedText, &e.SourceLang, &e.TargetLang, &e.ReviewerID, &e.Status, &e.CorrectedText, &e.ReviewedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []domain.ReviewEntry{}
	}
	return entries, nil
}

func (r *transRepo) ApproveReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID, correctedText string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE review_queue SET status='approved', reviewer_id=$2, corrected_text=$3, reviewed_at=NOW() WHERE id=$1`,
		id, reviewerID, correctedText)
	return err
}

func (r *transRepo) RejectReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE review_queue SET status='rejected', reviewer_id=$2, reviewed_at=NOW() WHERE id=$1`,
		id, reviewerID)
	return err
}

func (r *transRepo) LogUsage(ctx context.Context, provider string, charCount int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO provider_usage (provider, request_count, char_count, recorded_at) VALUES ($1,1,$2,NOW())`,
		provider, charCount)
	return err
}

func (r *transRepo) GetUsage(ctx context.Context) ([]domain.ProviderUsage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT provider, SUM(request_count) as request_count, SUM(char_count) as char_count FROM provider_usage GROUP BY provider`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usages []domain.ProviderUsage
	for rows.Next() {
		u := domain.ProviderUsage{}
		if err := rows.Scan(&u.Provider, &u.RequestCount, &u.CharCount); err != nil {
			return nil, err
		}
		usages = append(usages, u)
	}
	if usages == nil {
		usages = []domain.ProviderUsage{}
	}
	return usages, nil
}

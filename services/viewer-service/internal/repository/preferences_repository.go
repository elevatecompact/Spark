package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
)

type PreferencesRepository interface {
	GetByViewer(ctx context.Context, viewerID uuid.UUID) (*domain.ViewerPreferences, error)
	Upsert(ctx context.Context, prefs *domain.ViewerPreferences) error
	Patch(ctx context.Context, viewerID uuid.UUID, updates domain.UpdatePreferences) (*domain.ViewerPreferences, error)
}

type preferencesRepository struct {
	pool *pgxpool.Pool
}

func NewPreferencesRepository(pool *pgxpool.Pool) PreferencesRepository {
	return &preferencesRepository{pool: pool}
}

func (r *preferencesRepository) GetByViewer(ctx context.Context, viewerID uuid.UUID) (*domain.ViewerPreferences, error) {
	query := `
		SELECT viewer_id, preferred_categories, content_language, autoplay, mature_content_allowed, notification_prefs, created_at, updated_at
		FROM viewer_preferences WHERE viewer_id = $1`

	prefs := &domain.ViewerPreferences{}
	var notifPrefs []byte
	var categories []uuid.UUID
	err := r.pool.QueryRow(ctx, query, viewerID).Scan(
		&prefs.ViewerID, &categories, &prefs.ContentLanguage,
		&prefs.Autoplay, &prefs.MatureContentAllowed,
		&notifPrefs, &prefs.CreatedAt, &prefs.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	prefs.PreferredCategories = categories
	if len(notifPrefs) > 0 {
		json.Unmarshal(notifPrefs, &prefs.NotificationPrefs)
	}
	if prefs.NotificationPrefs == nil {
		prefs.NotificationPrefs = make(map[string]interface{})
	}

	return prefs, nil
}

func (r *preferencesRepository) Upsert(ctx context.Context, prefs *domain.ViewerPreferences) error {
	notifPrefs, err := json.Marshal(prefs.NotificationPrefs)
	if err != nil {
		return fmt.Errorf("failed to marshal notification prefs: %w", err)
	}

	query := `
		INSERT INTO viewer_preferences (viewer_id, preferred_categories, content_language, autoplay, mature_content_allowed, notification_prefs, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (viewer_id) DO UPDATE SET
			preferred_categories = $2, content_language = $3, autoplay = $4,
			mature_content_allowed = $5, notification_prefs = $6, updated_at = $8
		RETURNING created_at, updated_at`

	err = r.pool.QueryRow(ctx, query,
		prefs.ViewerID, prefs.PreferredCategories, prefs.ContentLanguage,
		prefs.Autoplay, prefs.MatureContentAllowed, notifPrefs,
		prefs.CreatedAt, prefs.UpdatedAt,
	).Scan(&prefs.CreatedAt, &prefs.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert preferences: %w", err)
	}
	return nil
}

func (r *preferencesRepository) Patch(ctx context.Context, viewerID uuid.UUID, updates domain.UpdatePreferences) (*domain.ViewerPreferences, error) {
	prefs, err := r.GetByViewer(ctx, viewerID)
	if err != nil {
		prefs = &domain.ViewerPreferences{
			ViewerID:          viewerID,
			PreferredCategories: []uuid.UUID{},
			ContentLanguage:    "en",
			Autoplay:           true,
			MatureContentAllowed: false,
			NotificationPrefs:  make(map[string]interface{}),
			CreatedAt:          time.Now().UTC(),
		}
	}

	changed := false
	if updates.PreferredCategories != nil {
		prefs.PreferredCategories = *updates.PreferredCategories
		changed = true
	}
	if updates.ContentLanguage != nil {
		prefs.ContentLanguage = *updates.ContentLanguage
		changed = true
	}
	if updates.Autoplay != nil {
		prefs.Autoplay = *updates.Autoplay
		changed = true
	}
	if updates.MatureContentAllowed != nil {
		prefs.MatureContentAllowed = *updates.MatureContentAllowed
		changed = true
	}
	if updates.NotificationPrefs != nil {
		prefs.NotificationPrefs = *updates.NotificationPrefs
		changed = true
	}

	if !changed {
		return prefs, nil
	}

	prefs.UpdatedAt = time.Now().UTC()
	if err := r.Upsert(ctx, prefs); err != nil {
		return nil, err
	}

	return prefs, nil
}

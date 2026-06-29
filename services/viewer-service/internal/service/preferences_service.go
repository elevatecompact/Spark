package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/viewer-service/internal/domain"
	"github.com/elevatecompact/spark/services/viewer-service/internal/repository"
)

type PreferencesService interface {
	Get(ctx context.Context, viewerID uuid.UUID) (*domain.ViewerPreferences, error)
	Replace(ctx context.Context, viewerID uuid.UUID, prefs *domain.ViewerPreferences) (*domain.ViewerPreferences, error)
	Patch(ctx context.Context, viewerID uuid.UUID, updates domain.UpdatePreferences) (*domain.ViewerPreferences, error)
}

type preferencesService struct {
	repo repository.PreferencesRepository
	eventPub interface{ PublishPreferencesUpdated(ctx context.Context, viewerID uuid.UUID) error }
}

func NewPreferencesService(repo repository.PreferencesRepository) PreferencesService {
	return &preferencesService{
		repo: repo,
	}
}

func (s *preferencesService) Get(ctx context.Context, viewerID uuid.UUID) (*domain.ViewerPreferences, error) {
	prefs, err := s.repo.GetByViewer(ctx, viewerID)
	if err != nil {
		if err == domain.ErrNotFound {
			return &domain.ViewerPreferences{
				ViewerID:            viewerID,
				PreferredCategories: []uuid.UUID{},
				ContentLanguage:     "en",
				Autoplay:            true,
				MatureContentAllowed: false,
				NotificationPrefs:   make(map[string]interface{}),
			}, nil
		}
		return nil, err
	}
	return prefs, nil
}

func (s *preferencesService) Replace(ctx context.Context, viewerID uuid.UUID, prefs *domain.ViewerPreferences) (*domain.ViewerPreferences, error) {
	prefs.ViewerID = viewerID
	prefs.CreatedAt = time.Now().UTC()
	prefs.UpdatedAt = time.Now().UTC()

	if prefs.PreferredCategories == nil {
		prefs.PreferredCategories = []uuid.UUID{}
	}
	if prefs.ContentLanguage == "" {
		prefs.ContentLanguage = "en"
	}
	if prefs.NotificationPrefs == nil {
		prefs.NotificationPrefs = make(map[string]interface{})
	}

	if err := s.repo.Upsert(ctx, prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

func (s *preferencesService) Patch(ctx context.Context, viewerID uuid.UUID, updates domain.UpdatePreferences) (*domain.ViewerPreferences, error) {
	prefs, err := s.repo.Patch(ctx, viewerID, updates)
	if err != nil {
		return nil, err
	}
	return prefs, nil
}

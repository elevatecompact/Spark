package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/translation-service/internal/domain"
	"github.com/elevatecompact/spark/services/translation-service/internal/events"
	"github.com/elevatecompact/spark/services/translation-service/internal/repository"
)

type TranslationService interface {
	Translate(ctx context.Context, req *domain.TranslationRequest) (*domain.TranslationResult, error)
	TranslateBatch(ctx context.Context, reqs []domain.TranslationRequest) ([]domain.TranslationResult, error)
	Detect(ctx context.Context, text string) (*domain.DetectionResult, error)
	DetectBatch(ctx context.Context, texts []string) ([]domain.DetectionResult, error)
	GetLanguages(ctx context.Context) ([]domain.SupportedLanguage, error)
	LookupMemory(ctx context.Context, text, sourceLang, targetLang string) (*domain.TranslationMemoryEntry, error)
	StoreMemory(ctx context.Context, entry *domain.TranslationMemoryEntry) error
	DeleteMemory(ctx context.Context, id uuid.UUID) error
	ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewEntry, error)
	ApproveReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID, correctedText string) error
	RejectReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error
	GetUsage(ctx context.Context) ([]domain.ProviderUsage, error)
}

type translationService struct {
	repo     repository.TranslationRepository
	eventPub events.EventProducer
}

func NewTranslationService(repo repository.TranslationRepository, eventPub events.EventProducer) TranslationService {
	return &translationService{repo: repo, eventPub: eventPub}
}

func (s *translationService) Translate(ctx context.Context, req *domain.TranslationRequest) (*domain.TranslationResult, error) {
	if len(req.Text) > 5000 {
		return nil, domain.ErrTextTooLong
	}

	sourceHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Text)))
	mem, err := s.repo.LookupMemory(ctx, sourceHash, req.SourceLang, req.TargetLang)
	if err == nil && mem != nil {
		return &domain.TranslationResult{
			TranslatedText: mem.TranslatedText,
			SourceLang:     mem.SourceLang,
			TargetLang:     mem.TargetLang,
			Provider:       mem.Provider,
			Confidence:     mem.QualityScore,
		}, nil
	}

	result := &domain.TranslationResult{
		TranslatedText: req.Text,
		SourceLang:     req.SourceLang,
		TargetLang:     req.TargetLang,
		Provider:       "noop",
		Confidence:     0.95,
	}

	s.repo.SaveMemory(ctx, &domain.TranslationMemoryEntry{
		SourceText:     req.Text,
		TranslatedText: result.TranslatedText,
		SourceLang:     result.SourceLang,
		TargetLang:     result.TargetLang,
		Provider:       result.Provider,
		QualityScore:   result.Confidence,
	})

	s.eventPub.PublishTranslationCompleted(ctx, &events.TranslationCompletedEvent{
		TranslationID:  uuid.New(),
		SourceText:     req.Text,
		TranslatedText: result.TranslatedText,
		SourceLang:     result.SourceLang,
		TargetLang:     result.TargetLang,
		Provider:       result.Provider,
		CharCount:      len(req.Text),
		Timestamp:      time.Now().UTC(),
	})

	return result, nil
}

func (s *translationService) TranslateBatch(ctx context.Context, reqs []domain.TranslationRequest) ([]domain.TranslationResult, error) {
	if len(reqs) > 100 {
		reqs = reqs[:100]
	}
	results := make([]domain.TranslationResult, len(reqs))
	for i, req := range reqs {
		res, err := s.Translate(ctx, &req)
		if err != nil {
			results[i] = domain.TranslationResult{SourceLang: req.SourceLang, TargetLang: req.TargetLang, Confidence: 0}
			continue
		}
		results[i] = *res
	}
	s.eventPub.PublishBatchCompleted(ctx, uuid.New(), len(results))
	return results, nil
}

func (s *translationService) Detect(ctx context.Context, text string) (*domain.DetectionResult, error) {
	if text == "" {
		return nil, domain.ErrValidation
	}
	return &domain.DetectionResult{Language: "en", Confidence: 0.95}, nil
}

func (s *translationService) DetectBatch(ctx context.Context, texts []string) ([]domain.DetectionResult, error) {
	results := make([]domain.DetectionResult, len(texts))
	for i, t := range texts {
		res, err := s.Detect(ctx, t)
		if err != nil {
			results[i] = domain.DetectionResult{Language: "unknown", Confidence: 0}
			continue
		}
		results[i] = *res
	}
	return results, nil
}

func (s *translationService) GetLanguages(ctx context.Context) ([]domain.SupportedLanguage, error) {
	return []domain.SupportedLanguage{
		{Code: "en", Name: "English", Coverage: 1.0},
		{Code: "es", Name: "Spanish", Coverage: 0.85},
		{Code: "fr", Name: "French", Coverage: 0.80},
		{Code: "de", Name: "German", Coverage: 0.78},
		{Code: "ja", Name: "Japanese", Coverage: 0.72},
		{Code: "zh", Name: "Chinese", Coverage: 0.70},
		{Code: "ko", Name: "Korean", Coverage: 0.68},
		{Code: "pt", Name: "Portuguese", Coverage: 0.75},
		{Code: "ar", Name: "Arabic", Coverage: 0.65},
		{Code: "ru", Name: "Russian", Coverage: 0.70},
	}, nil
}

func (s *translationService) LookupMemory(ctx context.Context, text, sourceLang, targetLang string) (*domain.TranslationMemoryEntry, error) {
	sourceHash := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	return s.repo.LookupMemory(ctx, sourceHash, sourceLang, targetLang)
}

func (s *translationService) StoreMemory(ctx context.Context, entry *domain.TranslationMemoryEntry) error {
	return s.repo.SaveMemory(ctx, entry)
}

func (s *translationService) DeleteMemory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMemory(ctx, id)
}

func (s *translationService) ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewEntry, error) {
	return s.repo.ListReviewQueue(ctx, status)
}

func (s *translationService) ApproveReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID, correctedText string) error {
	return s.repo.ApproveReview(ctx, id, reviewerID, correctedText)
}

func (s *translationService) RejectReview(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error {
	return s.repo.RejectReview(ctx, id, reviewerID)
}

func (s *translationService) GetUsage(ctx context.Context) ([]domain.ProviderUsage, error) {
	return s.repo.GetUsage(ctx)
}

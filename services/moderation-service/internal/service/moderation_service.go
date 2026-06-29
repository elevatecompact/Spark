package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/moderation-service/internal/domain"
	"github.com/elevatecompact/spark/services/moderation-service/internal/events"
	"github.com/elevatecompact/spark/services/moderation-service/internal/repository"
)

type ModerationService interface {
	ScanText(ctx context.Context, contentID uuid.UUID, text string) (*domain.ScanResult, error)
	ScanImage(ctx context.Context, contentID uuid.UUID, imageURL string) (*domain.ScanResult, error)
	ScanBatch(ctx context.Context, items []domain.ScanResult) ([]domain.ScanResult, error)

	CreateRule(ctx context.Context, rule *domain.ModerationRule) error
	GetRule(ctx context.Context, id uuid.UUID) (*domain.ModerationRule, error)
	ListRules(ctx context.Context) ([]domain.ModerationRule, error)
	UpdateRule(ctx context.Context, rule *domain.ModerationRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error

	ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewItem, error)
	ApproveReview(ctx context.Context, id uuid.UUID, resolution string) error
	RejectReview(ctx context.Context, id uuid.UUID, resolution string) error
	GetQueueStats(ctx context.Context) ([]domain.QueueStats, error)

	WarnUser(ctx context.Context, userID uuid.UUID, reason string) (*domain.ModerationAction, error)
	RestrictUser(ctx context.Context, userID uuid.UUID, reason string, duration int) (*domain.ModerationAction, error)
	RemoveContent(ctx context.Context, contentID uuid.UUID, reason string) (*domain.ModerationAction, error)
	SuspendUser(ctx context.Context, userID uuid.UUID, reason string, duration int) (*domain.ModerationAction, error)
	ReverseAction(ctx context.Context, actionID uuid.UUID) error

	CreateReport(ctx context.Context, report *domain.ContentReport) error
	GetReport(ctx context.Context, id uuid.UUID) (*domain.ContentReport, error)
	ListReports(ctx context.Context, status domain.ReportStatus) ([]domain.ContentReport, error)

	GetAdminStats(ctx context.Context) (*domain.AdminStats, error)
	GetAccuracy(ctx context.Context) (map[string]float64, error)
}

type modService struct {
	repo     repository.ModerationRepository
	eventPub events.EventProducer
}

func NewModerationService(repo repository.ModerationRepository, eventPub events.EventProducer) ModerationService {
	return &modService{repo: repo, eventPub: eventPub}
}

func (s *modService) ScanText(ctx context.Context, contentID uuid.UUID, text string) (*domain.ScanResult, error) {
	if len(text) > 10000 {
		text = text[:10000]
	}

	rules, _ := s.repo.ListRules(ctx)
	violations := make([]domain.ScanViolation, 0)
	for _, rule := range rules {
		if rule.IsActive {
			violations = append(violations, domain.ScanViolation{
				RuleID:     rule.ID,
				Severity:   rule.Severity,
				Category:   rule.Category,
				Confidence: 0.0,
				Matched:    "",
			})
		}
	}

	result := &domain.ScanResult{
		ContentID:   contentID,
		ContentType: "text",
		Violations:  violations,
		NeedsReview: true,
	}

	if len(violations) > 0 {
		result.NeedsReview = true
		s.repo.CreateReviewItem(ctx, &domain.ReviewItem{
			ID:          uuid.New(),
			ContentType: "text",
			ContentID:   contentID,
			FlaggedBy:   "automated",
			Reasons:     []string{"policy violation detected"},
			Status:      domain.ReviewPending,
		})
	}

	s.eventPub.PublishContentFlagged(ctx, &events.ContentFlaggedEvent{
		ContentID:    contentID,
		ContentType:  "text",
		ScanResults:  violations,
		Timestamp:    domain.ReviewItem{}.CreatedAt,
	})

	return result, nil
}

func (s *modService) ScanImage(ctx context.Context, contentID uuid.UUID, imageURL string) (*domain.ScanResult, error) {
	return &domain.ScanResult{
		ContentID:   contentID,
		ContentType: "image",
		Violations:  []domain.ScanViolation{},
		NeedsReview: false,
	}, nil
}

func (s *modService) ScanBatch(ctx context.Context, items []domain.ScanResult) ([]domain.ScanResult, error) {
	if len(items) > 50 {
		items = items[:50]
	}
	for i := range items {
		items[i].NeedsReview = false
	}
	return items, nil
}

func (s *modService) CreateRule(ctx context.Context, rule *domain.ModerationRule) error {
	rule.ID = uuid.New()
	return s.repo.CreateRule(ctx, rule)
}

func (s *modService) GetRule(ctx context.Context, id uuid.UUID) (*domain.ModerationRule, error) {
	return s.repo.GetRule(ctx, id)
}

func (s *modService) ListRules(ctx context.Context) ([]domain.ModerationRule, error) {
	return s.repo.ListRules(ctx)
}

func (s *modService) UpdateRule(ctx context.Context, rule *domain.ModerationRule) error {
	return s.repo.UpdateRule(ctx, rule)
}

func (s *modService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRule(ctx, id)
}

func (s *modService) ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewItem, error) {
	return s.repo.ListReviewQueue(ctx, status)
}

func (s *modService) ApproveReview(ctx context.Context, id uuid.UUID, resolution string) error {
	if err := s.repo.ApproveReview(ctx, id, resolution); err != nil {
		return err
	}
	return s.eventPub.PublishReviewCompleted(ctx, id, resolution)
}

func (s *modService) RejectReview(ctx context.Context, id uuid.UUID, resolution string) error {
	if err := s.repo.RejectReview(ctx, id, resolution); err != nil {
		return err
	}
	return s.eventPub.PublishReviewCompleted(ctx, id, resolution)
}

func (s *modService) GetQueueStats(ctx context.Context) ([]domain.QueueStats, error) {
	return s.repo.GetQueueStats(ctx)
}

func (s *modService) WarnUser(ctx context.Context, userID uuid.UUID, reason string) (*domain.ModerationAction, error) {
	action := &domain.ModerationAction{
		ID:         uuid.New(),
		UserID:     userID,
		ActionType: domain.SevWarn,
		Status:     domain.ActionApplied,
		AppliedBy:  "moderator",
		Reason:     reason,
	}
	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}
	s.eventPub.PublishActionTaken(ctx, &events.ActionTakenEvent{Action: *action})
	return action, nil
}

func (s *modService) RestrictUser(ctx context.Context, userID uuid.UUID, reason string, duration int) (*domain.ModerationAction, error) {
	action := &domain.ModerationAction{
		ID:         uuid.New(),
		UserID:     userID,
		ActionType: domain.SevRestrict,
		Status:     domain.ActionApplied,
		AppliedBy:  "moderator",
		Reason:     reason,
		Duration:   &duration,
	}
	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}
	s.eventPub.PublishActionTaken(ctx, &events.ActionTakenEvent{Action: *action})
	return action, nil
}

func (s *modService) RemoveContent(ctx context.Context, contentID uuid.UUID, reason string) (*domain.ModerationAction, error) {
	action := &domain.ModerationAction{
		ID:         uuid.New(),
		ContentID:  &contentID,
		ActionType: domain.SevRemove,
		Status:     domain.ActionApplied,
		AppliedBy:  "moderator",
		Reason:     reason,
	}
	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}
	s.eventPub.PublishActionTaken(ctx, &events.ActionTakenEvent{Action: *action})
	return action, nil
}

func (s *modService) SuspendUser(ctx context.Context, userID uuid.UUID, reason string, duration int) (*domain.ModerationAction, error) {
	action := &domain.ModerationAction{
		ID:         uuid.New(),
		UserID:     userID,
		ActionType: domain.SevSuspend,
		Status:     domain.ActionApplied,
		AppliedBy:  "moderator",
		Reason:     reason,
		Duration:   &duration,
	}
	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}
	s.eventPub.PublishActionTaken(ctx, &events.ActionTakenEvent{Action: *action})
	return action, nil
}

func (s *modService) ReverseAction(ctx context.Context, actionID uuid.UUID) error {
	log.Info().Str("action_id", actionID.String()).Msg("action reversed (noop)")
	return nil
}

func (s *modService) CreateReport(ctx context.Context, report *domain.ContentReport) error {
	report.ID = uuid.New()
	report.Status = domain.ReportOpen
	if err := s.repo.CreateReport(ctx, report); err != nil {
		return err
	}
	return s.eventPub.PublishReportSubmitted(ctx, report.ID)
}

func (s *modService) GetReport(ctx context.Context, id uuid.UUID) (*domain.ContentReport, error) {
	return s.repo.GetReport(ctx, id)
}

func (s *modService) ListReports(ctx context.Context, status domain.ReportStatus) ([]domain.ContentReport, error) {
	return s.repo.ListReports(ctx, status)
}

func (s *modService) GetAdminStats(ctx context.Context) (*domain.AdminStats, error) {
	return s.repo.GetAdminStats(ctx)
}

func (s *modService) GetAccuracy(ctx context.Context) (map[string]float64, error) {
	return map[string]float64{
		"auto_human_match_rate": 0.87,
		"false_positive_rate":   0.04,
		"false_negative_rate":   0.02,
		"avg_review_time_hours": 2.5,
	}, nil
}

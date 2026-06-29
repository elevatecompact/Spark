package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/licensing-service/internal/domain"
	"github.com/elevatecompact/spark/services/licensing-service/internal/repository"
)

type LicensingService struct {
	repo *repository.LicensingRepository
	evt  domain.EventProducer
}

func NewLicensingService(repo *repository.LicensingRepository, evt domain.EventProducer) *LicensingService {
	return &LicensingService{repo: repo, evt: evt}
}

func (s *LicensingService) CreateLicense(ctx context.Context, l *domain.License) (*domain.License, error) {
	if l.StartDate.After(l.EndDate) {
		return nil, errors.New("start date must be before end date")
	}
	l.ID = uuid.New()
	l.Status = domain.LicenseDraft
	l.CreatedAt = time.Now()
	if err := s.repo.CreateLicense(ctx, l); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "licensing.license.created", map[string]interface{}{
		"licenseId": l.ID, "rightsHolderId": l.RightsHolderID, "status": l.Status,
	})
	return l, nil
}

func (s *LicensingService) GetLicense(ctx context.Context, id uuid.UUID) (*domain.License, error) {
	return s.repo.GetLicense(ctx, id)
}

func (s *LicensingService) UpdateLicense(ctx context.Context, l *domain.License) error {
	if err := s.repo.UpdateLicense(ctx, l); err != nil {
		return err
	}
	return nil
}

func (s *LicensingService) DeleteLicense(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLicense(ctx, id)
}

func (s *LicensingService) ListLicenses(ctx context.Context, rightsHolderID, licenseeID, contentID *uuid.UUID, limit, offset int) ([]domain.License, error) {
	return s.repo.ListLicenses(ctx, rightsHolderID, licenseeID, contentID, limit, offset)
}

func (s *LicensingService) ApproveLicense(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UpdateLicenseStatus(ctx, id, domain.LicenseActive); err != nil {
		return err
	}
	s.evt.Publish(ctx, "licensing.license.activated", map[string]interface{}{
		"licenseId": id, "status": "active",
	})
	return nil
}

func (s *LicensingService) RejectLicense(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateLicenseStatus(ctx, id, domain.LicenseTerminated)
}

func (s *LicensingService) TerminateLicense(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UpdateLicenseStatus(ctx, id, domain.LicenseTerminated); err != nil {
		return err
	}
	s.evt.Publish(ctx, "licensing.license.terminated", map[string]interface{}{
		"licenseId": id,
	})
	return nil
}

func (s *LicensingService) RegisterContentRight(ctx context.Context, cr *domain.ContentRight) (*domain.ContentRight, error) {
	cr.ID = uuid.New()
	cr.RegisteredAt = time.Now()
	if err := s.repo.RegisterContentRight(ctx, cr); err != nil {
		return nil, err
	}
	return cr, nil
}

func (s *LicensingService) GetContentRights(ctx context.Context, contentID uuid.UUID) (*domain.ContentRight, error) {
	return s.repo.GetContentRights(ctx, contentID)
}

func (s *LicensingService) VerifyContentRights(ctx context.Context, contentID uuid.UUID, usageType string) (*domain.RightsCheckResult, error) {
	license, err := s.repo.GetActiveLicenseByContent(ctx, contentID)
	if err != nil {
		return &domain.RightsCheckResult{ContentID: contentID, Allowed: false, Reason: "no active license"}, nil
	}
	if license.EndDate.Before(time.Now()) {
		return &domain.RightsCheckResult{ContentID: contentID, Allowed: false, Reason: "license expired"}, nil
	}
	licenseID := license.ID
	return &domain.RightsCheckResult{
		ContentID: contentID,
		Allowed:   true,
		LicenseID: &licenseID,
		Reason:    fmt.Sprintf("active %s license", license.Type),
	}, nil
}

func (s *LicensingService) RecordUsage(ctx context.Context, licenseID, contentID uuid.UUID, usageType domain.UsageType, contextData map[string]interface{}) (*domain.UsageLog, error) {
	ctxJSON, _ := json.Marshal(contextData)
	u := &domain.UsageLog{
		ID:         uuid.New(),
		LicenseID:  licenseID,
		ContentID:  contentID,
		UsageType:  usageType,
		Context:    ctxJSON,
		RecordedAt: time.Now(),
	}
	if err := s.repo.RecordUsage(ctx, u); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "licensing.usage.recorded", map[string]interface{}{
		"usageId": u.ID, "licenseId": licenseID, "contentId": contentID, "usageType": usageType,
	})
	return u, nil
}

func (s *LicensingService) GetUsageReport(ctx context.Context, rightsHolderID uuid.UUID, start, end time.Time) (map[string]interface{}, error) {
	count, err := s.repo.GetUsageByRightsHolder(ctx, rightsHolderID, start, end)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"rightsHolderId": rightsHolderID,
		"periodStart":    start,
		"periodEnd":      end,
		"usageCount":     count,
	}, nil
}

func (s *LicensingService) GetUsageByContent(ctx context.Context, contentID uuid.UUID, limit, offset int) ([]domain.UsageLog, error) {
	return s.repo.GetUsageByContent(ctx, contentID, limit, offset)
}

func (s *LicensingService) CalculateProjectedRoyalty(ctx context.Context, licenseID uuid.UUID) (*domain.RoyaltyStatement, error) {
	return s.repo.CalculateProjectedRoyalty(ctx, licenseID)
}

func (s *LicensingService) GenerateRoyaltyStatement(ctx context.Context, licenseID uuid.UUID, periodStart, periodEnd time.Time) (*domain.RoyaltyStatement, error) {
	license, err := s.repo.GetLicense(ctx, licenseID)
	if err != nil {
		return nil, err
	}
	usage, err := s.repo.GetUsageByLicense(ctx, licenseID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	usageCount := int64(len(usage))
	totalCents := license.RateCents * usageCount
	if license.RateType == domain.RateTypeFlat {
		totalCents = license.RateCents
	}
	if license.MinGuaranteeCents > 0 && totalCents < license.MinGuaranteeCents {
		totalCents = license.MinGuaranteeCents
	}
	rs := &domain.RoyaltyStatement{
		ID:             uuid.New(),
		LicenseID:      licenseID,
		RightsHolderID: license.RightsHolderID,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		UsageCount:     usageCount,
		RateApplied:    license.RateCents,
		TotalCents:     totalCents,
		Status:         domain.RoyaltyPending,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.CreateRoyaltyStatement(ctx, rs); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "licensing.royalty.calculated", map[string]interface{}{
		"licenseId": licenseID, "periodStart": periodStart, "periodEnd": periodEnd,
		"usageCount": usageCount, "totalCents": totalCents, "rightsHolderId": license.RightsHolderID,
	})
	return rs, nil
}

func (s *LicensingService) GetRoyaltyStatements(ctx context.Context, rightsHolderID, licenseID *uuid.UUID, limit, offset int) ([]domain.RoyaltyStatement, error) {
	return s.repo.GetRoyaltyStatements(ctx, rightsHolderID, licenseID, limit, offset)
}

func (s *LicensingService) GetPendingRoyalties(ctx context.Context) ([]domain.RoyaltyStatement, error) {
	return s.repo.GetPendingRoyalties(ctx)
}

func (s *LicensingService) ProcessRoyaltyPayout(ctx context.Context, statementID uuid.UUID) error {
	if err := s.repo.UpdateRoyaltyStatus(ctx, statementID, domain.RoyaltyPaid); err != nil {
		return err
	}
	rs, _ := s.repo.GetRoyaltyStatements(ctx, nil, nil, 1, 0)
	if len(rs) > 0 {
		s.evt.Publish(ctx, "licensing.royalty.paid", map[string]interface{}{
			"statementId": statementID, "totalCents": rs[0].TotalCents, "rightsHolderId": rs[0].RightsHolderID,
		})
	}
	return nil
}

func (s *LicensingService) GetComplianceReport(ctx context.Context) (*domain.ComplianceReport, error) {
	return s.repo.GetComplianceReport(ctx)
}

func (s *LicensingService) HandleStreamSessionStarted(ctx context.Context, contentID uuid.UUID, userID uuid.UUID) error {
	result, err := s.VerifyContentRights(ctx, contentID, "stream")
	if err != nil {
		return err
	}
	if !result.Allowed && result.LicenseID != nil {
		s.evt.Publish(ctx, "licensing.compliance.flag", map[string]interface{}{
			"contentId": contentID, "reason": "stream attempted without valid license",
		})
	}
	return nil
}

func (s *LicensingService) HandleMediaContentUploaded(ctx context.Context, contentID uuid.UUID, uploaderID uuid.UUID) error {
	_, err := s.repo.GetContentRights(ctx, contentID)
	if err != nil {
		log.Warn().Str("contentId", contentID.String()).Msg("content uploaded without registered rights")
	}
	return nil
}

func (s *LicensingService) HandleCommerceOrderPlaced(ctx context.Context, contentID uuid.UUID) error {
	result, err := s.VerifyContentRights(ctx, contentID, "download")
	if err != nil {
		return err
	}
	if !result.Allowed {
		return errors.New("licensed content sale attempted without valid rights")
	}
	return nil
}

func (s *LicensingService) HandleContentFlagged(ctx context.Context, contentID uuid.UUID) error {
	s.evt.Publish(ctx, "licensing.compliance.flag", map[string]interface{}{
		"contentId": contentID, "reason": "copyright claim from moderation",
	})
	return nil
}

func (s *LicensingService) HandlePayoutCompleted(ctx context.Context, statementID uuid.UUID) error {
	return s.repo.UpdateRoyaltyStatus(ctx, statementID, domain.RoyaltyPaid)
}

func (s *LicensingService) GetUsageByLicense(ctx context.Context, licenseID uuid.UUID, start, end time.Time) ([]domain.UsageLog, error) {
	return s.repo.GetUsageByLicense(ctx, licenseID, start, end)
}

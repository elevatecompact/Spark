package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/advertising-service/internal/domain"
	"github.com/elevatecompact/spark/services/advertising-service/internal/events"
	"github.com/elevatecompact/spark/services/advertising-service/internal/repository"
)

type AdvertisingService interface {
	CreateCampaign(ctx context.Context, c *domain.Campaign) (*domain.Campaign, error)
	GetCampaign(ctx context.Context, id uuid.UUID) (*domain.Campaign, error)
	UpdateCampaign(ctx context.Context, c *domain.Campaign) error
	ListCampaigns(ctx context.Context, advertiserID uuid.UUID) ([]domain.Campaign, error)
	PauseCampaign(ctx context.Context, id uuid.UUID) error
	ResumeCampaign(ctx context.Context, id uuid.UUID) error

	CreateAdUnit(ctx context.Context, u *domain.AdUnit) (*domain.AdUnit, error)
	GetAdUnit(ctx context.Context, id uuid.UUID) (*domain.AdUnit, error)
	UpdateAdUnit(ctx context.Context, u *domain.AdUnit) error
	DeleteAdUnit(ctx context.Context, id uuid.UUID) error
	ListAdUnits(ctx context.Context, campaignID uuid.UUID) ([]domain.AdUnit, error)
	ApproveAdUnit(ctx context.Context, id uuid.UUID, approved bool, note string) error

	RequestAd(ctx context.Context, placementID string, userID *uuid.UUID) (*domain.AdUnit, error)
	RecordImpression(ctx context.Context, campaignID, adUnitID uuid.UUID, placementID string, userID *uuid.UUID, costMicro int64, device, geo string) error
	RecordClick(ctx context.Context, impressionID uuid.UUID) error

	GetCampaignPerformance(ctx context.Context, id uuid.UUID) (*domain.CampaignPerformance, error)
	GetRevenueStats(ctx context.Context) (*domain.RevenueStats, error)
}

type adService struct {
	repo     repository.AdvertisingRepository
	eventPub events.EventProducer
}

func NewAdvertisingService(repo repository.AdvertisingRepository, eventPub events.EventProducer) AdvertisingService {
	return &adService{repo: repo, eventPub: eventPub}
}

func (s *adService) CreateCampaign(ctx context.Context, c *domain.Campaign) (*domain.Campaign, error) {
	c.ID = uuid.New()
	c.Status = domain.CampDraft
	c.SpentCents = 0
	c.CreatedAt = time.Now().UTC()
	if err := s.repo.CreateCampaign(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *adService) GetCampaign(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	return s.repo.GetCampaign(ctx, id)
}

func (s *adService) UpdateCampaign(ctx context.Context, c *domain.Campaign) error {
	return s.repo.UpdateCampaign(ctx, c)
}

func (s *adService) ListCampaigns(ctx context.Context, advertiserID uuid.UUID) ([]domain.Campaign, error) {
	return s.repo.ListCampaigns(ctx, advertiserID)
}

func (s *adService) PauseCampaign(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return err
	}
	c.Status = domain.CampPaused
	return s.repo.UpdateCampaign(ctx, c)
}

func (s *adService) ResumeCampaign(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return err
	}
	c.Status = domain.CampActive
	return s.repo.UpdateCampaign(ctx, c)
}

func (s *adService) CreateAdUnit(ctx context.Context, u *domain.AdUnit) (*domain.AdUnit, error) {
	u.ID = uuid.New()
	u.Status = domain.AdPending
	if err := s.repo.CreateAdUnit(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *adService) GetAdUnit(ctx context.Context, id uuid.UUID) (*domain.AdUnit, error) {
	return s.repo.GetAdUnit(ctx, id)
}

func (s *adService) UpdateAdUnit(ctx context.Context, u *domain.AdUnit) error {
	return s.repo.UpdateAdUnit(ctx, u)
}

func (s *adService) DeleteAdUnit(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAdUnit(ctx, id)
}

func (s *adService) ListAdUnits(ctx context.Context, campaignID uuid.UUID) ([]domain.AdUnit, error) {
	return s.repo.ListAdUnits(ctx, campaignID)
}

func (s *adService) ApproveAdUnit(ctx context.Context, id uuid.UUID, approved bool, note string) error {
	return s.repo.ApproveAdUnit(ctx, id, approved, note)
}

func (s *adService) RequestAd(ctx context.Context, placementID string, userID *uuid.UUID) (*domain.AdUnit, error) {
	ads, err := s.repo.GetActiveAds(ctx, placementID, 5)
	if err != nil {
		return nil, err
	}
	if len(ads) == 0 {
		return nil, domain.ErrNoAdsAvailable
	}
	return &ads[rand.Intn(len(ads))], nil
}

func (s *adService) RecordImpression(ctx context.Context, campaignID, adUnitID uuid.UUID, placementID string, userID *uuid.UUID, costMicro int64, device, geo string) error {
	imp := &domain.Impression{
		ID:             uuid.New(),
		CampaignID:     campaignID,
		AdUnitID:       adUnitID,
		PlacementID:    placementID,
		UserID:         userID,
		CostMicroCents: costMicro,
		DeviceType:     device,
		Geo:            geo,
	}
	if err := s.repo.RecordImpression(ctx, imp); err != nil {
		return err
	}
	s.eventPub.PublishImpressionRecorded(ctx, &events.ImpressionRecordedEvent{
		ImpressionID: imp.ID,
		CampaignID:   imp.CampaignID,
		AdUnitID:     imp.AdUnitID,
		PlacementID:  imp.PlacementID,
		UserID:       imp.UserID,
		CostCents:    imp.CostMicroCents / 100000,
		ServedAt:     imp.ServedAt,
	})
	return nil
}

func (s *adService) RecordClick(ctx context.Context, impressionID uuid.UUID) error {
	click := &domain.Click{ID: uuid.New(), ImpressionID: impressionID}
	if err := s.repo.RecordClick(ctx, click); err != nil {
		return err
	}
	return s.eventPub.PublishClickRecorded(ctx, impressionID)
}

func (s *adService) GetCampaignPerformance(ctx context.Context, id uuid.UUID) (*domain.CampaignPerformance, error) {
	return s.repo.GetCampaignPerformance(ctx, id)
}

func (s *adService) GetRevenueStats(ctx context.Context) (*domain.RevenueStats, error) {
	return s.repo.GetRevenueStats(ctx)
}

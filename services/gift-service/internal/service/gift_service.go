package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
	"github.com/elevatecompact/spark/services/gift-service/internal/events"
	"github.com/elevatecompact/spark/services/gift-service/internal/repository"
)

type GiftService interface {
	// Catalog
	CreateGiftItem(ctx context.Context, item *domain.GiftItem) error
	GetGiftItem(ctx context.Context, id uuid.UUID) (*domain.GiftItem, error)
	ListGiftItems(ctx context.Context, admin bool) ([]*domain.GiftItem, error)
	UpdateGiftItem(ctx context.Context, item *domain.GiftItem) error
	DeleteGiftItem(ctx context.Context, id uuid.UUID) error

	// Sending
	SendGift(ctx context.Context, senderID uuid.UUID, req domain.SendGiftRequest) (*domain.Gift, error)
	SendBatchGift(ctx context.Context, senderID uuid.UUID, req domain.SendBatchGiftRequest) ([]*domain.Gift, error)
	SendSubscriptionGift(ctx context.Context, senderID uuid.UUID, req domain.SendSubscriptionGiftRequest) (*domain.Gift, error)

	// My Gifts
	GetGift(ctx context.Context, id uuid.UUID) (*domain.Gift, error)
	ListSent(ctx context.Context, senderID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error)
	ListReceived(ctx context.Context, recipientID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error)

	// Gift Cards
	PurchaseGiftCard(ctx context.Context, purchaserID uuid.UUID, req domain.PurchaseGiftCardRequest) (*domain.GiftCard, error)
	RedeemGiftCard(ctx context.Context, userID uuid.UUID, req domain.RedeemGiftCardRequest) (*domain.GiftCard, error)
	GetGiftCardByCode(ctx context.Context, code string) (*domain.GiftCard, error)

	// Campaigns
	CreateCampaign(ctx context.Context, creatorID uuid.UUID, req domain.CreateCampaignRequest) (*domain.GiftCampaign, error)
	ListCampaigns(ctx context.Context, creatorID uuid.UUID) ([]*domain.GiftCampaign, error)
	ApplyCampaignMatch(ctx context.Context, giftID, campaignID uuid.UUID) error

	// Analytics
	GetTopGifts(ctx context.Context, period string, limit int) ([]*domain.Gift, error)
	GetLeaderboard(ctx context.Context, period string, limit int) ([]domain.LeaderboardEntry, error)

	// Admin
	RefundGift(ctx context.Context, id uuid.UUID) error
}

type giftService struct {
	giftRepo     repository.GiftRepository
	itemRepo     repository.GiftItemRepository
	cardRepo     repository.GiftCardRepository
	campaignRepo repository.GiftCampaignRepository
	eventPub     events.EventProducer

	config GiftServiceConfig
}

type GiftServiceConfig struct {
	CardCodeLength     int
	CardExpiryDays     int
	MaxGiftAmountCents int64
	MinGiftAmountCents int64
	CampaignMaxDays    int
	GiftSendingEnabled   bool
	GiftCardsEnabled      bool
	CampaignMatching      bool
	AnonymousGifting      bool
	GiftLeaderboard       bool
	RateLimitGifts      int
	RateLimitCards      int
	RateLimitCampaigns  int
}

func NewGiftService(
	giftRepo repository.GiftRepository,
	itemRepo repository.GiftItemRepository,
	cardRepo repository.GiftCardRepository,
	campaignRepo repository.GiftCampaignRepository,
	eventPub events.EventProducer,
	cfg GiftServiceConfig,
) GiftService {
	return &giftService{
		giftRepo:     giftRepo,
		itemRepo:     itemRepo,
		cardRepo:     cardRepo,
		campaignRepo: campaignRepo,
		eventPub:     eventPub,
		config:       cfg,
	}
}

func (s *giftService) CreateGiftItem(ctx context.Context, item *domain.GiftItem) error {
	now := time.Now().UTC()
	item.ID = uuid.New()
	item.CreatedAt = now
	item.UpdatedAt = now
	return s.itemRepo.Create(ctx, item)
}

func (s *giftService) GetGiftItem(ctx context.Context, id uuid.UUID) (*domain.GiftItem, error) {
	return s.itemRepo.GetByID(ctx, id)
}

func (s *giftService) ListGiftItems(ctx context.Context, admin bool) ([]*domain.GiftItem, error) {
	if admin {
		return s.itemRepo.ListAll(ctx)
	}
	return s.itemRepo.ListActive(ctx)
}

func (s *giftService) UpdateGiftItem(ctx context.Context, item *domain.GiftItem) error {
	return s.itemRepo.Update(ctx, item)
}

func (s *giftService) DeleteGiftItem(ctx context.Context, id uuid.UUID) error {
	return s.itemRepo.SoftDelete(ctx, id)
}

func (s *giftService) SendGift(ctx context.Context, senderID uuid.UUID, req domain.SendGiftRequest) (*domain.Gift, error) {
	if !s.config.GiftSendingEnabled {
		return nil, domain.ErrGiftSendingDisabled
	}

	if err := s.checkGiftRateLimit(ctx, senderID); err != nil {
		return nil, err
	}

	if req.AmountCents < s.config.MinGiftAmountCents {
		return nil, domain.NewDomainErrorMsg(domain.ErrAmountTooSmall, fmt.Sprintf("minimum gift amount is %d cents", s.config.MinGiftAmountCents), 400)
	}
	if req.AmountCents > s.config.MaxGiftAmountCents {
		return nil, domain.NewDomainErrorMsg(domain.ErrAmountTooLarge, fmt.Sprintf("maximum gift amount is %d cents", s.config.MaxGiftAmountCents), 400)
	}

	if req.GiftItemID != nil {
		if _, err := s.itemRepo.GetByID(ctx, *req.GiftItemID); err != nil {
			return nil, err
		}
	}

	var campaign *domain.GiftCampaign
	if req.CampaignID != nil {
		if !s.config.CampaignMatching {
			return nil, domain.ErrCampaignMatchingDisabled
		}
		var err error
		campaign, err = s.campaignRepo.GetByID(ctx, *req.CampaignID)
		if err != nil {
			return nil, err
		}
		now := time.Now().UTC()
		if now.Before(campaign.StartAt) || now.After(campaign.EndAt) {
			return nil, domain.ErrCampaignInactive
		}
	}

	now := time.Now().UTC()
	gift := &domain.Gift{
		ID:          uuid.New(),
		SenderID:    senderID,
		RecipientID: req.RecipientID,
		GiftItemID:  req.GiftItemID,
		AmountCents: req.AmountCents,
		Message:     req.Message,
		CampaignID:  req.CampaignID,
		IsAnonymous: req.IsAnonymous && s.config.AnonymousGifting,
		Status:      domain.GiftCompleted,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.giftRepo.Create(ctx, gift); err != nil {
		return nil, fmt.Errorf("failed to create gift: %w", err)
	}

	if campaign != nil {
		matchAmount := int64(float64(req.AmountCents) * campaign.MatchRatio)
		remaining := campaign.MaxMatchCents - campaign.TotalMatched
		if matchAmount > remaining {
			matchAmount = remaining
		}
		if matchAmount > 0 {
			if err := s.campaignRepo.AddMatchAmount(ctx, campaign.ID, matchAmount); err == nil {
				if err := s.eventPub.PublishCampaignMatch(ctx, gift, campaign, matchAmount); err != nil {
					log.Warn().Err(err).Msg("failed to publish campaign match")
				}
			}
		}
	}

	if err := s.eventPub.PublishSent(ctx, gift); err != nil {
		log.Warn().Err(err).Msg("failed to publish gift.sent")
	}
	if err := s.eventPub.PublishReceived(ctx, gift); err != nil {
		log.Warn().Err(err).Msg("failed to publish gift.received")
	}

	return gift, nil
}

func (s *giftService) SendBatchGift(ctx context.Context, senderID uuid.UUID, req domain.SendBatchGiftRequest) ([]*domain.Gift, error) {
	if len(req.Gifts) > 50 {
		return nil, domain.ErrBatchLimitExceeded
	}

	gifts := make([]*domain.Gift, 0, len(req.Gifts))
	for _, g := range req.Gifts {
		gift, err := s.SendGift(ctx, senderID, g)
		if err != nil {
			return gifts, fmt.Errorf("batch send failed at gift %d: %w", len(gifts)+1, err)
		}
		gifts = append(gifts, gift)
	}
	return gifts, nil
}

func (s *giftService) SendSubscriptionGift(ctx context.Context, senderID uuid.UUID, req domain.SendSubscriptionGiftRequest) (*domain.Gift, error) {
	if !s.config.GiftSendingEnabled {
		return nil, domain.ErrGiftSendingDisabled
	}

	if s.config.MaxGiftAmountCents > 0 && !s.config.AnonymousGifting {
		// rate limit check
	}

	now := time.Now().UTC()
	gift := &domain.Gift{
		ID:          uuid.New(),
		SenderID:    senderID,
		RecipientID: req.RecipientID,
		AmountCents: 0,
		Message:     req.Message,
		IsAnonymous: req.IsAnonymous && s.config.AnonymousGifting,
		Status:      domain.GiftPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.giftRepo.Create(ctx, gift); err != nil {
		return nil, fmt.Errorf("failed to create subscription gift: %w", err)
	}

	if err := s.eventPub.PublishSubscriptionGifted(ctx, senderID, req.RecipientID, req.PlanID); err != nil {
		log.Warn().Err(err).Msg("failed to publish subscription.gifted")
	}

	return gift, nil
}

func (s *giftService) GetGift(ctx context.Context, id uuid.UUID) (*domain.Gift, error) {
	return s.giftRepo.GetByID(ctx, id)
}

func (s *giftService) ListSent(ctx context.Context, senderID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.giftRepo.ListBySender(ctx, senderID, cursor, limit)
}

func (s *giftService) ListReceived(ctx context.Context, recipientID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error) {
	if cursor.IsZero() {
		cursor = time.Now().UTC().Add(time.Hour)
	}
	return s.giftRepo.ListByRecipient(ctx, recipientID, cursor, limit)
}

func (s *giftService) PurchaseGiftCard(ctx context.Context, purchaserID uuid.UUID, req domain.PurchaseGiftCardRequest) (*domain.GiftCard, error) {
	if !s.config.GiftCardsEnabled {
		return nil, domain.ErrGiftCardsDisabled
	}

	if req.AmountCents < s.config.MinGiftAmountCents {
		return nil, domain.NewDomainErrorMsg(domain.ErrAmountTooSmall, fmt.Sprintf("minimum gift card amount is %d cents", s.config.MinGiftAmountCents), 400)
	}
	if req.AmountCents > s.config.MaxGiftAmountCents {
		return nil, domain.NewDomainErrorMsg(domain.ErrAmountTooLarge, fmt.Sprintf("maximum gift card amount is %d cents", s.config.MaxGiftAmountCents), 400)
	}

	since := time.Now().UTC().AddDate(0, 0, -1)
	count, err := s.cardRepo.CountByPurchaserSince(ctx, purchaserID, since)
	if err != nil {
		return nil, err
	}
	if count >= s.config.RateLimitCards {
		return nil, domain.ErrRateLimitExceeded
	}

	code, err := generateCode(s.config.CardCodeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate gift card code: %w", err)
	}

	now := time.Now().UTC()
	card := &domain.GiftCard{
		ID:           uuid.New(),
		Code:         code,
		PurchaserID:  purchaserID,
		BalanceCents: req.AmountCents,
		ExpiresAt:    now.AddDate(0, 0, s.config.CardExpiryDays),
		CreatedAt:    now,
	}

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("failed to create gift card: %w", err)
	}

	if err := s.eventPub.PublishCardPurchased(ctx, card); err != nil {
		log.Warn().Err(err).Msg("failed to publish card.purchased")
	}

	return card, nil
}

func (s *giftService) RedeemGiftCard(ctx context.Context, userID uuid.UUID, req domain.RedeemGiftCardRequest) (*domain.GiftCard, error) {
	if !s.config.GiftCardsEnabled {
		return nil, domain.ErrGiftCardsDisabled
	}

	card, err := s.cardRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	if card.RedeemedAt != nil {
		return nil, domain.ErrGiftCardRedeemed
	}

	if time.Now().UTC().After(card.ExpiresAt) {
		return nil, domain.ErrGiftCardExpired
	}

	if err := s.cardRepo.MarkRedeemed(ctx, card.ID); err != nil {
		return nil, err
	}

	card.RedeemedAt = timePtr(time.Now().UTC())

	if err := s.eventPub.PublishCardRedeemed(ctx, card); err != nil {
		log.Warn().Err(err).Msg("failed to publish card.redeemed")
	}

	return card, nil
}

func (s *giftService) GetGiftCardByCode(ctx context.Context, code string) (*domain.GiftCard, error) {
	return s.cardRepo.GetByCode(ctx, code)
}

func (s *giftService) CreateCampaign(ctx context.Context, creatorID uuid.UUID, req domain.CreateCampaignRequest) (*domain.GiftCampaign, error) {
	if !s.config.CampaignMatching {
		return nil, domain.ErrCampaignMatchingDisabled
	}

	duration := req.EndAt.Sub(req.StartAt)
	if duration.Hours() > float64(s.config.CampaignMaxDays*24) {
		return nil, domain.NewDomainErrorMsg(domain.ErrValidation, fmt.Sprintf("campaign duration cannot exceed %d days", s.config.CampaignMaxDays), 400)
	}

	since := time.Now().UTC().AddDate(0, -1, 0)
	count, err := s.campaignRepo.CountByCreatorSince(ctx, creatorID, since)
	if err != nil {
		return nil, err
	}
	if count >= s.config.RateLimitCampaigns {
		return nil, domain.ErrRateLimitExceeded
	}

	now := time.Now().UTC()
	campaign := &domain.GiftCampaign{
		ID:            uuid.New(),
		CreatorID:     creatorID,
		MatchRatio:    req.MatchRatio,
		MaxMatchCents: req.MaxMatchCents,
		TotalMatched:  0,
		StartAt:       req.StartAt,
		EndAt:         req.EndAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.campaignRepo.Create(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

func (s *giftService) ListCampaigns(ctx context.Context, creatorID uuid.UUID) ([]*domain.GiftCampaign, error) {
	return s.campaignRepo.ListByCreator(ctx, creatorID)
}

func (s *giftService) ApplyCampaignMatch(ctx context.Context, giftID, campaignID uuid.UUID) error {
	if !s.config.CampaignMatching {
		return domain.ErrCampaignMatchingDisabled
	}

	gift, err := s.giftRepo.GetByID(ctx, giftID)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.Before(campaign.StartAt) || now.After(campaign.EndAt) {
		return domain.ErrCampaignInactive
	}

	matchAmount := int64(float64(gift.AmountCents) * campaign.MatchRatio)
	remaining := campaign.MaxMatchCents - campaign.TotalMatched
	if matchAmount > remaining {
		matchAmount = remaining
	}
	if matchAmount <= 0 {
		return domain.ErrCampaignBudgetExhausted
	}

	if err := s.campaignRepo.AddMatchAmount(ctx, campaignID, matchAmount); err != nil {
		return err
	}

	return s.eventPub.PublishCampaignMatch(ctx, gift, campaign, matchAmount)
}

func (s *giftService) RefundGift(ctx context.Context, id uuid.UUID) error {
	gift, err := s.giftRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if gift.Status != domain.GiftCompleted {
		return domain.ErrGiftNotCompleted
	}
	return s.giftRepo.UpdateStatus(ctx, id, domain.GiftRefunded)
}

func (s *giftService) GetTopGifts(ctx context.Context, period string, limit int) ([]*domain.Gift, error) {
	periodStart := periodToStart(period)
	return s.giftRepo.GetTopGifts(ctx, periodStart, limit)
}

func (s *giftService) GetLeaderboard(ctx context.Context, period string, limit int) ([]domain.LeaderboardEntry, error) {
	if !s.config.GiftLeaderboard {
		return []domain.LeaderboardEntry{}, nil
	}
	periodStart := periodToStart(period)
	return s.giftRepo.GetLeaderboard(ctx, periodStart, limit)
}

func (s *giftService) checkGiftRateLimit(ctx context.Context, senderID uuid.UUID) error {
	since := time.Now().UTC().Add(-1 * time.Hour)
	count, err := s.giftRepo.CountBySenderSince(ctx, senderID, since)
	if err != nil {
		return err
	}
	if count >= s.config.RateLimitGifts {
		return domain.ErrRateLimitExceeded
	}
	return nil
}

func periodToStart(period string) time.Time {
	now := time.Now().UTC()
	switch period {
	case "weekly":
		return now.AddDate(0, 0, -7)
	case "monthly":
		return now.AddDate(0, -1, 0)
	case "all":
		return time.Time{}
	default:
		return now.AddDate(0, -1, 0)
	}
}

func generateCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

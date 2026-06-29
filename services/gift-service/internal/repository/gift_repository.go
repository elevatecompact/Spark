package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/gift-service/internal/domain"
)

type GiftRepository interface {
	Create(ctx context.Context, gift *domain.Gift) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Gift, error)
	ListBySender(ctx context.Context, senderID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error)
	ListByRecipient(ctx context.Context, recipientID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.GiftStatus) error
	CountBySenderSince(ctx context.Context, senderID uuid.UUID, since time.Time) (int, error)
	GetTopGifts(ctx context.Context, periodStart time.Time, limit int) ([]*domain.Gift, error)
	GetLeaderboard(ctx context.Context, periodStart time.Time, limit int) ([]domain.LeaderboardEntry, error)
}

type giftRepository struct {
	pool *pgxpool.Pool
}

func NewGiftRepository(pool *pgxpool.Pool) GiftRepository {
	return &giftRepository{pool: pool}
}

func (r *giftRepository) Create(ctx context.Context, gift *domain.Gift) error {
	query := `INSERT INTO gifts (id, sender_id, recipient_id, gift_item_id, amount_cents, message, campaign_id, is_anonymous, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, gift.ID, gift.SenderID, gift.RecipientID, gift.GiftItemID, gift.AmountCents, gift.Message, gift.CampaignID, gift.IsAnonymous, gift.Status, gift.CreatedAt, gift.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create gift: %w", err)
	}
	return nil
}

func (r *giftRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Gift, error) {
	query := `SELECT id, sender_id, recipient_id, gift_item_id, amount_cents, message, campaign_id, is_anonymous, status, created_at, updated_at
		FROM gifts WHERE id = $1`
	g := &domain.Gift{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&g.ID, &g.SenderID, &g.RecipientID, &g.GiftItemID, &g.AmountCents, &g.Message, &g.CampaignID, &g.IsAnonymous, &g.Status, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGiftNotFound
		}
		return nil, fmt.Errorf("failed to get gift: %w", err)
	}
	return g, nil
}

func (r *giftRepository) ListBySender(ctx context.Context, senderID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id, sender_id, recipient_id, gift_item_id, amount_cents, message, campaign_id, is_anonymous, status, created_at, updated_at
		FROM gifts WHERE sender_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`, senderID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list gifts by sender: %w", err)
	}
	defer rows.Close()
	return scanGifts(rows)
}

func (r *giftRepository) ListByRecipient(ctx context.Context, recipientID uuid.UUID, cursor time.Time, limit int) ([]*domain.Gift, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id, sender_id, recipient_id, gift_item_id, amount_cents, message, campaign_id, is_anonymous, status, created_at, updated_at
		FROM gifts WHERE recipient_id = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`, recipientID, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list gifts by recipient: %w", err)
	}
	defer rows.Close()
	return scanGifts(rows)
}

func (r *giftRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.GiftStatus) error {
	tag, err := r.pool.Exec(ctx, `UPDATE gifts SET status=$2, updated_at=NOW() WHERE id=$1`, id, status)
	if err != nil {
		return fmt.Errorf("failed to update gift status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrGiftNotFound
	}
	return nil
}

func (r *giftRepository) CountBySenderSince(ctx context.Context, senderID uuid.UUID, since time.Time) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM gifts WHERE sender_id = $1 AND created_at >= $2`, senderID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count gifts: %w", err)
	}
	return count, nil
}

func (r *giftRepository) GetTopGifts(ctx context.Context, periodStart time.Time, limit int) ([]*domain.Gift, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT id, sender_id, recipient_id, gift_item_id, amount_cents, message, campaign_id, is_anonymous, status, created_at, updated_at
		FROM gifts WHERE created_at >= $1 AND status = 'completed' ORDER BY amount_cents DESC LIMIT $2`, periodStart, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top gifts: %w", err)
	}
	defer rows.Close()
	return scanGifts(rows)
}

func (r *giftRepository) GetLeaderboard(ctx context.Context, periodStart time.Time, limit int) ([]domain.LeaderboardEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `SELECT recipient_id AS user_id, COUNT(*) AS gift_count, SUM(amount_cents) AS total_cents
		FROM gifts WHERE created_at >= $1 AND status = 'completed' AND is_anonymous = false
		GROUP BY recipient_id ORDER BY total_cents DESC LIMIT $2`, periodStart, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []domain.LeaderboardEntry
	rank := 1
	for rows.Next() {
		e := domain.LeaderboardEntry{Rank: rank}
		if err := rows.Scan(&e.UserID, &e.GiftCount, &e.TotalCents); err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}
		entries = append(entries, e)
		rank++
	}
	if entries == nil {
		entries = []domain.LeaderboardEntry{}
	}
	return entries, nil
}

func scanGifts(rows pgx.Rows) ([]*domain.Gift, error) {
	var gifts []*domain.Gift
	for rows.Next() {
		g := &domain.Gift{}
		if err := rows.Scan(&g.ID, &g.SenderID, &g.RecipientID, &g.GiftItemID, &g.AmountCents, &g.Message, &g.CampaignID, &g.IsAnonymous, &g.Status, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan gift: %w", err)
		}
		gifts = append(gifts, g)
	}
	if gifts == nil {
		gifts = []*domain.Gift{}
	}
	return gifts, nil
}

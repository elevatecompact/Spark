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

type GiftCardRepository interface {
	Create(ctx context.Context, card *domain.GiftCard) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftCard, error)
	GetByCode(ctx context.Context, code string) (*domain.GiftCard, error)
	ListByPurchaser(ctx context.Context, purchaserID uuid.UUID) ([]*domain.GiftCard, error)
	MarkRedeemed(ctx context.Context, id uuid.UUID) error
	CountByPurchaserSince(ctx context.Context, purchaserID uuid.UUID, since time.Time) (int, error)
}

type giftCardRepository struct {
	pool *pgxpool.Pool
}

func NewGiftCardRepository(pool *pgxpool.Pool) GiftCardRepository {
	return &giftCardRepository{pool: pool}
}

func (r *giftCardRepository) Create(ctx context.Context, card *domain.GiftCard) error {
	query := `INSERT INTO gift_cards (id, code, purchaser_id, balance_cents, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, card.ID, card.Code, card.PurchaserID, card.BalanceCents, card.ExpiresAt, card.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create gift card: %w", err)
	}
	return nil
}

func (r *giftCardRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GiftCard, error) {
	query := `SELECT id, code, purchaser_id, balance_cents, expires_at, redeemed_at, created_at
		FROM gift_cards WHERE id = $1`
	card := &domain.GiftCard{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&card.ID, &card.Code, &card.PurchaserID, &card.BalanceCents, &card.ExpiresAt, &card.RedeemedAt, &card.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGiftCardNotFound
		}
		return nil, fmt.Errorf("failed to get gift card: %w", err)
	}
	return card, nil
}

func (r *giftCardRepository) GetByCode(ctx context.Context, code string) (*domain.GiftCard, error) {
	query := `SELECT id, code, purchaser_id, balance_cents, expires_at, redeemed_at, created_at
		FROM gift_cards WHERE code = $1`
	card := &domain.GiftCard{}
	err := r.pool.QueryRow(ctx, query, code).Scan(&card.ID, &card.Code, &card.PurchaserID, &card.BalanceCents, &card.ExpiresAt, &card.RedeemedAt, &card.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrGiftCardNotFound
		}
		return nil, fmt.Errorf("failed to get gift card by code: %w", err)
	}
	return card, nil
}

func (r *giftCardRepository) ListByPurchaser(ctx context.Context, purchaserID uuid.UUID) ([]*domain.GiftCard, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, code, purchaser_id, balance_cents, expires_at, redeemed_at, created_at
		FROM gift_cards WHERE purchaser_id = $1 ORDER BY created_at DESC`, purchaserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list gift cards: %w", err)
	}
	defer rows.Close()

	var cards []*domain.GiftCard
	for rows.Next() {
		card := &domain.GiftCard{}
		if err := rows.Scan(&card.ID, &card.Code, &card.PurchaserID, &card.BalanceCents, &card.ExpiresAt, &card.RedeemedAt, &card.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan gift card: %w", err)
		}
		cards = append(cards, card)
	}
	if cards == nil {
		cards = []*domain.GiftCard{}
	}
	return cards, nil
}

func (r *giftCardRepository) MarkRedeemed(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE gift_cards SET redeemed_at=$2 WHERE id=$1 AND redeemed_at IS NULL`, id, now)
	if err != nil {
		return fmt.Errorf("failed to mark gift card redeemed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrGiftCardRedeemed
	}
	return nil
}

func (r *giftCardRepository) CountByPurchaserSince(ctx context.Context, purchaserID uuid.UUID, since time.Time) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM gift_cards WHERE purchaser_id = $1 AND created_at >= $2`, purchaserID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count gift cards: %w", err)
	}
	return count, nil
}

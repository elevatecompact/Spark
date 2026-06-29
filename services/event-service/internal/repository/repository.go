package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/event-service/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, e *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, e *domain.Event) error
	List(ctx context.Context, category string, status domain.EventStatus, offset, limit int) ([]domain.Event, error)

	CreateTicketTier(ctx context.Context, t *domain.EventTicketTier) error
	ListTicketTiers(ctx context.Context, eventID uuid.UUID) ([]domain.EventTicketTier, error)

	RegisterAttendee(ctx context.Context, a *domain.EventAttendee) error
	ListAttendees(ctx context.Context, eventID uuid.UUID) ([]domain.EventAttendee, error)

	CreateSession(ctx context.Context, s *domain.EventSession) error
	UpdateSession(ctx context.Context, s *domain.EventSession) error
	ListSessions(ctx context.Context, eventID uuid.UUID) ([]domain.EventSession, error)

	CreateSeries(ctx context.Context, s *domain.EventSeries) error
	GetSeries(ctx context.Context, id uuid.UUID) (*domain.EventSeries, error)
	UpdateSeries(ctx context.Context, s *domain.EventSeries) error
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	GetAdminStats(ctx context.Context) (*domain.EventAdminStats, error)
}

type eventRepo struct{ pool *pgxpool.Pool }

func NewEventRepository(pool *pgxpool.Pool) EventRepository { return &eventRepo{pool} }

func (r *eventRepo) Create(ctx context.Context, e *domain.Event) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO events (id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,NOW())`,
		e.ID, e.CreatorID, e.Title, e.Description, e.Category, e.Type, e.StartAt, e.EndAt, e.Timezone, e.MaxAttendees, e.StreamID, e.Status, e.CoverImageURL)
	return err
}

func (r *eventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	e := &domain.Event{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at FROM events WHERE id=$1`, id).
		Scan(&e.ID, &e.CreatorID, &e.Title, &e.Description, &e.Category, &e.Type, &e.StartAt, &e.EndAt, &e.Timezone, &e.MaxAttendees, &e.StreamID, &e.Status, &e.CoverImageURL, &e.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return e, err
}

func (r *eventRepo) Update(ctx context.Context, e *domain.Event) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE events SET title=$2, description=$3, category=$4, type=$5, start_at=$6, end_at=$7, timezone=$8, max_attendees=$9, stream_id=$10, status=$11, cover_image_url=$12 WHERE id=$1`,
		e.ID, e.Title, e.Description, e.Category, e.Type, e.StartAt, e.EndAt, e.Timezone, e.MaxAttendees, e.StreamID, e.Status, e.CoverImageURL)
	return err
}

func (r *eventRepo) List(ctx context.Context, category string, status domain.EventStatus, offset, limit int) ([]domain.Event, error) {
	if limit <= 0 || limit > 50 {
		limit = 25
	}
	var rows pgx.Rows
	var err error
	if category != "" && status != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at FROM events WHERE category=$1 AND status=$2 ORDER BY start_at DESC OFFSET $3 LIMIT $4`,
			category, status, offset, limit)
	} else if category != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at FROM events WHERE category=$1 ORDER BY start_at DESC OFFSET $2 LIMIT $3`,
			category, offset, limit)
	} else if status != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at FROM events WHERE status=$1 ORDER BY start_at DESC OFFSET $2 LIMIT $3`,
			status, offset, limit)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, category, type, start_at, end_at, timezone, max_attendees, stream_id, status, cover_image_url, created_at FROM events ORDER BY start_at DESC OFFSET $1 LIMIT $2`,
			offset, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.CreatorID, &e.Title, &e.Description, &e.Category, &e.Type, &e.StartAt, &e.EndAt, &e.Timezone, &e.MaxAttendees, &e.StreamID, &e.Status, &e.CoverImageURL, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if events == nil {
		events = []domain.Event{}
	}
	return events, nil
}

func (r *eventRepo) CreateTicketTier(ctx context.Context, t *domain.EventTicketTier) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_ticket_tiers (id, event_id, name, price_cents, quantity_total, quantity_sold, benefits, sales_start_at, sales_end_at) VALUES ($1,$2,$3,$4,$5,0,$6,$7,$8)`,
		t.ID, t.EventID, t.Name, t.PriceCents, t.QuantityTotal, t.Benefits, t.SalesStartAt, t.SalesEndAt)
	return err
}

func (r *eventRepo) ListTicketTiers(ctx context.Context, eventID uuid.UUID) ([]domain.EventTicketTier, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, name, price_cents, quantity_total, quantity_sold, benefits, sales_start_at, sales_end_at FROM event_ticket_tiers WHERE event_id=$1 ORDER BY price_cents ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tiers []domain.EventTicketTier
	for rows.Next() {
		var t domain.EventTicketTier
		if err := rows.Scan(&t.ID, &t.EventID, &t.Name, &t.PriceCents, &t.QuantityTotal, &t.QuantitySold, &t.Benefits, &t.SalesStartAt, &t.SalesEndAt); err != nil {
			return nil, err
		}
		tiers = append(tiers, t)
	}
	if tiers == nil {
		tiers = []domain.EventTicketTier{}
	}
	return tiers, nil
}

func (r *eventRepo) RegisterAttendee(ctx context.Context, a *domain.EventAttendee) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_attendees (event_id, ticket_tier_id, user_id, status, registered_at) VALUES ($1,$2,$3,$4,NOW()) ON CONFLICT DO NOTHING`,
		a.EventID, a.TicketTierID, a.UserID, a.Status)
	return err
}

func (r *eventRepo) ListAttendees(ctx context.Context, eventID uuid.UUID) ([]domain.EventAttendee, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT event_id, ticket_tier_id, user_id, status, registered_at, attended_at FROM event_attendees WHERE event_id=$1 ORDER BY registered_at ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var attendees []domain.EventAttendee
	for rows.Next() {
		var a domain.EventAttendee
		if err := rows.Scan(&a.EventID, &a.TicketTierID, &a.UserID, &a.Status, &a.RegisteredAt, &a.AttendedAt); err != nil {
			return nil, err
		}
		attendees = append(attendees, a)
	}
	if attendees == nil {
		attendees = []domain.EventAttendee{}
	}
	return attendees, nil
}

func (r *eventRepo) CreateSession(ctx context.Context, s *domain.EventSession) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_sessions (id, event_id, title, speaker, start_at, end_at, stream_id) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.EventID, s.Title, s.Speaker, s.StartAt, s.EndAt, s.StreamID)
	return err
}

func (r *eventRepo) UpdateSession(ctx context.Context, s *domain.EventSession) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE event_sessions SET title=$2, speaker=$3, start_at=$4, end_at=$5, stream_id=$6 WHERE id=$1`,
		s.ID, s.Title, s.Speaker, s.StartAt, s.EndAt, s.StreamID)
	return err
}

func (r *eventRepo) ListSessions(ctx context.Context, eventID uuid.UUID) ([]domain.EventSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, title, speaker, start_at, end_at, stream_id FROM event_sessions WHERE event_id=$1 ORDER BY start_at ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []domain.EventSession
	for rows.Next() {
		var s domain.EventSession
		if err := rows.Scan(&s.ID, &s.EventID, &s.Title, &s.Speaker, &s.StartAt, &s.EndAt, &s.StreamID); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if sessions == nil {
		sessions = []domain.EventSession{}
	}
	return sessions, nil
}

func (r *eventRepo) CreateSeries(ctx context.Context, s *domain.EventSeries) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_series (id, creator_id, title, description, frequency, day_of_week, start_time, timezone, next_event_at, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,true)`,
		s.ID, s.CreatorID, s.Title, s.Description, s.Frequency, s.DayOfWeek, s.StartTime, s.Timezone, s.NextEventAt)
	return err
}

func (r *eventRepo) GetSeries(ctx context.Context, id uuid.UUID) (*domain.EventSeries, error) {
	s := &domain.EventSeries{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, creator_id, title, description, frequency, day_of_week, start_time, timezone, next_event_at, is_active FROM event_series WHERE id=$1`, id).
		Scan(&s.ID, &s.CreatorID, &s.Title, &s.Description, &s.Frequency, &s.DayOfWeek, &s.StartTime, &s.Timezone, &s.NextEventAt, &s.IsActive)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return s, err
}

func (r *eventRepo) UpdateSeries(ctx context.Context, s *domain.EventSeries) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE event_series SET title=$2, description=$3, frequency=$4, day_of_week=$5, start_time=$6, timezone=$7, next_event_at=$8, is_active=$9 WHERE id=$1`,
		s.ID, s.Title, s.Description, s.Frequency, s.DayOfWeek, s.StartTime, s.Timezone, s.NextEventAt, s.IsActive)
	return err
}

func (r *eventRepo) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM event_series WHERE id=$1`, id)
	return err
}

func (r *eventRepo) GetAdminStats(ctx context.Context) (*domain.EventAdminStats, error) {
	s := &domain.EventAdminStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
		 (SELECT COUNT(*) FROM events) as total_events,
		 (SELECT COUNT(*) FROM events WHERE status='published') as published_events,
		 (SELECT COUNT(*) FROM event_attendees WHERE status='registered') as total_attendees,
		 (SELECT COALESCE(SUM(ett.price_cents), 0) FROM event_attendees ea JOIN event_ticket_tiers ett ON ea.ticket_tier_id = ett.id WHERE ea.status='registered') as revenue_cents`).
		Scan(&s.TotalEvents, &s.PublishedEvents, &s.TotalAttendees, &s.RevenueCents)
	if err != nil {
		return nil, err
	}
	return s, nil
}

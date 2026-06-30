package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/notification-service/internal/domain"
)

type NotificationRepository interface {
	Insert(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Notification, error)
	MarkRead(ctx context.Context, id uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
	DeliveryStats(ctx context.Context) (*DeliveryStatsResult, error)
}

type DeliveryStatsResult struct {
	PushDelivered  int64 `json:"push_delivered"`
	EmailDelivered int64 `json:"email_delivered"`
	SMSDelivered   int64 `json:"sms_delivered"`
	InAppDelivered int64 `json:"inapp_delivered"`
}

type PreferenceRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreference, error)
	Upsert(ctx context.Context, pref *domain.NotificationPreference) error
}

type DeviceRepository interface {
	Register(ctx context.Context, device *domain.PushDevice) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PushDevice, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PushDevice, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type TemplateRepository interface {
	Create(ctx context.Context, t *domain.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
	GetByType(ctx context.Context, typ string) (*domain.Template, error)
	List(ctx context.Context) ([]*domain.Template, error)
	Update(ctx context.Context, t *domain.Template) error
}

type notifRepo struct{ pool *pgxpool.Pool }
type prefRepo struct{ pool *pgxpool.Pool }
type devRepo struct{ pool *pgxpool.Pool }
type tmplRepo struct{ pool *pgxpool.Pool }

func NewNotificationRepository(pool *pgxpool.Pool) NotificationRepository { return &notifRepo{pool} }
func NewPreferenceRepository(pool *pgxpool.Pool) PreferenceRepository      { return &prefRepo{pool} }
func NewDeviceRepository(pool *pgxpool.Pool) DeviceRepository              { return &devRepo{pool} }
func NewTemplateRepository(pool *pgxpool.Pool) TemplateRepository          { return &tmplRepo{pool} }

func (r *notifRepo) Insert(ctx context.Context, n *domain.Notification) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications (id, user_id, type, title, body, data, channel, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		n.ID, n.UserID, n.Type, n.Title, n.Body, n.Data, n.Channel, n.CreatedAt)
	return err
}

func (r *notifRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	n := &domain.Notification{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, type, title, body, data, channel, read_at, created_at FROM notifications WHERE id=$1`, id).
		Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Data, &n.Channel, &n.ReadAt, &n.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotifNotFound
	}
	return n, err
}

func (r *notifRepo) ListByUser(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]*domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, type, title, body, data, channel, read_at, created_at FROM notifications WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3`, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ns []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Data, &n.Channel, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	if ns == nil {
		ns = []*domain.Notification{}
	}
	return ns, nil
}

func (r *notifRepo) MarkRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at=$2 WHERE id=$1 AND read_at IS NULL`, id, now)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotifNotFound
	}
	return nil
}

func (r *notifRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at=NOW() WHERE user_id=$1 AND read_at IS NULL`, userID)
	return err
}

func (r *notifRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM notifications WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotifNotFound
	}
	return nil
}

func (r *notifRepo) DeliveryStats(ctx context.Context) (*DeliveryStatsResult, error) {
	result := &DeliveryStatsResult{}
	row := r.pool.QueryRow(ctx, `SELECT
		COALESCE(SUM(CASE WHEN channel='push' THEN 1 ELSE 0 END), 0) AS push_delivered,
		COALESCE(SUM(CASE WHEN channel='email' THEN 1 ELSE 0 END), 0) AS email_delivered,
		COALESCE(SUM(CASE WHEN channel='sms' THEN 1 ELSE 0 END), 0) AS sms_delivered,
		COALESCE(SUM(CASE WHEN channel='inapp' THEN 1 ELSE 0 END), 0) AS inapp_delivered
	FROM notifications`)
	err := row.Scan(&result.PushDelivered, &result.EmailDelivered, &result.SMSDelivered, &result.InAppDelivered)
	if err == pgx.ErrNoRows {
		return result, nil
	}
	return result, err
}

func (r *notifRepo) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notifications WHERE user_id=$1`, userID)
	return err
}

func (r *prefRepo) Get(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreference, error) {
	p := &domain.NotificationPreference{}
	err := r.pool.QueryRow(ctx, `SELECT user_id, preferences FROM notification_preferences WHERE user_id=$1`, userID).Scan(&p.UserID, &p.Preferences)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *prefRepo) Upsert(ctx context.Context, p *domain.NotificationPreference) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notification_preferences (user_id, preferences) VALUES ($1,$2) ON CONFLICT (user_id) DO UPDATE SET preferences=$2`, p.UserID, p.Preferences)
	return err
}

func (r *devRepo) Register(ctx context.Context, d *domain.PushDevice) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO push_devices (id, user_id, platform, token, is_active, created_at) VALUES ($1,$2,$3,$4,$5,$6)`, d.ID, d.UserID, d.Platform, d.Token, d.IsActive, d.CreatedAt)
	return err
}

func (r *devRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PushDevice, error) {
	d := &domain.PushDevice{}
	err := r.pool.QueryRow(ctx, `SELECT id, user_id, platform, token, is_active, created_at FROM push_devices WHERE id=$1`, id).
		Scan(&d.ID, &d.UserID, &d.Platform, &d.Token, &d.IsActive, &d.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrDeviceNotFound
	}
	return d, err
}

func (r *devRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.PushDevice, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, platform, token, is_active, created_at FROM push_devices WHERE user_id=$1 AND is_active=true`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ds []*domain.PushDevice
	for rows.Next() {
		d := &domain.PushDevice{}
		if err := rows.Scan(&d.ID, &d.UserID, &d.Platform, &d.Token, &d.IsActive, &d.CreatedAt); err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}
	if ds == nil {
		ds = []*domain.PushDevice{}
	}
	return ds, nil
}

func (r *devRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `UPDATE push_devices SET is_active=false WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrDeviceNotFound
	}
	return nil
}

func (r *tmplRepo) Create(ctx context.Context, t *domain.Template) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notification_templates (id, type, subject_template, body_template, channels, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		t.ID, t.Type, t.SubjectTemplate, t.BodyTemplate, t.Channels, t.CreatedAt)
	return err
}

func (r *tmplRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	t := &domain.Template{}
	err := r.pool.QueryRow(ctx, `SELECT id, type, subject_template, body_template, channels, created_at FROM notification_templates WHERE id=$1`, id).
		Scan(&t.ID, &t.Type, &t.SubjectTemplate, &t.BodyTemplate, &t.Channels, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrTemplateNotFound
	}
	return t, err
}

func (r *tmplRepo) GetByType(ctx context.Context, typ string) (*domain.Template, error) {
	t := &domain.Template{}
	err := r.pool.QueryRow(ctx, `SELECT id, type, subject_template, body_template, channels, created_at FROM notification_templates WHERE type=$1`, typ).
		Scan(&t.ID, &t.Type, &t.SubjectTemplate, &t.BodyTemplate, &t.Channels, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrTemplateNotFound
	}
	return t, err
}

func (r *tmplRepo) List(ctx context.Context) ([]*domain.Template, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, type, subject_template, body_template, channels, created_at FROM notification_templates ORDER BY type ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ts []*domain.Template
	for rows.Next() {
		t := &domain.Template{}
		if err := rows.Scan(&t.ID, &t.Type, &t.SubjectTemplate, &t.BodyTemplate, &t.Channels, &t.CreatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	if ts == nil {
		ts = []*domain.Template{}
	}
	return ts, nil
}

func (r *tmplRepo) Update(ctx context.Context, t *domain.Template) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notification_templates SET subject_template=$2, body_template=$3, channels=$4 WHERE id=$1`, t.ID, t.SubjectTemplate, t.BodyTemplate, t.Channels)
	return err
}

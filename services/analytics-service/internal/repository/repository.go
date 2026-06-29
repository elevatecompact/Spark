package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/analytics-service/internal/domain"
)

type TrackedEventRepository interface {
	Insert(ctx context.Context, event *domain.TrackedEvent) error
	InsertBatch(ctx context.Context, events []*domain.TrackedEvent) error
	Query(ctx context.Context, metricName string, start, end time.Time) ([]*domain.MetricAggregate, error)
}

type DashboardRepository interface {
	GetByUserAndType(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType) (*domain.Dashboard, error)
	Upsert(ctx context.Context, dash *domain.Dashboard) error
}

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Report, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Report, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReportStatus, downloadURL string) error
}

type ReportTemplateRepository interface {
	List(ctx context.Context) ([]*domain.ReportTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ReportTemplate, error)
}

type FunnelRepository interface {
	Create(ctx context.Context, funnel *domain.Funnel) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Funnel, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Funnel, error)
	UpdateResults(ctx context.Context, id uuid.UUID, results json.RawMessage) error
}

type eventRepository struct {
	pool *pgxpool.Pool
}

type dashboardRepository struct {
	pool *pgxpool.Pool
}

type reportRepository struct {
	pool *pgxpool.Pool
}

type templateRepository struct {
	pool *pgxpool.Pool
}

type funnelRepository struct {
	pool *pgxpool.Pool
}

func NewTrackedEventRepository(pool *pgxpool.Pool) TrackedEventRepository {
	return &eventRepository{pool: pool}
}

func NewDashboardRepository(pool *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{pool: pool}
}

func NewReportRepository(pool *pgxpool.Pool) ReportRepository {
	return &reportRepository{pool: pool}
}

func NewReportTemplateRepository(pool *pgxpool.Pool) ReportTemplateRepository {
	return &templateRepository{pool: pool}
}

func NewFunnelRepository(pool *pgxpool.Pool) FunnelRepository {
	return &funnelRepository{pool: pool}
}

func (r *eventRepository) Insert(ctx context.Context, event *domain.TrackedEvent) error {
	query := `INSERT INTO tracked_events (id, event_name, user_id, session_id, properties, context, event_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, event.ID, event.EventName, event.UserID, event.SessionID, event.Properties, event.Context, event.EventTime, event.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	return nil
}

func (r *eventRepository) InsertBatch(ctx context.Context, events []*domain.TrackedEvent) error {
	if len(events) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, e := range events {
		batch.Queue(`INSERT INTO tracked_events (id, event_name, user_id, session_id, properties, context, event_time, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, e.ID, e.EventName, e.UserID, e.SessionID, e.Properties, e.Context, e.EventTime, e.CreatedAt)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range events {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("failed to insert batch event: %w", err)
		}
	}
	return nil
}

func (r *eventRepository) Query(ctx context.Context, metricName string, start, end time.Time) ([]*domain.MetricAggregate, error) {
	rows, err := r.pool.Query(ctx, `SELECT metric_name, time_bucket, dimensions, count, sum, avg, p50, p95, p99
		FROM metric_aggregates WHERE metric_name = $1 AND time_bucket >= $2 AND time_bucket <= $3 ORDER BY time_bucket ASC`, metricName, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var aggregates []*domain.MetricAggregate
	for rows.Next() {
		a := &domain.MetricAggregate{}
		var dims []byte
		if err := rows.Scan(&a.MetricName, &a.TimeBucket, &dims, &a.Count, &a.Sum, &a.Avg, &a.P50, &a.P95, &a.P99); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}
		a.Dimensions = json.RawMessage(dims)
		aggregates = append(aggregates, a)
	}
	if aggregates == nil {
		aggregates = []*domain.MetricAggregate{}
	}
	return aggregates, nil
}

func (r *dashboardRepository) GetByUserAndType(ctx context.Context, userID uuid.UUID, dashType domain.DashboardType) (*domain.Dashboard, error) {
	query := `SELECT id, user_id, dash_type, config, data, cache_until, created_at, updated_at
		FROM dashboards WHERE user_id = $1 AND dash_type = $2`
	d := &domain.Dashboard{}
	var config, data []byte
	err := r.pool.QueryRow(ctx, query, userID, dashType).Scan(&d.ID, &d.UserID, &d.DashType, &config, &data, &d.CacheUntil, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDashboardNotFound
		}
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	d.Config = json.RawMessage(config)
	d.Data = json.RawMessage(data)
	return d, nil
}

func (r *dashboardRepository) Upsert(ctx context.Context, d *domain.Dashboard) error {
	query := `INSERT INTO dashboards (id, user_id, dash_type, config, data, cache_until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, dash_type) DO UPDATE SET config=$4, data=$5, cache_until=$6, updated_at=NOW()`
	_, err := r.pool.Exec(ctx, query, d.ID, d.UserID, d.DashType, d.Config, d.Data, d.CacheUntil, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert dashboard: %w", err)
	}
	return nil
}

func (r *reportRepository) Create(ctx context.Context, report *domain.Report) error {
	query := `INSERT INTO reports (id, user_id, name, type, config, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, report.ID, report.UserID, report.Name, report.Type, report.Config, report.Status, report.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}
	return nil
}

func (r *reportRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	query := `SELECT id, user_id, name, type, config, status, download_url, created_at, completed_at
		FROM reports WHERE id = $1`
	rep := &domain.Report{}
	var config []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&rep.ID, &rep.UserID, &rep.Name, &rep.Type, &config, &rep.Status, &rep.DownloadURL, &rep.CreatedAt, &rep.CompletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrReportNotFound
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	rep.Config = json.RawMessage(config)
	return rep, nil
}

func (r *reportRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Report, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, name, type, config, status, download_url, created_at, completed_at
		FROM reports WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}
	defer rows.Close()

	var reports []*domain.Report
	for rows.Next() {
		rep := &domain.Report{}
		var config []byte
		if err := rows.Scan(&rep.ID, &rep.UserID, &rep.Name, &rep.Type, &config, &rep.Status, &rep.DownloadURL, &rep.CreatedAt, &rep.CompletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		rep.Config = json.RawMessage(config)
		reports = append(reports, rep)
	}
	if reports == nil {
		reports = []*domain.Report{}
	}
	return reports, nil
}

func (r *reportRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReportStatus, downloadURL string) error {
	now := time.Now().UTC()
	tag, err := r.pool.Exec(ctx, `UPDATE reports SET status=$2, download_url=$3, completed_at=CASE WHEN $2='ready' THEN $4 ELSE completed_at END WHERE id=$1`, id, status, downloadURL, now)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrReportNotFound
	}
	return nil
}

func (r *templateRepository) List(ctx context.Context) ([]*domain.ReportTemplate, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, config, created_at FROM report_templates ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	var templates []*domain.ReportTemplate
	for rows.Next() {
		t := &domain.ReportTemplate{}
		var config []byte
		if err := rows.Scan(&t.ID, &t.Name, &config, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		t.Config = json.RawMessage(config)
		templates = append(templates, t)
	}
	if templates == nil {
		templates = []*domain.ReportTemplate{}
	}
	return templates, nil
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReportTemplate, error) {
	query := `SELECT id, name, config, created_at FROM report_templates WHERE id = $1`
	t := &domain.ReportTemplate{}
	var config []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &config, &t.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	t.Config = json.RawMessage(config)
	return t, nil
}

func (r *funnelRepository) Create(ctx context.Context, funnel *domain.Funnel) error {
	query := `INSERT INTO funnels (id, user_id, name, steps, results, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, funnel.ID, funnel.UserID, funnel.Name, funnel.Steps, funnel.Results, funnel.CreatedAt, funnel.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create funnel: %w", err)
	}
	return nil
}

func (r *funnelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Funnel, error) {
	query := `SELECT id, user_id, name, steps, results, created_at, updated_at FROM funnels WHERE id = $1`
	f := &domain.Funnel{}
	var steps, results []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&f.ID, &f.UserID, &f.Name, &steps, &results, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrFunnelNotFound
		}
		return nil, fmt.Errorf("failed to get funnel: %w", err)
	}
	f.Steps = json.RawMessage(steps)
	f.Results = json.RawMessage(results)
	return f, nil
}

func (r *funnelRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Funnel, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, name, steps, results, created_at, updated_at FROM funnels WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list funnels: %w", err)
	}
	defer rows.Close()

	var funnels []*domain.Funnel
	for rows.Next() {
		f := &domain.Funnel{}
		var steps, results []byte
		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &steps, &results, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan funnel: %w", err)
		}
		f.Steps = json.RawMessage(steps)
		f.Results = json.RawMessage(results)
		funnels = append(funnels, f)
	}
	if funnels == nil {
		funnels = []*domain.Funnel{}
	}
	return funnels, nil
}

func (r *funnelRepository) UpdateResults(ctx context.Context, id uuid.UUID, results json.RawMessage) error {
	tag, err := r.pool.Exec(ctx, `UPDATE funnels SET results=$2, updated_at=NOW() WHERE id=$1`, id, results)
	if err != nil {
		return fmt.Errorf("failed to update funnel results: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrFunnelNotFound
	}
	return nil
}

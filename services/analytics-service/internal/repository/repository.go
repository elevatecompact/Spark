package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/analytics-service/internal/domain"
)

type TrackedEventRepository interface {
	Insert(ctx context.Context, event *domain.TrackedEvent) error
	InsertBatch(ctx context.Context, events []*domain.TrackedEvent) error
	Query(ctx context.Context, metricName string, start, end time.Time) ([]*domain.MetricAggregate, error)
	CountByEventSince(ctx context.Context, eventName string, sinceMinutes int) (int64, error)
	UniqueUsersSince(ctx context.Context, sinceMinutes int) (int64, error)
	TopEventNamesSince(ctx context.Context, sinceMinutes, limit int) ([]string, int64, error)
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
	pool   *pgxpool.Pool
	chConn driver.Conn
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

func NewTrackedEventRepository(pool *pgxpool.Pool, chConn driver.Conn) TrackedEventRepository {
	return &eventRepository{pool: pool, chConn: chConn}
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

// OpenClickHouse opens a ClickHouse connection using the standard DSN format
// (clickhouse://user:password@host:port/database). When the URL is empty a nil
// connection is returned; callers should treat that as "ClickHouse disabled".
func OpenClickHouse(ctx context.Context, dsn string) (driver.Conn, error) {
	if dsn == "" {
		return nil, nil
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{dsn},
		Auth: clickhouse.Auth{
			Database: "default",
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("open clickhouse: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping clickhouse: %w", err)
	}
	return conn, nil
}

func (r *eventRepository) Insert(ctx context.Context, event *domain.TrackedEvent) error {
	query := `INSERT INTO tracked_events (id, event_name, user_id, session_id, properties, context, event_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, event.ID, event.EventName, event.UserID, event.SessionID, event.Properties, event.Context, event.EventTime, event.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	if r.chConn != nil {
		_ = r.chConn.AsyncInsert(ctx, fmt.Sprintf(
			"INSERT INTO spark.tracked_events (id, event_name, user_id, session_id, properties, event_time) VALUES ('%s','%s','%s','%s','%s','%s')",
			event.ID, escapeCH(event.EventName), event.UserID, escapeCH(event.SessionID), escapeCH(string(event.Properties)), event.EventTime.UTC().Format("2006-01-02 15:04:05.000"),
		), false)
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
	if r.chConn != nil {
		var buf []string
		for _, e := range events {
			buf = append(buf, fmt.Sprintf(
				"('%s','%s','%s','%s','%s','%s')",
				e.ID, escapeCH(e.EventName), e.UserID, escapeCH(e.SessionID), escapeCH(string(e.Properties)), e.EventTime.UTC().Format("2006-01-02 15:04:05.000"),
			))
		}
		stmt := "INSERT INTO spark.tracked_events (id, event_name, user_id, session_id, properties, event_time) VALUES " + joinComma(buf)
		_ = r.chConn.AsyncInsert(ctx, stmt, false)
	}
	return nil
}

func (r *eventRepository) Query(ctx context.Context, metricName string, start, end time.Time) ([]*domain.MetricAggregate, error) {
	// Try ClickHouse first for richer aggregates; fall back to Postgres.
	if r.chConn != nil {
		rows, err := r.chConn.Query(ctx, `
SELECT
    event_name AS metric_name,
    toStartOfMinute(event_time) AS time_bucket,
    count() AS count,
    sum(0) AS sum,
    avg(0) AS avg,
    quantile(0.5)(0) AS p50,
    quantile(0.95)(0) AS p95,
    quantile(0.99)(0) AS p99
FROM spark.tracked_events
WHERE event_name = ? AND event_time BETWEEN ? AND ?
GROUP BY metric_name, time_bucket
ORDER BY time_bucket ASC`, metricName, start, end)
		if err == nil {
			defer rows.Close()
			var out []*domain.MetricAggregate
			for rows.Next() {
				a := &domain.MetricAggregate{}
				if err := rows.ScanStruct(&a); err == nil {
					a.Dimensions = json.RawMessage("{}")
					out = append(out, a)
				}
			}
			if out != nil {
				return out, nil
			}
		}
	}

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

// CountByEventSince returns the number of events with the given name in the
// last `sinceMinutes` minutes. Backed by ClickHouse when available.
func (r *eventRepository) CountByEventSince(ctx context.Context, eventName string, sinceMinutes int) (int64, error) {
	if r.chConn != nil {
		var count uint64
		err := r.chConn.QueryRow(ctx, `SELECT count() FROM spark.tracked_events WHERE event_name = ? AND event_time >= now() - INTERVAL ? MINUTE`, eventName, sinceMinutes).Scan(&count)
		if err == nil {
			return int64(count), nil
		}
	}
	return r.countByEventPG(ctx, eventName, sinceMinutes)
}

func (r *eventRepository) countByEventPG(ctx context.Context, eventName string, sinceMinutes int) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT count(*) FROM tracked_events WHERE event_name = $1 AND event_time >= NOW() - ($2::int * INTERVAL '1 minute')`, eventName, sinceMinutes).Scan(&count)
	return count, err
}

// UniqueUsersSince returns the number of distinct users active in the window.
func (r *eventRepository) UniqueUsersSince(ctx context.Context, sinceMinutes int) (int64, error) {
	if r.chConn != nil {
		var count uint64
		err := r.chConn.QueryRow(ctx, `SELECT uniqExact(user_id) FROM spark.tracked_events WHERE event_time >= now() - INTERVAL ? MINUTE`, sinceMinutes).Scan(&count)
		if err == nil {
			return int64(count), nil
		}
	}
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT count(DISTINCT user_id) FROM tracked_events WHERE event_time >= NOW() - ($1::int * INTERVAL '1 minute')`, sinceMinutes).Scan(&count)
	return count, err
}

// TopEventNamesSince returns the most frequent event names in the window
// along with the total number of events.
func (r *eventRepository) TopEventNamesSince(ctx context.Context, sinceMinutes, limit int) ([]string, int64, error) {
	if r.chConn != nil {
		rows, err := r.chConn.Query(ctx, `SELECT event_name, count() AS c FROM spark.tracked_events WHERE event_time >= now() - INTERVAL ? MINUTE GROUP BY event_name ORDER BY c DESC LIMIT ?`, sinceMinutes, limit)
		if err == nil {
			defer rows.Close()
			var names []string
			for rows.Next() {
				var name string
				var c uint64
				if err := rows.Scan(&name, &c); err == nil {
					names = append(names, name)
				}
			}
			var total uint64
			_ = r.chConn.QueryRow(ctx, `SELECT count() FROM spark.tracked_events WHERE event_time >= now() - INTERVAL ? MINUTE`, sinceMinutes).Scan(&total)
			return names, int64(total), nil
		}
	}
	rows, err := r.pool.Query(ctx, `SELECT event_name, count(*) AS c FROM tracked_events WHERE event_time >= NOW() - ($1::int * INTERVAL '1 minute') GROUP BY event_name ORDER BY c DESC LIMIT $2`, sinceMinutes, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		var c int64
		if err := rows.Scan(&name, &c); err == nil {
			names = append(names, name)
		}
	}
	var total int64
	_ = r.pool.QueryRow(ctx, `SELECT count(*) FROM tracked_events WHERE event_time >= NOW() - ($1::int * INTERVAL '1 minute')`, sinceMinutes).Scan(&total)
	return names, total, nil
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

func escapeCH(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' || c == '\\' {
			out = append(out, '\\', c)
			continue
		}
		if c == '\n' {
			out = append(out, '\\', 'n')
			continue
		}
		if c == '\r' {
			out = append(out, '\\', 'r')
			continue
		}
		out = append(out, c)
	}
	return string(out)
}

func joinComma(items []string) string {
	if len(items) == 0 {
		return ""
	}
	out := items[0]
	for i := 1; i < len(items); i++ {
		out += "," + items[i]
	}
	return out
}

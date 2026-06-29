package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/moderation-service/internal/domain"
)

type ModerationRepository interface {
	CreateRule(ctx context.Context, r *domain.ModerationRule) error
	GetRule(ctx context.Context, id uuid.UUID) (*domain.ModerationRule, error)
	ListRules(ctx context.Context) ([]domain.ModerationRule, error)
	UpdateRule(ctx context.Context, r *domain.ModerationRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error

	CreateAction(ctx context.Context, a *domain.ModerationAction) error
	ListActions(ctx context.Context, userID uuid.UUID) ([]domain.ModerationAction, error)

	CreateReviewItem(ctx context.Context, r *domain.ReviewItem) error
	GetReviewItem(ctx context.Context, id uuid.UUID) (*domain.ReviewItem, error)
	ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewItem, error)
	ApproveReview(ctx context.Context, id uuid.UUID, resolution string) error
	RejectReview(ctx context.Context, id uuid.UUID, resolution string) error
	GetQueueStats(ctx context.Context) ([]domain.QueueStats, error)

	CreateReport(ctx context.Context, r *domain.ContentReport) error
	GetReport(ctx context.Context, id uuid.UUID) (*domain.ContentReport, error)
	ListReports(ctx context.Context, status domain.ReportStatus) ([]domain.ContentReport, error)

	GetAdminStats(ctx context.Context) (*domain.AdminStats, error)
}

type modRepo struct{ pool *pgxpool.Pool }

func NewModerationRepository(pool *pgxpool.Pool) ModerationRepository { return &modRepo{pool} }

func (r *modRepo) CreateRule(ctx context.Context, rule *domain.ModerationRule) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO moderation_rules (id, name, category, severity, conditions, is_active, priority, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`,
		rule.ID, rule.Name, rule.Category, rule.Severity, rule.Conditions, rule.IsActive, rule.Priority)
	return err
}

func (r *modRepo) GetRule(ctx context.Context, id uuid.UUID) (*domain.ModerationRule, error) {
	rule := &domain.ModerationRule{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, category, severity, conditions, is_active, priority, created_at FROM moderation_rules WHERE id=$1`, id).
		Scan(&rule.ID, &rule.Name, &rule.Category, &rule.Severity, &rule.Conditions, &rule.IsActive, &rule.Priority, &rule.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func (r *modRepo) ListRules(ctx context.Context) ([]domain.ModerationRule, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, category, severity, conditions, is_active, priority, created_at FROM moderation_rules ORDER BY priority ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []domain.ModerationRule
	for rows.Next() {
		var rule domain.ModerationRule
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.Category, &rule.Severity, &rule.Conditions, &rule.IsActive, &rule.Priority, &rule.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if rules == nil {
		rules = []domain.ModerationRule{}
	}
	return rules, nil
}

func (r *modRepo) UpdateRule(ctx context.Context, rule *domain.ModerationRule) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE moderation_rules SET name=$2, category=$3, severity=$4, conditions=$5, is_active=$6, priority=$7 WHERE id=$1`,
		rule.ID, rule.Name, rule.Category, rule.Severity, rule.Conditions, rule.IsActive, rule.Priority)
	return err
}

func (r *modRepo) DeleteRule(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM moderation_rules WHERE id=$1`, id)
	return err
}

func (r *modRepo) CreateAction(ctx context.Context, a *domain.ModerationAction) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO moderation_actions (id, user_id, content_id, rule_id, action_type, status, applied_by, reason, duration, applied_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())`,
		a.ID, a.UserID, a.ContentID, a.RuleID, a.ActionType, a.Status, a.AppliedBy, a.Reason, a.Duration)
	return err
}

func (r *modRepo) ListActions(ctx context.Context, userID uuid.UUID) ([]domain.ModerationAction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, content_id, rule_id, action_type, status, applied_by, reason, duration, applied_at FROM moderation_actions WHERE user_id=$1 ORDER BY applied_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var actions []domain.ModerationAction
	for rows.Next() {
		var a domain.ModerationAction
		if err := rows.Scan(&a.ID, &a.UserID, &a.ContentID, &a.RuleID, &a.ActionType, &a.Status, &a.AppliedBy, &a.Reason, &a.Duration, &a.AppliedAt); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	if actions == nil {
		actions = []domain.ModerationAction{}
	}
	return actions, nil
}

func (r *modRepo) CreateReviewItem(ctx context.Context, ri *domain.ReviewItem) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO review_queue (id, content_type, content_id, flagged_by, reasons, assigned_moderator, status, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`,
		ri.ID, ri.ContentType, ri.ContentID, ri.FlaggedBy, ri.Reasons, ri.AssignedModerator, ri.Status)
	return err
}

func (r *modRepo) GetReviewItem(ctx context.Context, id uuid.UUID) (*domain.ReviewItem, error) {
	ri := &domain.ReviewItem{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, content_type, content_id, flagged_by, reasons, assigned_moderator, status, resolution, resolved_at, created_at FROM review_queue WHERE id=$1`, id).
		Scan(&ri.ID, &ri.ContentType, &ri.ContentID, &ri.FlaggedBy, &ri.Reasons, &ri.AssignedModerator, &ri.Status, &ri.Resolution, &ri.ResolvedAt, &ri.CreatedAt)
	if err != nil {
		return nil, err
	}
	return ri, nil
}

func (r *modRepo) ListReviewQueue(ctx context.Context, status domain.ReviewStatus) ([]domain.ReviewItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, content_type, content_id, flagged_by, reasons, assigned_moderator, status, resolution, resolved_at, created_at FROM review_queue WHERE status=$1 ORDER BY created_at ASC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.ReviewItem
	for rows.Next() {
		var ri domain.ReviewItem
		if err := rows.Scan(&ri.ID, &ri.ContentType, &ri.ContentID, &ri.FlaggedBy, &ri.Reasons, &ri.AssignedModerator, &ri.Status, &ri.Resolution, &ri.ResolvedAt, &ri.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, ri)
	}
	if items == nil {
		items = []domain.ReviewItem{}
	}
	return items, nil
}

func (r *modRepo) ApproveReview(ctx context.Context, id uuid.UUID, resolution string) error {
	_, err := r.pool.Exec(ctx, `UPDATE review_queue SET status='resolved', resolution=$2, resolved_at=NOW() WHERE id=$1`, id, resolution)
	return err
}

func (r *modRepo) RejectReview(ctx context.Context, id uuid.UUID, resolution string) error {
	_, err := r.pool.Exec(ctx, `UPDATE review_queue SET status='resolved', resolution=$2, resolved_at=NOW() WHERE id=$1`, id, resolution)
	return err
}

func (r *modRepo) GetQueueStats(ctx context.Context) ([]domain.QueueStats, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT content_type, COUNT(*) FILTER (WHERE status='pending') as pending, COUNT(*) FILTER (WHERE status='reviewing') as reviewing, COUNT(*) FILTER (WHERE status='resolved') as resolved FROM review_queue GROUP BY content_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []domain.QueueStats
	for rows.Next() {
		var s domain.QueueStats
		if err := rows.Scan(&s.ContentType, &s.Pending, &s.Reviewing, &s.Resolved); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if stats == nil {
		stats = []domain.QueueStats{}
	}
	return stats, nil
}

func (r *modRepo) CreateReport(ctx context.Context, report *domain.ContentReport) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO content_reports (id, reporter_id, content_type, content_id, reason, description, status, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`,
		report.ID, report.ReporterID, report.ContentType, report.ContentID, report.Reason, report.Description, report.Status)
	return err
}

func (r *modRepo) GetReport(ctx context.Context, id uuid.UUID) (*domain.ContentReport, error) {
	report := &domain.ContentReport{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, reporter_id, content_type, content_id, reason, description, status, created_at FROM content_reports WHERE id=$1`, id).
		Scan(&report.ID, &report.ReporterID, &report.ContentType, &report.ContentID, &report.Reason, &report.Description, &report.Status, &report.CreatedAt)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (r *modRepo) ListReports(ctx context.Context, status domain.ReportStatus) ([]domain.ContentReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, reporter_id, content_type, content_id, reason, description, status, created_at FROM content_reports WHERE status=$1 ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reports []domain.ContentReport
	for rows.Next() {
		var rep domain.ContentReport
		if err := rows.Scan(&rep.ID, &rep.ReporterID, &rep.ContentType, &rep.ContentID, &rep.Reason, &rep.Description, &rep.Status, &rep.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	if reports == nil {
		reports = []domain.ContentReport{}
	}
	return reports, nil
}

func (r *modRepo) GetAdminStats(ctx context.Context) (*domain.AdminStats, error) {
	stats := &domain.AdminStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT (SELECT COUNT(*) FROM moderation_actions) as total_actions,
		        (SELECT COUNT(*) FROM review_queue) as queue_depth,
		        (SELECT COUNT(*) FROM moderation_actions WHERE applied_by='automated') as auto_actions`).
		Scan(&stats.TotalActions, &stats.QueueDepth, &stats.AutoActionRate)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

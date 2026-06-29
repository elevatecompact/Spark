package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/licensing-service/internal/domain"
)

type LicensingRepository struct {
	pool *pgxpool.Pool
}

func NewLicensingRepository(pool *pgxpool.Pool) *LicensingRepository {
	return &LicensingRepository{pool: pool}
}

func (r *LicensingRepository) CreateLicense(ctx context.Context, l *domain.License) error {
	territory, _ := json.Marshal(l.Territory)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO licenses (id, rights_holder_id, licensee_id, content_id, type, scope, territory, start_date, end_date, auto_renew, rate_type, rate_cents, revenue_share_percent, min_guarantee_cents, status, terms_url, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`, l.ID, l.RightsHolderID, l.LicenseeID, l.ContentID, string(l.Type), string(l.Scope), territory, l.StartDate, l.EndDate, l.AutoRenew, string(l.RateType), l.RateCents, l.RevenueSharePercent, l.MinGuaranteeCents, string(l.Status), l.TermsURL, l.CreatedAt)
	return err
}

func (r *LicensingRepository) GetLicense(ctx context.Context, id uuid.UUID) (*domain.License, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, rights_holder_id, licensee_id, content_id, type, scope, territory, start_date, end_date, auto_renew, rate_type, rate_cents, revenue_share_percent, min_guarantee_cents, status, terms_url, created_at
		FROM licenses WHERE id=$1
	`, id)
	l := &domain.License{}
	var territory []byte
	err := row.Scan(&l.ID, &l.RightsHolderID, &l.LicenseeID, &l.ContentID, &l.Type, &l.Scope, &territory, &l.StartDate, &l.EndDate, &l.AutoRenew, &l.RateType, &l.RateCents, &l.RevenueSharePercent, &l.MinGuaranteeCents, &l.Status, &l.TermsURL, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(territory, &l.Territory)
	return l, nil
}

func (r *LicensingRepository) UpdateLicense(ctx context.Context, l *domain.License) error {
	territory, _ := json.Marshal(l.Territory)
	_, err := r.pool.Exec(ctx, `
		UPDATE licenses SET type=$3, scope=$4, territory=$5, start_date=$6, end_date=$7, auto_renew=$8, rate_type=$9, rate_cents=$10, revenue_share_percent=$11, min_guarantee_cents=$12, status=$13, terms_url=$14
		WHERE id=$1 AND rights_holder_id=$2
	`, l.ID, l.RightsHolderID, string(l.Type), string(l.Scope), territory, l.StartDate, l.EndDate, l.AutoRenew, string(l.RateType), l.RateCents, l.RevenueSharePercent, l.MinGuaranteeCents, string(l.Status), l.TermsURL)
	return err
}

func (r *LicensingRepository) UpdateLicenseStatus(ctx context.Context, id uuid.UUID, status domain.LicenseStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE licenses SET status=$2 WHERE id=$1`, id, string(status))
	return err
}

func (r *LicensingRepository) DeleteLicense(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM licenses WHERE id=$1`, id)
	return err
}

func (r *LicensingRepository) ListLicenses(ctx context.Context, rightsHolderID, licenseeID, contentID *uuid.UUID, limit, offset int) ([]domain.License, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if rightsHolderID != nil {
		where += " AND rights_holder_id=$" + string(rune('0'+idx))
		args = append(args, *rightsHolderID)
		idx++
	}
	if licenseeID != nil {
		where += " AND licensee_id=$" + string(rune('0'+idx))
		args = append(args, *licenseeID)
		idx++
	}
	if contentID != nil {
		where += " AND content_id=$" + string(rune('0'+idx))
		args = append(args, *contentID)
		idx++
	}
	args = append(args, limit, offset)
	q := `SELECT id, rights_holder_id, licensee_id, content_id, type, scope, territory, start_date, end_date, auto_renew, rate_type, rate_cents, revenue_share_percent, min_guarantee_cents, status, terms_url, created_at
	      FROM licenses ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.License
	for rows.Next() {
		var l domain.License
		var territory []byte
		if err := rows.Scan(&l.ID, &l.RightsHolderID, &l.LicenseeID, &l.ContentID, &l.Type, &l.Scope, &territory, &l.StartDate, &l.EndDate, &l.AutoRenew, &l.RateType, &l.RateCents, &l.RevenueSharePercent, &l.MinGuaranteeCents, &l.Status, &l.TermsURL, &l.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(territory, &l.Territory)
		res = append(res, l)
	}
	return res, nil
}

func (r *LicensingRepository) RegisterContentRight(ctx context.Context, cr *domain.ContentRight) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO content_rights (id, content_id, rights_holder_id, license_id, restrictions, registered_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, cr.ID, cr.ContentID, cr.RightsHolderID, cr.LicenseID, cr.Restrictions, cr.RegisteredAt)
	return err
}

func (r *LicensingRepository) GetContentRights(ctx context.Context, contentID uuid.UUID) (*domain.ContentRight, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, content_id, rights_holder_id, license_id, restrictions, registered_at
		FROM content_rights WHERE content_id=$1
	`, contentID)
	cr := &domain.ContentRight{}
	err := row.Scan(&cr.ID, &cr.ContentID, &cr.RightsHolderID, &cr.LicenseID, &cr.Restrictions, &cr.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (r *LicensingRepository) RecordUsage(ctx context.Context, u *domain.UsageLog) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO usage_log (id, license_id, content_id, usage_type, context, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, u.ID, u.LicenseID, u.ContentID, string(u.UsageType), u.Context, u.RecordedAt)
	return err
}

func (r *LicensingRepository) GetUsageByLicense(ctx context.Context, licenseID uuid.UUID, start, end time.Time) ([]domain.UsageLog, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, license_id, content_id, usage_type, context, recorded_at
		FROM usage_log WHERE license_id=$1 AND recorded_at>=$2 AND recorded_at<=$3 ORDER BY recorded_at DESC
	`, licenseID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.UsageLog
	for rows.Next() {
		var u domain.UsageLog
		if err := rows.Scan(&u.ID, &u.LicenseID, &u.ContentID, &u.UsageType, &u.Context, &u.RecordedAt); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func (r *LicensingRepository) GetUsageByRightsHolder(ctx context.Context, rightsHolderID uuid.UUID, start, end time.Time) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM usage_log ul JOIN licenses l ON ul.license_id=l.id
		WHERE l.rights_holder_id=$1 AND ul.recorded_at>=$2 AND ul.recorded_at<=$3
	`, rightsHolderID, start, end).Scan(&count)
	return count, err
}

func (r *LicensingRepository) GetUsageByContent(ctx context.Context, contentID uuid.UUID, limit, offset int) ([]domain.UsageLog, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, license_id, content_id, usage_type, context, recorded_at
		FROM usage_log WHERE content_id=$1 ORDER BY recorded_at DESC LIMIT $2 OFFSET $3
	`, contentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.UsageLog
	for rows.Next() {
		var u domain.UsageLog
		if err := rows.Scan(&u.ID, &u.LicenseID, &u.ContentID, &u.UsageType, &u.Context, &u.RecordedAt); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func (r *LicensingRepository) CreateRoyaltyStatement(ctx context.Context, rs *domain.RoyaltyStatement) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO royalty_statements (id, license_id, rights_holder_id, period_start, period_end, usage_count, rate_applied, total_cents, status, paid_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, rs.ID, rs.LicenseID, rs.RightsHolderID, rs.PeriodStart, rs.PeriodEnd, rs.UsageCount, rs.RateApplied, rs.TotalCents, string(rs.Status), rs.PaidAt, rs.CreatedAt)
	return err
}

func (r *LicensingRepository) GetRoyaltyStatements(ctx context.Context, rightsHolderID, licenseID *uuid.UUID, limit, offset int) ([]domain.RoyaltyStatement, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if rightsHolderID != nil {
		where += " AND rights_holder_id=$" + string(rune('0'+idx))
		args = append(args, *rightsHolderID)
		idx++
	}
	if licenseID != nil {
		where += " AND license_id=$" + string(rune('0'+idx))
		args = append(args, *licenseID)
		idx++
	}
	args = append(args, limit, offset)
	q := `SELECT id, license_id, rights_holder_id, period_start, period_end, usage_count, rate_applied, total_cents, status, paid_at, created_at
	      FROM royalty_statements ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.RoyaltyStatement
	for rows.Next() {
		var rs domain.RoyaltyStatement
		if err := rows.Scan(&rs.ID, &rs.LicenseID, &rs.RightsHolderID, &rs.PeriodStart, &rs.PeriodEnd, &rs.UsageCount, &rs.RateApplied, &rs.TotalCents, &rs.Status, &rs.PaidAt, &rs.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rs)
	}
	return res, nil
}

func (r *LicensingRepository) GetPendingRoyalties(ctx context.Context) ([]domain.RoyaltyStatement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, license_id, rights_holder_id, period_start, period_end, usage_count, rate_applied, total_cents, status, paid_at, created_at
		FROM royalty_statements WHERE status='pending' ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.RoyaltyStatement
	for rows.Next() {
		var rs domain.RoyaltyStatement
		if err := rows.Scan(&rs.ID, &rs.LicenseID, &rs.RightsHolderID, &rs.PeriodStart, &rs.PeriodEnd, &rs.UsageCount, &rs.RateApplied, &rs.TotalCents, &rs.Status, &rs.PaidAt, &rs.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rs)
	}
	return res, nil
}

func (r *LicensingRepository) UpdateRoyaltyStatus(ctx context.Context, id uuid.UUID, status domain.RoyaltyStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE royalty_statements SET status=$2, paid_at=CASE WHEN $2='paid' THEN NOW() ELSE NULL END WHERE id=$1`, id, string(status))
	return err
}

func (r *LicensingRepository) GetComplianceReport(ctx context.Context) (*domain.ComplianceReport, error) {
	var report domain.ComplianceReport
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM licenses`).Scan(&report.TotalLicenses)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM licenses WHERE status='active'`).Scan(&report.ActiveLicenses)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM licenses WHERE status='active' AND end_date <= NOW() + INTERVAL '30 days'`).Scan(&report.ExpiringLicenses)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM usage_log`).Scan(&report.UsageRecords)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM royalty_statements WHERE status='pending'`).Scan(&report.PendingRoyalties)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM royalty_statements WHERE status='paid'`).Scan(&report.PaidRoyalties)
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(total_cents),0) FROM royalty_statements WHERE status='paid'`).Scan(&report.TotalRoyaltyCents)
	return &report, nil
}

func (r *LicensingRepository) CalculateProjectedRoyalty(ctx context.Context, licenseID uuid.UUID) (*domain.RoyaltyStatement, error) {
	l, err := r.GetLicense(ctx, licenseID)
	if err != nil {
		return nil, err
	}
	count, err := r.GetUsageByLicense(ctx, licenseID, l.StartDate, l.EndDate)
	if err != nil {
		return nil, err
	}
	usageCount := int64(len(count))
	totalCents := l.RateCents * usageCount
	if l.RateType == domain.RateTypeFlat {
		totalCents = l.RateCents
	}
	if l.MinGuaranteeCents > 0 && totalCents < l.MinGuaranteeCents {
		totalCents = l.MinGuaranteeCents
	}
	return &domain.RoyaltyStatement{
		LicenseID:      licenseID,
		RightsHolderID: l.RightsHolderID,
		UsageCount:     usageCount,
		RateApplied:    l.RateCents,
		TotalCents:     totalCents,
		Status:         domain.RoyaltyPending,
		CreatedAt:      l.CreatedAt,
	}, nil
}

func (r *LicensingRepository) GetActiveLicenseByContent(ctx context.Context, contentID uuid.UUID) (*domain.License, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, rights_holder_id, licensee_id, content_id, type, scope, territory, start_date, end_date, auto_renew, rate_type, rate_cents, revenue_share_percent, min_guarantee_cents, status, terms_url, created_at
		FROM licenses WHERE content_id=$1 AND status='active' AND start_date<=NOW() AND end_date>=NOW() ORDER BY created_at DESC LIMIT 1
	`, contentID)
	l := &domain.License{}
	var territory []byte
	err := row.Scan(&l.ID, &l.RightsHolderID, &l.LicenseeID, &l.ContentID, &l.Type, &l.Scope, &territory, &l.StartDate, &l.EndDate, &l.AutoRenew, &l.RateType, &l.RateCents, &l.RevenueSharePercent, &l.MinGuaranteeCents, &l.Status, &l.TermsURL, &l.CreatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(territory, &l.Territory)
	return l, nil
}

package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/trust-service/internal/domain"
)

type TrustRepository struct {
	pool *pgxpool.Pool
}

func NewTrustRepository(pool *pgxpool.Pool) *TrustRepository {
	return &TrustRepository{pool: pool}
}

func (r *TrustRepository) GetReputation(ctx context.Context, userID uuid.UUID) (*domain.ReputationScore, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT user_id, overall_score, trust_level, positive_signal_weight, negative_signal_weight, score_decay_factor, model_version, calculated_at, next_recalculation_at
		FROM reputation_scores WHERE user_id=$1
	`, userID)
	rs := &domain.ReputationScore{}
	err := row.Scan(&rs.UserID, &rs.OverallScore, &rs.TrustLevel, &rs.PositiveSignalWeight, &rs.NegativeSignalWeight, &rs.ScoreDecayFactor, &rs.ModelVersion, &rs.CalculatedAt, &rs.NextRecalculationAt)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *TrustRepository) UpsertReputation(ctx context.Context, rs *domain.ReputationScore) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO reputation_scores (user_id, overall_score, trust_level, positive_signal_weight, negative_signal_weight, score_decay_factor, model_version, calculated_at, next_recalculation_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (user_id) DO UPDATE SET overall_score=EXCLUDED.overall_score, trust_level=EXCLUDED.trust_level, positive_signal_weight=EXCLUDED.positive_signal_weight, negative_signal_weight=EXCLUDED.negative_signal_weight, score_decay_factor=EXCLUDED.score_decay_factor, model_version=EXCLUDED.model_version, calculated_at=EXCLUDED.calculated_at, next_recalculation_at=EXCLUDED.next_recalculation_at
	`, rs.UserID, rs.OverallScore, string(rs.TrustLevel), rs.PositiveSignalWeight, rs.NegativeSignalWeight, rs.ScoreDecayFactor, rs.ModelVersion, rs.CalculatedAt, rs.NextRecalculationAt)
	return err
}

func (r *TrustRepository) GetSignals(ctx context.Context, userID uuid.UUID) ([]domain.TrustSignal, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, signal_type, category, weight, description, source_entity_type, source_entity_id, expires_at, recorded_at
		FROM trust_signals WHERE user_id=$1 AND (expires_at IS NULL OR expires_at>NOW()) ORDER BY recorded_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.TrustSignal
	for rows.Next() {
		var s domain.TrustSignal
		if err := rows.Scan(&s.ID, &s.UserID, &s.SignalType, &s.Category, &s.Weight, &s.Description, &s.SourceEntityType, &s.SourceEntityID, &s.ExpiresAt, &s.RecordedAt); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *TrustRepository) RecordSignal(ctx context.Context, s *domain.TrustSignal) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO trust_signals (id, user_id, signal_type, category, weight, description, source_entity_type, source_entity_id, expires_at, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, s.ID, s.UserID, string(s.SignalType), string(s.Category), s.Weight, s.Description, s.SourceEntityType, s.SourceEntityID, s.ExpiresAt, s.RecordedAt)
	return err
}

func (r *TrustRepository) CreateRiskAssessment(ctx context.Context, ra *domain.RiskAssessment) error {
	triggers, _ := json.Marshal(ra.TriggeredRules)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO risk_assessments (id, user_id, action_type, context, risk_score, risk_level, triggered_rules, recommended_action, assessed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, ra.ID, ra.UserID, ra.ActionType, ra.Context, ra.RiskScore, string(ra.RiskLevel), triggers, string(ra.RecommendedAction), ra.AssessedAt)
	return err
}

func (r *TrustRepository) GetRiskAssessment(ctx context.Context, id uuid.UUID) (*domain.RiskAssessment, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, action_type, context, risk_score, risk_level, triggered_rules, recommended_action, assessed_at
		FROM risk_assessments WHERE id=$1
	`, id)
	ra := &domain.RiskAssessment{}
	var triggers []byte
	err := row.Scan(&ra.ID, &ra.UserID, &ra.ActionType, &ra.Context, &ra.RiskScore, &ra.RiskLevel, &triggers, &ra.RecommendedAction, &ra.AssessedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(triggers, &ra.TriggeredRules)
	return ra, nil
}

func (r *TrustRepository) GetActiveRiskRules(ctx context.Context) ([]domain.RiskRule, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, category, conditions, risk_score_impact, action, is_active, priority, created_at
		FROM risk_rules WHERE is_active=true ORDER BY priority DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.RiskRule
	for rows.Next() {
		var rr domain.RiskRule
		if err := rows.Scan(&rr.ID, &rr.Name, &rr.Category, &rr.Conditions, &rr.RiskScoreImpact, &rr.Action, &rr.IsActive, &rr.Priority, &rr.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rr)
	}
	return res, nil
}

func (r *TrustRepository) CreateRiskRule(ctx context.Context, rr *domain.RiskRule) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO risk_rules (id, name, category, conditions, risk_score_impact, action, is_active, priority, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, rr.ID, rr.Name, rr.Category, rr.Conditions, rr.RiskScoreImpact, rr.Action, rr.IsActive, rr.Priority, rr.CreatedAt)
	return err
}

func (r *TrustRepository) UpdateRiskRule(ctx context.Context, rr *domain.RiskRule) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE risk_rules SET name=$2, category=$3, conditions=$4, risk_score_impact=$5, action=$6, is_active=$7, priority=$8 WHERE id=$1
	`, rr.ID, rr.Name, rr.Category, rr.Conditions, rr.RiskScoreImpact, rr.Action, rr.IsActive, rr.Priority)
	return err
}

func (r *TrustRepository) CreateFraudCase(ctx context.Context, fc *domain.FraudCase) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO fraud_cases (id, user_id, case_type, status, evidence, automated_decision, reviewed_by, resolved_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, fc.ID, fc.UserID, string(fc.CaseType), string(fc.Status), fc.Evidence, fc.AutomatedDecision, fc.ReviewedBy, fc.ResolvedAt, fc.CreatedAt)
	return err
}

func (r *TrustRepository) ListFraudCases(ctx context.Context, status *domain.FraudCaseStatus, limit, offset int) ([]domain.FraudCase, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if status != nil {
		where += " AND status=$" + string(rune('0'+idx))
		args = append(args, string(*status))
		idx++
	}
	args = append(args, limit, offset)
	q := `SELECT id, user_id, case_type, status, evidence, automated_decision, reviewed_by, resolved_at, created_at
	      FROM fraud_cases ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.FraudCase
	for rows.Next() {
		var fc domain.FraudCase
		if err := rows.Scan(&fc.ID, &fc.UserID, &fc.CaseType, &fc.Status, &fc.Evidence, &fc.AutomatedDecision, &fc.ReviewedBy, &fc.ResolvedAt, &fc.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, fc)
	}
	return res, nil
}

func (r *TrustRepository) UpdateFraudCaseStatus(ctx context.Context, id uuid.UUID, status domain.FraudCaseStatus, reviewedBy *uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE fraud_cases SET status=$2, reviewed_by=$3, resolved_at=CASE WHEN $2 IN ('confirmed','false_positive','resolved') THEN NOW() ELSE NULL END WHERE id=$1`, id, string(status), reviewedBy)
	return err
}

func (r *TrustRepository) GetDashboard(ctx context.Context) (*domain.TrustDashboard, error) {
	var d domain.TrustDashboard
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores`).Scan(&d.TotalUsers)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(overall_score),0) FROM reputation_scores`).Scan(&d.AvgReputationScore)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE trust_level='low'`).Scan(&d.LowTrustUsers)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE trust_level='medium'`).Scan(&d.MediumTrustUsers)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE trust_level='high'`).Scan(&d.HighTrustUsers)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE trust_level='verified'`).Scan(&d.VerifiedUsers)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM fraud_cases WHERE status IN ('open','investigating')`).Scan(&d.OpenFraudCases)
	return &d, nil
}

func (r *TrustRepository) GetScoreDistribution(ctx context.Context) (*domain.ScoreDistribution, error) {
	var sd domain.ScoreDistribution
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE overall_score BETWEEN 0 AND 200`).Scan(&sd.Range0_200)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE overall_score BETWEEN 201 AND 400`).Scan(&sd.Range201_400)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE overall_score BETWEEN 401 AND 600`).Scan(&sd.Range401_600)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE overall_score BETWEEN 601 AND 800`).Scan(&sd.Range601_800)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM reputation_scores WHERE overall_score BETWEEN 801 AND 1000`).Scan(&sd.Range801_1000)
	return &sd, nil
}

func (r *TrustRepository) GetFlaggedUsers(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT user_id FROM reputation_scores WHERE trust_level='low' ORDER BY overall_score ASC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func (r *TrustRepository) GetReputationHistory(ctx context.Context, userID uuid.UUID) ([]domain.ReputationScore, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT user_id, overall_score, trust_level, positive_signal_weight, negative_signal_weight, score_decay_factor, model_version, calculated_at, next_recalculation_at
		FROM reputation_score_history WHERE user_id=$1 ORDER BY calculated_at DESC LIMIT 30
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.ReputationScore
	for rows.Next() {
		var rs domain.ReputationScore
		if err := rows.Scan(&rs.UserID, &rs.OverallScore, &rs.TrustLevel, &rs.PositiveSignalWeight, &rs.NegativeSignalWeight, &rs.ScoreDecayFactor, &rs.ModelVersion, &rs.CalculatedAt, &rs.NextRecalculationAt); err != nil {
			return nil, err
		}
		res = append(res, rs)
	}
	return res, nil
}

func (r *TrustRepository) SaveReputationHistory(ctx context.Context, rs *domain.ReputationScore) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO reputation_score_history (user_id, overall_score, trust_level, positive_signal_weight, negative_signal_weight, score_decay_factor, model_version, calculated_at, next_recalculation_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, rs.UserID, rs.OverallScore, string(rs.TrustLevel), rs.PositiveSignalWeight, rs.NegativeSignalWeight, rs.ScoreDecayFactor, rs.ModelVersion, rs.CalculatedAt, rs.NextRecalculationAt)
	return err
}

func (r *TrustRepository) GetUserIDByRiskAssessment(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT user_id FROM risk_assessments WHERE id=$1`, id).Scan(&userID)
	return userID, err
}

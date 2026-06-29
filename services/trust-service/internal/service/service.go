package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/trust-service/internal/domain"
	"github.com/elevatecompact/spark/services/trust-service/internal/repository"
)

type TrustService struct {
	repo *repository.TrustRepository
	evt  domain.EventProducer
}

func NewTrustService(repo *repository.TrustRepository, evt domain.EventProducer) *TrustService {
	return &TrustService{repo: repo, evt: evt}
}

func (s *TrustService) GetReputation(ctx context.Context, userID uuid.UUID) (*domain.ReputationScore, error) {
	rs, err := s.repo.GetReputation(ctx, userID)
	if err != nil {
		rs = s.initializeReputation(userID)
		rs.NextRecalculationAt = time.Now().Add(24 * time.Hour)
		if err := s.repo.UpsertReputation(ctx, rs); err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func (s *TrustService) GetReputationHistory(ctx context.Context, userID uuid.UUID) ([]domain.ReputationScore, error) {
	return s.repo.GetReputationHistory(ctx, userID)
}

func (s *TrustService) RecalculateReputation(ctx context.Context, userID uuid.UUID) (*domain.ReputationScore, error) {
	signals, err := s.repo.GetSignals(ctx, userID)
	if err != nil {
		return nil, err
	}
	var posWeight, negWeight int
	for _, sig := range signals {
		if sig.SignalType == domain.SignalPositive {
			posWeight += sig.Weight
		} else {
			negWeight += sig.Weight
		}
	}
	score := posWeight - negWeight
	if score < 0 {
		score = 0
	}
	if score > 1000 {
		score = 1000
	}
	trustLevel := s.determineTrustLevel(score)
	rs := &domain.ReputationScore{
		UserID:               userID,
		OverallScore:         score,
		TrustLevel:           trustLevel,
		PositiveSignalWeight: posWeight,
		NegativeSignalWeight: negWeight,
		ScoreDecayFactor:     0.95,
		ModelVersion:         "v1",
		CalculatedAt:         time.Now(),
		NextRecalculationAt:  time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.SaveReputationHistory(ctx, rs); err != nil {
		log.Error().Err(err).Msg("failed to save reputation history")
	}
	if err := s.repo.UpsertReputation(ctx, rs); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "trust.reputation.updated", map[string]interface{}{
		"userId": userID, "overallScore": score, "trustLevel": trustLevel,
	})
	return rs, nil
}

func (s *TrustService) determineTrustLevel(score int) domain.TrustLevel {
	switch {
	case score >= 900:
		return domain.TrustVerified
	case score >= 700:
		return domain.TrustHigh
	case score >= 500:
		return domain.TrustMedium
	default:
		return domain.TrustLow
	}
}

func (s *TrustService) initializeReputation(userID uuid.UUID) *domain.ReputationScore {
	return &domain.ReputationScore{
		UserID:               userID,
		OverallScore:         500,
		TrustLevel:           domain.TrustMedium,
		PositiveSignalWeight: 0,
		NegativeSignalWeight: 0,
		ScoreDecayFactor:     0.95,
		ModelVersion:         "v1",
		CalculatedAt:         time.Now(),
	}
}

func (s *TrustService) GetTrustSignals(ctx context.Context, userID uuid.UUID) ([]domain.TrustSignal, error) {
	return s.repo.GetSignals(ctx, userID)
}

func (s *TrustService) RecordSignal(ctx context.Context, signal *domain.TrustSignal) (*domain.TrustSignal, error) {
	signal.ID = uuid.New()
	signal.RecordedAt = time.Now()
	if err := s.repo.RecordSignal(ctx, signal); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "trust.signal.recorded", map[string]interface{}{
		"userId": signal.UserID, "signalType": signal.SignalType, "category": signal.Category,
	})
	return signal, nil
}

func (s *TrustService) GetTrustLevel(ctx context.Context, userID uuid.UUID) (*domain.ReputationScore, error) {
	return s.GetReputation(ctx, userID)
}

func (s *TrustService) AssessRisk(ctx context.Context, userID uuid.UUID, actionType string, contextData map[string]interface{}) (*domain.RiskAssessment, error) {
	ctxJSON, _ := json.Marshal(contextData)
	score := rand.Float64()
	riskLevel := s.determineRiskLevel(score)
	action := s.determineAction(riskLevel)
	ra := &domain.RiskAssessment{
		ID:                uuid.New(),
		UserID:            userID,
		ActionType:        actionType,
		Context:           ctxJSON,
		RiskScore:         score,
		RiskLevel:         riskLevel,
		TriggeredRules:    []string{},
		RecommendedAction: action,
		AssessedAt:        time.Now(),
	}
	if err := s.repo.CreateRiskAssessment(ctx, ra); err != nil {
		return nil, err
	}
	if riskLevel == domain.RiskHigh || riskLevel == domain.RiskCritical {
		s.evt.Publish(ctx, "trust.risk.alert", map[string]interface{}{
			"assessmentId": ra.ID, "userId": userID, "riskScore": score, "riskLevel": riskLevel,
		})
	}
	return ra, nil
}

func (s *TrustService) GetRiskAssessment(ctx context.Context, id uuid.UUID) (*domain.RiskAssessment, error) {
	return s.repo.GetRiskAssessment(ctx, id)
}

func (s *TrustService) CreateRiskRule(ctx context.Context, rr *domain.RiskRule) (*domain.RiskRule, error) {
	rr.ID = uuid.New()
	rr.CreatedAt = time.Now()
	if !rr.IsActive {
		rr.IsActive = true
	}
	if err := s.repo.CreateRiskRule(ctx, rr); err != nil {
		return nil, err
	}
	return rr, nil
}

func (s *TrustService) UpdateRiskRule(ctx context.Context, rr *domain.RiskRule) error {
	return s.repo.UpdateRiskRule(ctx, rr)
}

func (s *TrustService) CheckPaymentFraud(ctx context.Context, contextData map[string]interface{}) (*domain.FraudDetectionResult, error) {
	score := rand.Float64()
	isFraud := score > 0.8
	result := &domain.FraudDetectionResult{
		IsFraudulent: isFraud,
		Score:        score,
		Reasons:      []string{},
	}
	if isFraud {
		result.Reasons = append(result.Reasons, "suspicious payment pattern detected")
		caseID := uuid.New()
		evidence, _ := json.Marshal(contextData)
		fc := &domain.FraudCase{
			ID:               caseID,
			CaseType:         domain.FraudPaymentFraud,
			Status:           domain.FraudOpen,
			Evidence:         evidence,
			AutomatedDecision: "flag_for_review",
			CreatedAt:        time.Now(),
		}
		if uid, ok := contextData["userId"].(uuid.UUID); ok {
			fc.UserID = uid
		}
		s.repo.CreateFraudCase(ctx, fc)
		s.evt.Publish(ctx, "trust.fraud.case.opened", map[string]interface{}{
			"caseId": caseID, "caseType": fc.CaseType,
		})
		result.CaseID = &caseID
	}
	return result, nil
}

func (s *TrustService) CheckAccountFraud(ctx context.Context, contextData map[string]interface{}) (*domain.FraudDetectionResult, error) {
	return s.CheckPaymentFraud(ctx, contextData)
}

func (s *TrustService) ReportFraud(ctx context.Context, userID uuid.UUID, reason string, evidence map[string]interface{}) (*domain.FraudCase, error) {
	evidenceJSON, _ := json.Marshal(evidence)
	fc := &domain.FraudCase{
		ID:               uuid.New(),
		UserID:           userID,
		CaseType:         s.classifyFraudType(reason),
		Status:           domain.FraudOpen,
		Evidence:         evidenceJSON,
		AutomatedDecision: "reported",
		CreatedAt:        time.Now(),
	}
	if err := s.repo.CreateFraudCase(ctx, fc); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "trust.fraud.case.opened", map[string]interface{}{
		"caseId": fc.ID, "userId": userID, "reason": reason,
	})
	return fc, nil
}

func (s *TrustService) ListFraudCases(ctx context.Context, status *domain.FraudCaseStatus, limit, offset int) ([]domain.FraudCase, error) {
	return s.repo.ListFraudCases(ctx, status, limit, offset)
}

func (s *TrustService) ResolveFraudCase(ctx context.Context, caseID uuid.UUID, status domain.FraudCaseStatus, reviewerID uuid.UUID) error {
	if err := s.repo.UpdateFraudCaseStatus(ctx, caseID, status, &reviewerID); err != nil {
		return err
	}
	s.evt.Publish(ctx, "trust.fraud.case.resolved", map[string]interface{}{
		"caseId": caseID, "status": status,
	})
	return nil
}

func (s *TrustService) GetDashboard(ctx context.Context) (*domain.TrustDashboard, error) {
	return s.repo.GetDashboard(ctx)
}

func (s *TrustService) GetScoreDistribution(ctx context.Context) (*domain.ScoreDistribution, error) {
	return s.repo.GetScoreDistribution(ctx)
}

func (s *TrustService) GetFlaggedUsers(ctx context.Context, limit, offset int) ([]uuid.UUID, error) {
	return s.repo.GetFlaggedUsers(ctx, limit, offset)
}

func (s *TrustService) UpdateThresholds(ctx context.Context, thresholds map[string]int) error {
	log.Info().Interface("thresholds", thresholds).Msg("thresholds updated (noop)")
	return nil
}

func (s *TrustService) determineRiskLevel(score float64) domain.RiskLevel {
	switch {
	case score >= 0.8:
		return domain.RiskCritical
	case score >= 0.6:
		return domain.RiskHigh
	case score >= 0.3:
		return domain.RiskMedium
	default:
		return domain.RiskLow
	}
}

func (s *TrustService) determineAction(level domain.RiskLevel) domain.RecommendedAction {
	switch level {
	case domain.RiskCritical:
		return domain.ActionBlock
	case domain.RiskHigh:
		return domain.ActionChallenge
	case domain.RiskMedium:
		return domain.ActionReview
	default:
		return domain.ActionAllow
	}
}

func (s *TrustService) classifyFraudType(reason string) domain.FraudCaseType {
	switch reason {
	case "payment_fraud":
		return domain.FraudPaymentFraud
	case "account_takeover":
		return domain.FraudAccountTakeover
	case "coordinated_behavior":
		return domain.FraudCoordinatedBehavior
	case "spam":
		return domain.FraudSpam
	default:
		return domain.FraudSpam
	}
}

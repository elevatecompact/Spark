package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/elevatecompact/spark/services/competition-service/internal/domain"
	"github.com/elevatecompact/spark/services/competition-service/internal/events"
	"github.com/elevatecompact/spark/services/competition-service/internal/repository"
)

type CompetitionService interface {
	Create(ctx context.Context, c *domain.Competition) (*domain.Competition, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Competition, error)
	Update(ctx context.Context, c *domain.Competition) error
	List(ctx context.Context, status domain.CompetitionStatus, page, size int) ([]domain.Competition, error)
	Start(ctx context.Context, id uuid.UUID) error
	End(ctx context.Context, id uuid.UUID) error

	Register(ctx context.Context, competitionID, userID uuid.UUID) error
	Withdraw(ctx context.Context, competitionID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, competitionID uuid.UUID) ([]domain.CompetitionParticipant, error)

	GetBracket(ctx context.Context, competitionID uuid.UUID) ([]domain.Match, error)
	SubmitScore(ctx context.Context, matchID uuid.UUID, winnerID uuid.UUID, scores map[string]interface{}) error
	ConfirmMatch(ctx context.Context, matchID uuid.UUID) error
	DisputeMatch(ctx context.Context, matchID uuid.UUID) error

	AssignJudge(ctx context.Context, competitionID, judgeID uuid.UUID) error
	ScoreSubmission(ctx context.Context, submissionID uuid.UUID, judgeID uuid.UUID, score float64) error
	ListSubmissions(ctx context.Context, competitionID uuid.UUID) ([]domain.JudgingSubmission, error)

	GetLeaderboard(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error)
	GetResults(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error)
	GetPrizes(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error)
	DistributePrizes(ctx context.Context, competitionID uuid.UUID) error

	Cancel(ctx context.Context, id uuid.UUID) error
	OverrideMatch(ctx context.Context, matchID uuid.UUID, winnerID uuid.UUID) error
	GetAdminStats(ctx context.Context) (*domain.CompetitionAdminStats, error)
}

type compService struct {
	repo     repository.CompetitionRepository
	eventPub events.EventProducer
}

func NewCompetitionService(repo repository.CompetitionRepository, eventPub events.EventProducer) CompetitionService {
	return &compService{repo: repo, eventPub: eventPub}
}

func (s *compService) Create(ctx context.Context, c *domain.Competition) (*domain.Competition, error) {
	c.ID = uuid.New()
	c.Status = domain.CompDraft
	c.CreatedAt = time.Now().UTC()
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.eventPub.PublishCompetitionCreated(ctx, c.ID)
	return c, nil
}

func (s *compService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Competition, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *compService) Update(ctx context.Context, c *domain.Competition) error {
	return s.repo.Update(ctx, c)
}

func (s *compService) List(ctx context.Context, status domain.CompetitionStatus, page, size int) ([]domain.Competition, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 25
	}
	offset := (page - 1) * size
	return s.repo.List(ctx, status, offset, size)
}

func (s *compService) Start(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.Status = domain.CompStarted
	now := time.Now().UTC()
	c.CompetitionStart = &now
	if err := s.repo.Update(ctx, c); err != nil {
		return err
	}
	return s.eventPub.PublishCompetitionStarted(ctx, id)
}

func (s *compService) End(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.Status = domain.CompCompleted
	now := time.Now().UTC()
	c.CompetitionEnd = &now
	if err := s.repo.Update(ctx, c); err != nil {
		return err
	}
	return s.eventPub.PublishCompetitionEnded(ctx, id)
}

func (s *compService) Register(ctx context.Context, competitionID, userID uuid.UUID) error {
	count, err := s.repo.GetParticipantCount(ctx, competitionID)
	if err != nil {
		return err
	}
	comp, err := s.repo.GetByID(ctx, competitionID)
	if err != nil {
		return err
	}
	if count >= comp.MaxParticipants {
		return domain.ErrCompetitionFull
	}
	participant := &domain.CompetitionParticipant{
		ID:            uuid.New(),
		CompetitionID: competitionID,
		UserID:        userID,
		SeedNumber:    count + 1,
		Status:        domain.PartRegistered,
	}
	if err := s.repo.RegisterParticipant(ctx, participant); err != nil {
		return err
	}
	return s.eventPub.PublishParticipantRegistered(ctx, competitionID, userID)
}

func (s *compService) Withdraw(ctx context.Context, competitionID, userID uuid.UUID) error {
	return s.repo.WithdrawParticipant(ctx, competitionID, userID)
}

func (s *compService) ListParticipants(ctx context.Context, competitionID uuid.UUID) ([]domain.CompetitionParticipant, error) {
	return s.repo.ListParticipants(ctx, competitionID)
}

func (s *compService) GetBracket(ctx context.Context, competitionID uuid.UUID) ([]domain.Match, error) {
	return s.repo.ListMatches(ctx, competitionID)
}

func (s *compService) SubmitScore(ctx context.Context, matchID uuid.UUID, winnerID uuid.UUID, scores map[string]interface{}) error {
	m, err := s.repo.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	m.WinnerID = &winnerID
	m.Scores = scores
	m.Status = domain.MatchCompleted
	now := time.Now().UTC()
	m.CompletedAt = &now
	if err := s.repo.UpdateMatch(ctx, m); err != nil {
		return err
	}
	s.eventPub.PublishMatchCompleted(ctx, &events.MatchCompletedEvent{
		CompetitionID:  m.CompetitionID,
		MatchID:        matchID,
		WinnerID:       winnerID,
		Round:          m.Round,
		BracketPosition: m.BracketPosition,
		CompletedAt:    now,
	})
	return nil
}

func (s *compService) ConfirmMatch(ctx context.Context, matchID uuid.UUID) error {
	return nil
}

func (s *compService) DisputeMatch(ctx context.Context, matchID uuid.UUID) error {
	m, err := s.repo.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	m.Status = domain.MatchDisputed
	if err := s.repo.UpdateMatch(ctx, m); err != nil {
		return err
	}
	return s.eventPub.PublishMatchDisputed(ctx, matchID)
}

func (s *compService) AssignJudge(ctx context.Context, competitionID, judgeID uuid.UUID) error {
	return nil
}

func (s *compService) ScoreSubmission(ctx context.Context, submissionID uuid.UUID, judgeID uuid.UUID, score float64) error {
	sub, err := s.repo.GetSubmission(ctx, submissionID)
	if err != nil {
		return err
	}
	if sub.Scores == nil {
		sub.Scores = make(map[string]float64)
	}
	sub.Scores[judgeID.String()] = score
	total := 0.0
	for _, v := range sub.Scores {
		total += v
	}
	sub.AverageScore = total / float64(len(sub.Scores))
	sub.Status = domain.SubScored
	return s.repo.UpdateSubmission(ctx, sub)
}

func (s *compService) ListSubmissions(ctx context.Context, competitionID uuid.UUID) ([]domain.JudgingSubmission, error) {
	return s.repo.ListSubmissions(ctx, competitionID)
}

func (s *compService) GetLeaderboard(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error) {
	submissions, err := s.repo.ListSubmissions(ctx, competitionID)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.LeaderboardEntry, len(submissions))
	for i, sub := range submissions {
		entries[i] = domain.LeaderboardEntry{
			ParticipantID: sub.ParticipantID,
			Score:         sub.AverageScore,
			Rank:          i + 1,
		}
	}
	if entries == nil {
		entries = []domain.LeaderboardEntry{}
	}
	return entries, nil
}

func (s *compService) GetResults(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error) {
	return s.GetLeaderboard(ctx, competitionID)
}

func (s *compService) GetPrizes(ctx context.Context, competitionID uuid.UUID) ([]domain.LeaderboardEntry, error) {
	return s.GetLeaderboard(ctx, competitionID)
}

func (s *compService) DistributePrizes(ctx context.Context, competitionID uuid.UUID) error {
	return s.eventPub.PublishPrizeDistributed(ctx, competitionID)
}

func (s *compService) Cancel(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	c.Status = domain.CompCancelled
	return s.repo.Update(ctx, c)
}

func (s *compService) OverrideMatch(ctx context.Context, matchID uuid.UUID, winnerID uuid.UUID) error {
	return s.SubmitScore(ctx, matchID, winnerID, map[string]interface{}{"override": true})
}

func (s *compService) GetAdminStats(ctx context.Context) (*domain.CompetitionAdminStats, error) {
	return s.repo.GetAdminStats(ctx)
}

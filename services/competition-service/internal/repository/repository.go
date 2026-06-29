package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/competition-service/internal/domain"
)

type CompetitionRepository interface {
	Create(ctx context.Context, c *domain.Competition) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Competition, error)
	Update(ctx context.Context, c *domain.Competition) error
	List(ctx context.Context, status domain.CompetitionStatus, offset, limit int) ([]domain.Competition, error)

	RegisterParticipant(ctx context.Context, p *domain.CompetitionParticipant) error
	WithdrawParticipant(ctx context.Context, competitionID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, competitionID uuid.UUID) ([]domain.CompetitionParticipant, error)
	GetParticipantCount(ctx context.Context, competitionID uuid.UUID) (int, error)

	CreateMatch(ctx context.Context, m *domain.Match) error
	GetMatch(ctx context.Context, id uuid.UUID) (*domain.Match, error)
	UpdateMatch(ctx context.Context, m *domain.Match) error
	ListMatches(ctx context.Context, competitionID uuid.UUID) ([]domain.Match, error)

	CreateSubmission(ctx context.Context, s *domain.JudgingSubmission) error
	GetSubmission(ctx context.Context, id uuid.UUID) (*domain.JudgingSubmission, error)
	UpdateSubmission(ctx context.Context, s *domain.JudgingSubmission) error
	ListSubmissions(ctx context.Context, competitionID uuid.UUID) ([]domain.JudgingSubmission, error)

	GetAdminStats(ctx context.Context) (*domain.CompetitionAdminStats, error)
}

type compRepo struct{ pool *pgxpool.Pool }

func NewCompetitionRepository(pool *pgxpool.Pool) CompetitionRepository { return &compRepo{pool} }

func (r *compRepo) Create(ctx context.Context, c *domain.Competition) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO competitions (id, creator_id, title, description, format, bracket_type, status, max_participants, registration_start, registration_end, competition_start, competition_end, prize_pool_cents, rules, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,NOW())`,
		c.ID, c.CreatorID, c.Title, c.Description, c.Format, c.BracketType, c.Status, c.MaxParticipants, c.RegistrationStart, c.RegistrationEnd, c.CompetitionStart, c.CompetitionEnd, c.PrizePoolCents, c.Rules)
	return err
}

func (r *compRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Competition, error) {
	c := &domain.Competition{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, creator_id, title, description, format, bracket_type, status, max_participants, registration_start, registration_end, competition_start, competition_end, prize_pool_cents, rules, created_at FROM competitions WHERE id=$1`, id).
		Scan(&c.ID, &c.CreatorID, &c.Title, &c.Description, &c.Format, &c.BracketType, &c.Status, &c.MaxParticipants, &c.RegistrationStart, &c.RegistrationEnd, &c.CompetitionStart, &c.CompetitionEnd, &c.PrizePoolCents, &c.Rules, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return c, err
}

func (r *compRepo) Update(ctx context.Context, c *domain.Competition) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE competitions SET title=$2, description=$3, bracket_type=$4, status=$5, competition_start=$6, competition_end=$7, prize_pool_cents=$8, rules=$9 WHERE id=$1`,
		c.ID, c.Title, c.Description, c.BracketType, c.Status, c.CompetitionStart, c.CompetitionEnd, c.PrizePoolCents, c.Rules)
	return err
}

func (r *compRepo) List(ctx context.Context, status domain.CompetitionStatus, offset, limit int) ([]domain.Competition, error) {
	if limit <= 0 || limit > 50 {
		limit = 25
	}
	var rows pgx.Rows
	var err error
	if status != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, format, bracket_type, status, max_participants, registration_start, registration_end, competition_start, competition_end, prize_pool_cents, rules, created_at FROM competitions WHERE status=$1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`,
			status, offset, limit)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, creator_id, title, description, format, bracket_type, status, max_participants, registration_start, registration_end, competition_start, competition_end, prize_pool_cents, rules, created_at FROM competitions ORDER BY created_at DESC OFFSET $1 LIMIT $2`,
			offset, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comps []domain.Competition
	for rows.Next() {
		var c domain.Competition
		if err := rows.Scan(&c.ID, &c.CreatorID, &c.Title, &c.Description, &c.Format, &c.BracketType, &c.Status, &c.MaxParticipants, &c.RegistrationStart, &c.RegistrationEnd, &c.CompetitionStart, &c.CompetitionEnd, &c.PrizePoolCents, &c.Rules, &c.CreatedAt); err != nil {
			return nil, err
		}
		comps = append(comps, c)
	}
	if comps == nil {
		comps = []domain.Competition{}
	}
	return comps, nil
}

func (r *compRepo) RegisterParticipant(ctx context.Context, p *domain.CompetitionParticipant) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO competition_participants (id, competition_id, user_id, seed_number, status, bracket_position, registered_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		p.ID, p.CompetitionID, p.UserID, p.SeedNumber, p.Status, p.BracketPosition)
	return err
}

func (r *compRepo) WithdrawParticipant(ctx context.Context, competitionID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE competition_participants SET status='withdrawn' WHERE competition_id=$1 AND user_id=$2`,
		competitionID, userID)
	return err
}

func (r *compRepo) ListParticipants(ctx context.Context, competitionID uuid.UUID) ([]domain.CompetitionParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, competition_id, user_id, seed_number, status, bracket_position, registered_at FROM competition_participants WHERE competition_id=$1 ORDER BY seed_number ASC`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parts []domain.CompetitionParticipant
	for rows.Next() {
		var p domain.CompetitionParticipant
		if err := rows.Scan(&p.ID, &p.CompetitionID, &p.UserID, &p.SeedNumber, &p.Status, &p.BracketPosition, &p.RegisteredAt); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	if parts == nil {
		parts = []domain.CompetitionParticipant{}
	}
	return parts, nil
}

func (r *compRepo) GetParticipantCount(ctx context.Context, competitionID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM competition_participants WHERE competition_id=$1 AND status NOT IN ('withdrawn')`, competitionID).Scan(&count)
	return count, err
}

func (r *compRepo) CreateMatch(ctx context.Context, m *domain.Match) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO matches (id, competition_id, round, bracket_position, participant1_id, participant2_id, winner_id, scores, status, scheduled_at, completed_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		m.ID, m.CompetitionID, m.Round, m.BracketPosition, m.Participant1ID, m.Participant2ID, m.WinnerID, m.Scores, m.Status, m.ScheduledAt, m.CompletedAt)
	return err
}

func (r *compRepo) GetMatch(ctx context.Context, id uuid.UUID) (*domain.Match, error) {
	m := &domain.Match{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, competition_id, round, bracket_position, participant1_id, participant2_id, winner_id, scores, status, scheduled_at, completed_at FROM matches WHERE id=$1`, id).
		Scan(&m.ID, &m.CompetitionID, &m.Round, &m.BracketPosition, &m.Participant1ID, &m.Participant2ID, &m.WinnerID, &m.Scores, &m.Status, &m.ScheduledAt, &m.CompletedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return m, err
}

func (r *compRepo) UpdateMatch(ctx context.Context, m *domain.Match) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE matches SET participant1_id=$2, participant2_id=$3, winner_id=$4, scores=$5, status=$6, completed_at=$7 WHERE id=$1`,
		m.ID, m.Participant1ID, m.Participant2ID, m.WinnerID, m.Scores, m.Status, m.CompletedAt)
	return err
}

func (r *compRepo) ListMatches(ctx context.Context, competitionID uuid.UUID) ([]domain.Match, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, competition_id, round, bracket_position, participant1_id, participant2_id, winner_id, scores, status, scheduled_at, completed_at FROM matches WHERE competition_id=$1 ORDER BY round ASC, bracket_position ASC`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		if err := rows.Scan(&m.ID, &m.CompetitionID, &m.Round, &m.BracketPosition, &m.Participant1ID, &m.Participant2ID, &m.WinnerID, &m.Scores, &m.Status, &m.ScheduledAt, &m.CompletedAt); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	if matches == nil {
		matches = []domain.Match{}
	}
	return matches, nil
}

func (r *compRepo) CreateSubmission(ctx context.Context, s *domain.JudgingSubmission) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO judging_submissions (id, competition_id, participant_id, scores, average_score, status) VALUES ($1,$2,$3,$4,$5,$6)`,
		s.ID, s.CompetitionID, s.ParticipantID, s.Scores, s.AverageScore, s.Status)
	return err
}

func (r *compRepo) GetSubmission(ctx context.Context, id uuid.UUID) (*domain.JudgingSubmission, error) {
	s := &domain.JudgingSubmission{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, competition_id, participant_id, scores, average_score, status FROM judging_submissions WHERE id=$1`, id).
		Scan(&s.ID, &s.CompetitionID, &s.ParticipantID, &s.Scores, &s.AverageScore, &s.Status)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return s, err
}

func (r *compRepo) UpdateSubmission(ctx context.Context, s *domain.JudgingSubmission) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE judging_submissions SET scores=$2, average_score=$3, status=$4 WHERE id=$1`,
		s.ID, s.Scores, s.AverageScore, s.Status)
	return err
}

func (r *compRepo) ListSubmissions(ctx context.Context, competitionID uuid.UUID) ([]domain.JudgingSubmission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, competition_id, participant_id, scores, average_score, status FROM judging_submissions WHERE competition_id=$1 ORDER BY average_score DESC`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var submissions []domain.JudgingSubmission
	for rows.Next() {
		var s domain.JudgingSubmission
		if err := rows.Scan(&s.ID, &s.CompetitionID, &s.ParticipantID, &s.Scores, &s.AverageScore, &s.Status); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	if submissions == nil {
		submissions = []domain.JudgingSubmission{}
	}
	return submissions, nil
}

func (r *compRepo) GetAdminStats(ctx context.Context) (*domain.CompetitionAdminStats, error) {
	s := &domain.CompetitionAdminStats{}
	err := r.pool.QueryRow(ctx,
		`SELECT
		 (SELECT COUNT(*) FROM competitions) as total,
		 (SELECT COALESCE(SUM(c.max_participants),0) FROM competitions WHERE status='open' OR status='started') as total_parts,
		 (SELECT COUNT(*) FROM competitions WHERE status='started') as active,
		 (SELECT COALESCE(SUM(prize_pool_cents),0) FROM competitions) as prizes`).
		Scan(&s.TotalCompetitions, &s.TotalParticipants, &s.ActiveCompetitions, &s.PrizePoolTotal)
	if err != nil {
		return nil, err
	}
	return s, nil
}

CREATE TABLE IF NOT EXISTS competitions (
    id UUID PRIMARY KEY,
    creator_id UUID NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    format VARCHAR(16) NOT NULL,
    bracket_type VARCHAR(8),
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    max_participants INT NOT NULL DEFAULT 128,
    registration_start TIMESTAMPTZ,
    registration_end TIMESTAMPTZ,
    competition_start TIMESTAMPTZ,
    competition_end TIMESTAMPTZ,
    prize_pool_cents BIGINT NOT NULL DEFAULT 0,
    rules TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS competition_participants (
    id UUID PRIMARY KEY,
    competition_id UUID NOT NULL REFERENCES competitions(id),
    user_id UUID NOT NULL,
    seed_number INT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'registered',
    bracket_position INT,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comp_parts ON competition_participants (competition_id);

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY,
    competition_id UUID NOT NULL REFERENCES competitions(id),
    round INT NOT NULL,
    bracket_position INT NOT NULL,
    participant1_id UUID,
    participant2_id UUID,
    winner_id UUID,
    scores JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(16) NOT NULL DEFAULT 'scheduled',
    scheduled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_matches_comp ON matches (competition_id, round);

CREATE TABLE IF NOT EXISTS judging_submissions (
    id UUID PRIMARY KEY,
    competition_id UUID NOT NULL REFERENCES competitions(id),
    participant_id UUID NOT NULL REFERENCES competition_participants(id),
    scores JSONB NOT NULL DEFAULT '{}'::jsonb,
    average_score FLOAT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'pending'
);

CREATE INDEX IF NOT EXISTS idx_judging_comp ON judging_submissions (competition_id);

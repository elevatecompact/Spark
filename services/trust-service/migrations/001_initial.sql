CREATE TABLE IF NOT EXISTS reputation_scores (
    user_id UUID PRIMARY KEY,
    overall_score INT NOT NULL DEFAULT 500,
    trust_level VARCHAR(16) NOT NULL DEFAULT 'medium',
    positive_signal_weight INT NOT NULL DEFAULT 0,
    negative_signal_weight INT NOT NULL DEFAULT 0,
    score_decay_factor DOUBLE PRECISION NOT NULL DEFAULT 0.95,
    model_version VARCHAR(16) NOT NULL DEFAULT 'v1',
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    next_recalculation_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS reputation_score_history (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    overall_score INT NOT NULL,
    trust_level VARCHAR(16) NOT NULL,
    positive_signal_weight INT NOT NULL,
    negative_signal_weight INT NOT NULL,
    score_decay_factor DOUBLE PRECISION NOT NULL,
    model_version VARCHAR(16) NOT NULL,
    calculated_at TIMESTAMPTZ NOT NULL,
    next_recalculation_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_reputation_history_user ON reputation_score_history (user_id, calculated_at DESC);

CREATE TABLE IF NOT EXISTS trust_signals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    signal_type VARCHAR(16) NOT NULL,
    category VARCHAR(16) NOT NULL,
    weight INT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    source_entity_type VARCHAR(64) NOT NULL DEFAULT '',
    source_entity_id VARCHAR(64) NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trust_signals_user ON trust_signals (user_id);

CREATE TABLE IF NOT EXISTS risk_assessments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    action_type VARCHAR(64) NOT NULL,
    context JSONB DEFAULT '{}'::jsonb,
    risk_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    risk_level VARCHAR(16) NOT NULL,
    triggered_rules TEXT[] DEFAULT '{}',
    recommended_action VARCHAR(16) NOT NULL,
    assessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fraud_cases (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    case_type VARCHAR(32) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'open',
    evidence JSONB DEFAULT '{}'::jsonb,
    automated_decision VARCHAR(32) NOT NULL DEFAULT '',
    reviewed_by UUID,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS risk_rules (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    category VARCHAR(64) NOT NULL,
    conditions JSONB DEFAULT '{}'::jsonb,
    risk_score_impact DOUBLE PRECISION NOT NULL DEFAULT 0,
    action VARCHAR(32) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

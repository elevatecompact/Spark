CREATE TABLE IF NOT EXISTS moderation_rules (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    category VARCHAR(32) NOT NULL CHECK (category IN ('harassment','spam','nsfw','violence','hate_speech')),
    severity VARCHAR(16) NOT NULL CHECK (severity IN ('warn','restrict','remove','suspend')),
    conditions JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS moderation_actions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    content_id UUID,
    rule_id UUID REFERENCES moderation_rules(id),
    action_type VARCHAR(16) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    applied_by VARCHAR(64) NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    duration INT,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mod_actions_user ON moderation_actions (user_id);
CREATE INDEX IF NOT EXISTS idx_mod_actions_status ON moderation_actions (status);

CREATE TABLE IF NOT EXISTS review_queue (
    id UUID PRIMARY KEY,
    content_type VARCHAR(32) NOT NULL,
    content_id UUID NOT NULL,
    flagged_by VARCHAR(32) NOT NULL DEFAULT 'automated',
    reasons TEXT[] NOT NULL DEFAULT '{}',
    assigned_moderator UUID,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    resolution TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_review_queue_status ON review_queue (status);

CREATE TABLE IF NOT EXISTS content_reports (
    id UUID PRIMARY KEY,
    reporter_id UUID NOT NULL,
    content_type VARCHAR(32) NOT NULL,
    content_id UUID NOT NULL,
    reason VARCHAR(64) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_content_reports_status ON content_reports (status);

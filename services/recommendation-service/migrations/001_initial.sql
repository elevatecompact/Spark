CREATE TABLE IF NOT EXISTS user_embeddings (
    user_id UUID PRIMARY KEY,
    embedding JSONB NOT NULL DEFAULT '[]'::jsonb,
    model_version VARCHAR(64) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_embeddings (
    content_id UUID PRIMARY KEY,
    embedding JSONB NOT NULL DEFAULT '[]'::jsonb,
    model_version VARCHAR(64) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_content_interactions (
    user_id UUID NOT NULL,
    content_id UUID NOT NULL,
    interaction_type VARCHAR(32) NOT NULL,
    weight FLOAT NOT NULL DEFAULT 1.0,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_interactions_user_ts ON user_content_interactions (user_id, timestamp DESC);

CREATE TABLE IF NOT EXISTS model_versions (
    version VARCHAR(64) PRIMARY KEY,
    deployed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metrics TEXT,
    is_active BOOLEAN NOT NULL DEFAULT false
);

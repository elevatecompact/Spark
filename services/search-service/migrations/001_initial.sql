CREATE TABLE IF NOT EXISTS search_analytics (
    id BIGSERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    result_ids TEXT[] DEFAULT '{}',
    latency_ms BIGINT NOT NULL DEFAULT 0,
    user_id UUID,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_search_analytics_ts ON search_analytics (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_search_analytics_query ON search_analytics (query);

CREATE TABLE IF NOT EXISTS synonym_dictionary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    term TEXT NOT NULL UNIQUE,
    synonyms TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

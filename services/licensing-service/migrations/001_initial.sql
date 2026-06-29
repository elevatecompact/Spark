CREATE TABLE IF NOT EXISTS licenses (
    id UUID PRIMARY KEY,
    rights_holder_id UUID NOT NULL,
    licensee_id UUID NOT NULL,
    content_id UUID,
    type VARCHAR(16) NOT NULL DEFAULT 'non_exclusive',
    scope VARCHAR(16) NOT NULL DEFAULT 'global',
    territory JSONB DEFAULT '[]'::jsonb,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    rate_type VARCHAR(16) NOT NULL DEFAULT 'flat',
    rate_cents BIGINT NOT NULL DEFAULT 0,
    revenue_share_percent DECIMAL(5,2) NOT NULL DEFAULT 0,
    min_guarantee_cents BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    terms_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_rights (
    id UUID PRIMARY KEY,
    content_id UUID NOT NULL UNIQUE,
    rights_holder_id UUID NOT NULL,
    license_id UUID NOT NULL REFERENCES licenses(id),
    restrictions JSONB DEFAULT '{}'::jsonb,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS usage_log (
    id UUID PRIMARY KEY,
    license_id UUID NOT NULL REFERENCES licenses(id),
    content_id UUID NOT NULL,
    usage_type VARCHAR(16) NOT NULL,
    context JSONB DEFAULT '{}'::jsonb,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_usage_license ON usage_log (license_id, recorded_at);
CREATE INDEX IF NOT EXISTS idx_usage_content ON usage_log (content_id);

CREATE TABLE IF NOT EXISTS royalty_statements (
    id UUID PRIMARY KEY,
    license_id UUID NOT NULL REFERENCES licenses(id),
    rights_holder_id UUID NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    usage_count BIGINT NOT NULL DEFAULT 0,
    rate_applied BIGINT NOT NULL DEFAULT 0,
    total_cents BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

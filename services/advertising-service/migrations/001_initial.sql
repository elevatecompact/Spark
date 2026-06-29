CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY,
    advertiser_id UUID NOT NULL,
    name VARCHAR(256) NOT NULL,
    budget_cents BIGINT NOT NULL DEFAULT 0,
    spent_cents BIGINT NOT NULL DEFAULT 0,
    daily_budget_cents BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    targeting JSONB DEFAULT '{}'::jsonb,
    bid_strategy VARCHAR(8) NOT NULL DEFAULT 'cpm',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ad_units (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES campaigns(id),
    type VARCHAR(16) NOT NULL,
    format VARCHAR(8) NOT NULL,
    creative_url TEXT NOT NULL,
    destination_url TEXT NOT NULL DEFAULT '',
    width INT NOT NULL DEFAULT 0,
    height INT NOT NULL DEFAULT 0,
    duration_seconds INT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    approval_note TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS impressions (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL,
    ad_unit_id UUID NOT NULL,
    placement_id VARCHAR(64) NOT NULL,
    user_id UUID,
    cost_micro_cents BIGINT NOT NULL DEFAULT 0,
    device_type VARCHAR(32) NOT NULL DEFAULT '',
    geo VARCHAR(8) NOT NULL DEFAULT '',
    served_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_impressions_campaign ON impressions (campaign_id);

CREATE TABLE IF NOT EXISTS clicks (
    id UUID PRIMARY KEY,
    impression_id UUID NOT NULL REFERENCES impressions(id),
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ad_inventory (
    placement_id VARCHAR(64) PRIMARY KEY,
    content_type VARCHAR(32) NOT NULL,
    available_from TIMESTAMPTZ NOT NULL,
    available_to TIMESTAMPTZ NOT NULL,
    floor_price_micro_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);

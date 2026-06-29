CREATE TABLE gift_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL,
    match_ratio DOUBLE PRECISION NOT NULL CHECK (match_ratio > 0 AND match_ratio <= 1),
    max_match_cents BIGINT NOT NULL CHECK (max_match_cents > 0),
    total_matched BIGINT NOT NULL DEFAULT 0 CHECK (total_matched >= 0),
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_campaign_dates CHECK (end_at > start_at)
);

CREATE INDEX idx_campaigns_creator ON gift_campaigns(creator_id);
CREATE INDEX idx_campaigns_active ON gift_campaigns(start_at, end_at);

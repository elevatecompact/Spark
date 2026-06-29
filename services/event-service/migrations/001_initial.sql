CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    creator_id UUID NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category VARCHAR(64) NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL DEFAULT 'virtual',
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',
    max_attendees INT NOT NULL DEFAULT 0,
    stream_id UUID,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    cover_image_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_start ON events (start_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_creator ON events (creator_id);

CREATE TABLE IF NOT EXISTS event_ticket_tiers (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id),
    name VARCHAR(128) NOT NULL,
    price_cents BIGINT NOT NULL DEFAULT 0,
    quantity_total INT NOT NULL DEFAULT 0,
    quantity_sold INT NOT NULL DEFAULT 0,
    benefits TEXT[] NOT NULL DEFAULT '{}',
    sales_start_at TIMESTAMPTZ,
    sales_end_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS event_attendees (
    event_id UUID NOT NULL REFERENCES events(id),
    ticket_tier_id UUID REFERENCES event_ticket_tiers(id),
    user_id UUID NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'registered',
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attended_at TIMESTAMPTZ,
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS event_sessions (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id),
    title VARCHAR(256) NOT NULL,
    speaker VARCHAR(128) NOT NULL DEFAULT '',
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    stream_id UUID
);

CREATE TABLE IF NOT EXISTS event_series (
    id UUID PRIMARY KEY,
    creator_id UUID NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    frequency VARCHAR(16) NOT NULL,
    day_of_week INT,
    start_time VARCHAR(8) NOT NULL,
    timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',
    next_event_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true
);

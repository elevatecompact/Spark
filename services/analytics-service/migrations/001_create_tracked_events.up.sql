CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tracked_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    session_id VARCHAR(255) NOT NULL DEFAULT '',
    properties JSONB NOT NULL DEFAULT '{}',
    context JSONB NOT NULL DEFAULT '{}',
    event_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_name ON tracked_events(event_name);
CREATE INDEX idx_events_user ON tracked_events(user_id);
CREATE INDEX idx_events_time ON tracked_events(event_time);
CREATE INDEX idx_events_name_time ON tracked_events(event_name, event_time);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    category VARCHAR(50) DEFAULT '',
    tags TEXT[] DEFAULT '{}',
    thumbnail_url VARCHAR(500) DEFAULT '',
    stream_key VARCHAR(100) UNIQUE NOT NULL,
    rtmp_endpoint VARCHAR(500) NOT NULL,
    ingest_protocol VARCHAR(10) DEFAULT 'rtmp',
    status VARCHAR(20) DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    duration INT DEFAULT 0,
    width INT DEFAULT 0,
    height INT DEFAULT 0,
    frame_rate FLOAT DEFAULT 0,
    bitrate INT DEFAULT 0,
    codec VARCHAR(50) DEFAULT '',
    available_qualities TEXT[] DEFAULT '{}',
    viewer_count INT DEFAULT 0,
    peak_viewers INT DEFAULT 0,
    total_views BIGINT DEFAULT 0,
    record_enabled BOOLEAN DEFAULT FALSE,
    recording_id UUID,
    chat_enabled BOOLEAN DEFAULT TRUE,
    age_restricted BOOLEAN DEFAULT FALSE,
    delay_seconds INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_streams_creator_id ON streams(creator_id);
CREATE INDEX IF NOT EXISTS idx_streams_status ON streams(status);
CREATE INDEX IF NOT EXISTS idx_streams_category ON streams(category);
CREATE INDEX IF NOT EXISTS idx_streams_stream_key ON streams(stream_key);
CREATE INDEX IF NOT EXISTS idx_streams_created_at ON streams(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_streams_live ON streams(status, viewer_count DESC) WHERE status = 'live';

CREATE TABLE IF NOT EXISTS streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    stream_key VARCHAR(255) UNIQUE NOT NULL,
    playback_url TEXT,
    rtmp_url TEXT,
    status VARCHAR(50) DEFAULT 'idle',
    category VARCHAR(100),
    tags TEXT[] DEFAULT '{}',
    is_live BOOLEAN DEFAULT FALSE,
    viewer_count BIGINT DEFAULT 0,
    peak_viewer_count BIGINT DEFAULT 0,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_streams_creator ON streams(creator_id);
CREATE INDEX idx_streams_status ON streams(status);
CREATE INDEX idx_streams_category ON streams(category);
CREATE INDEX idx_streams_live ON streams(is_live) WHERE is_live = TRUE;

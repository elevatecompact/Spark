CREATE TABLE IF NOT EXISTS recordings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id UUID NOT NULL,
    creator_id UUID NOT NULL,
    title VARCHAR(200) DEFAULT '',
    s3_key VARCHAR(500) NOT NULL,
    bucket VARCHAR(100) NOT NULL,
    duration INT DEFAULT 0,
    file_size BIGINT DEFAULT 0,
    width INT DEFAULT 0,
    height INT DEFAULT 0,
    codec VARCHAR(50) DEFAULT '',
    status VARCHAR(20) DEFAULT 'processing',
    processing_progress FLOAT DEFAULT 0,
    thumbnail_key VARCHAR(500) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recordings_creator_id ON recordings(creator_id);
CREATE INDEX IF NOT EXISTS idx_recordings_stream_id ON recordings(stream_id);

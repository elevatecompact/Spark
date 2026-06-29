CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY,
    uploader_id UUID NOT NULL,
    content_type VARCHAR(16) NOT NULL,
    source_filename VARCHAR(512) NOT NULL,
    file_size_bytes BIGINT NOT NULL DEFAULT 0,
    mime_type VARCHAR(64) NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'uploading',
    storage_path TEXT NOT NULL DEFAULT '',
    cdn_url TEXT NOT NULL DEFAULT '',
    duration_seconds DOUBLE PRECISION NOT NULL DEFAULT 0,
    width INT NOT NULL DEFAULT 0,
    height INT NOT NULL DEFAULT 0,
    checksum VARCHAR(64) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS media_renditions (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL REFERENCES media_assets(id),
    profile VARCHAR(32) NOT NULL,
    format VARCHAR(8) NOT NULL,
    file_size_bytes BIGINT NOT NULL DEFAULT 0,
    storage_path TEXT NOT NULL DEFAULT '',
    cdn_url TEXT NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'processing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transcoding_jobs (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL REFERENCES media_assets(id),
    profiles JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    worker_id VARCHAR(64) NOT NULL DEFAULT '',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS drm_policies (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    content_id UUID,
    key_system VARCHAR(16) NOT NULL DEFAULT 'widevine',
    license_duration_seconds BIGINT NOT NULL DEFAULT 86400,
    security_level VARCHAR(32) NOT NULL DEFAULT 'L1',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS upload_sessions (
    id UUID PRIMARY KEY,
    uploader_id UUID NOT NULL,
    filename VARCHAR(512) NOT NULL,
    file_size_bytes BIGINT NOT NULL DEFAULT 0,
    content_type VARCHAR(64) NOT NULL DEFAULT '',
    chunks_total INT NOT NULL DEFAULT 0,
    chunks_done INT NOT NULL DEFAULT 0,
    checksum VARCHAR(64) NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'initiated',
    storage_path TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

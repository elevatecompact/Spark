CREATE TABLE IF NOT EXISTS translation_memory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_hash VARCHAR(64) NOT NULL,
    source_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    provider VARCHAR(20) NOT NULL DEFAULT 'noop',
    quality_score FLOAT NOT NULL DEFAULT 0.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source_hash, source_lang, target_lang)
);

CREATE TABLE IF NOT EXISTS translation_jobs (
    id UUID PRIMARY KEY,
    content_type VARCHAR(64) NOT NULL,
    content_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    languages TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS review_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    translation_id UUID NOT NULL REFERENCES translation_memory(id),
    original_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    reviewer_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    corrected_text TEXT,
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS provider_usage (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    request_count BIGINT NOT NULL DEFAULT 0,
    char_count BIGINT NOT NULL DEFAULT 0,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

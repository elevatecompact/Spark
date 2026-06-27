CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS creators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    bio TEXT DEFAULT '',
    avatar_url TEXT DEFAULT '',
    banner_url TEXT DEFAULT '',
    categories JSONB DEFAULT '[]'::jsonb,
    tags JSONB DEFAULT '[]'::jsonb,
    language VARCHAR(2) NOT NULL DEFAULT 'en',
    country VARCHAR(2) NOT NULL DEFAULT 'US',
    timezone VARCHAR(50) DEFAULT '',
    social_links JSONB DEFAULT '{}'::jsonb,
    verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    follower_count INTEGER DEFAULT 0,
    subscriber_count INTEGER DEFAULT 0,
    total_views BIGINT DEFAULT 0,
    total_streams INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    rank INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_creators_user_id ON creators(user_id);
CREATE INDEX IF NOT EXISTS idx_creators_language ON creators(language);
CREATE INDEX IF NOT EXISTS idx_creators_country ON creators(country);
CREATE INDEX IF NOT EXISTS idx_creators_status ON creators(status);
CREATE INDEX IF NOT EXISTS idx_creators_follower_count ON creators(follower_count DESC);
CREATE INDEX IF NOT EXISTS idx_creators_rank ON creators(rank ASC);
CREATE INDEX IF NOT EXISTS idx_creators_created_at ON creators(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_creators_categories ON creators USING GIN(categories);
CREATE INDEX IF NOT EXISTS idx_creators_tags ON creators USING GIN(tags);

CREATE INDEX IF NOT EXISTS idx_creators_search ON creators USING GIN(
    to_tsvector('english', coalesce(display_name, '') || ' ' || coalesce(bio, ''))
);

CREATE TABLE IF NOT EXISTS creator_followers (
    follower_id UUID NOT NULL,
    creator_id UUID NOT NULL,
    followed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, creator_id),
    FOREIGN KEY (creator_id) REFERENCES creators(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_creator_followers_creator ON creator_followers(creator_id);
CREATE INDEX IF NOT EXISTS idx_creator_followers_follower ON creator_followers(follower_id);
